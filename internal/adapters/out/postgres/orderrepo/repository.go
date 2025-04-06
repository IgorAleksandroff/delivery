package orderrepo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/IgorAleksandroff/delivery/internal/adapters/out/postgres"
	"github.com/IgorAleksandroff/delivery/internal/core/domain"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/order"
	"github.com/IgorAleksandroff/delivery/internal/core/ports"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

var _ ports.OrderRepository = &Repository{}

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

func (r *Repository) Add(ctx context.Context, aggregate *order.Order) error {
	dto := DomainToDTO(aggregate)

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}
	err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&dto).Error
	if err != nil {
		return err
	}

	return r.publishDomainEvents(ctx, aggregate)
}

func (r *Repository) Update(ctx context.Context, aggregate *order.Order) error {
	dto := DomainToDTO(aggregate)

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}
	err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error
	if err != nil {
		return err
	}

	return r.publishDomainEvents(ctx, aggregate)
}

func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*order.Order, error) {
	dto := OrderDTO{}

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}
	result := tx.
		Preload(clause.Associations).
		Find(&dto, ID)
	if result.RowsAffected == 0 {
		return nil, nil
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

func (r *Repository) GetFirstInCreatedStatus(ctx context.Context) (*order.Order, error) {
	dto := OrderDTO{}

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}
	result := tx.
		Preload(clause.Associations).
		Where("status = ?", order.StatusCreated).
		First(&dto)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NewObjectNotFoundError("Free courier", nil)
		}
		return nil, result.Error
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

func (r *Repository) GetAllInAssignedStatus(ctx context.Context) ([]*order.Order, error) {
	var dtos []OrderDTO

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}
	result := tx.
		Preload(clause.Associations).
		Where("status = ?", order.StatusAssigned).
		Find(&dtos)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errs.NewObjectNotFoundError("Assigned orders", nil)
	}

	aggregates := make([]*order.Order, len(dtos))
	for i, dto := range dtos {
		aggregates[i] = DtoToDomain(dto)
	}

	return aggregates, nil
}

func (r *Repository) publishDomainEvents(ctx context.Context, aggregate domain.AggregateRoot) error {
	for _, event := range aggregate.GetDomainEvents() {
		switch event.(type) {
		case order.CompletedDomainEvent:
			err := mediatr.Publish[order.CompletedDomainEvent](ctx, event.(order.CompletedDomainEvent))
			if err != nil {
				return err
			}
		}
	}
	aggregate.ClearDomainEvents()
	return nil
}
