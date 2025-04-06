package outbox

import (
	"context"

	"gorm.io/gorm"

	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

type OutboxRepository interface {
	Update(ctx context.Context, event *Message) error
	GetNotPublishedMessages() ([]*Message, error)
}

var _ OutboxRepository = &Repository{}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Update(ctx context.Context, outboxEvent *Message) error {
	err := r.db.WithContext(ctx).Save(&outboxEvent).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetNotPublishedMessages() ([]*Message, error) {
	var events []*Message
	result := r.db.
		Order("occurred_at ASC").
		Limit(20).
		Where("processed_at IS NULL").Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return events, nil
}
