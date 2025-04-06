package orderrepo

import (
	"context"
	"errors"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/order"
	"github.com/IgorAleksandroff/delivery/internal/core/usecases/commands"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
	"github.com/IgorAleksandroff/delivery/internal/repository"
	"github.com/IgorAleksandroff/delivery/internal/repository/outbox"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ commands.OrderRepository = &Repository{}

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
	outboxEvents, err := outbox.EncodeDomainEvents(aggregate.GetDomainEvents())
	if err != nil {
		return err
	}

	tx := repository.GetTxFromContext(ctx)
	isTransaction := tx == nil
	if isTransaction {
		tx = r.db.Begin()
		defer tx.Rollback()
	}

	err = tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&dto).Error
	if err != nil {
		return err
	}

	if len(outboxEvents) > 0 {
		err = tx.Create(&outboxEvents).Error
		if err != nil {
			return err
		}
	}

	if isTransaction {
		return tx.Commit().Error
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, aggregate *order.Order) error {
	dto := DomainToDTO(aggregate)
	outboxEvents, err := outbox.EncodeDomainEvents(aggregate.GetDomainEvents())
	if err != nil {
		return err
	}

	tx := repository.GetTxFromContext(ctx)
	isTransaction := tx == nil
	if isTransaction {
		tx = r.db.Begin()
		defer tx.Rollback()
	}

	err = tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error
	if err != nil {
		return err
	}

	if len(outboxEvents) > 0 {
		err = tx.Create(&outboxEvents).Error
		if err != nil {
			return err
		}
	}

	if isTransaction {
		return tx.Commit().Error
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*order.Order, error) {
	dto := OrderDTO{}

	tx := repository.GetTxFromContext(ctx)
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

	tx := repository.GetTxFromContext(ctx)
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

	tx := repository.GetTxFromContext(ctx)
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
