package courierrepo

import (
	"context"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/postgres"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/courier"
	"github.com/IgorAleksandroff/delivery/internal/core/ports"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

var _ ports.CourierRepository = &Repository{}

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

func (r *Repository) Add(ctx context.Context, aggregate *courier.Courier) error {
	dto := DomainToDTO(aggregate)

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}
	err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&dto).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, aggregate *courier.Courier) error {
	dto := DomainToDTO(aggregate)

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}
	err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*courier.Courier, error) {
	dto := CourierDTO{}

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}
	result := tx.
		Preload(clause.Associations).
		Find(&dto, ID)
	if result.RowsAffected == 0 {
		return nil, errs.NewObjectNotFoundError(ID.String(), ID)
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

func (r *Repository) GetAllInFreeStatus(ctx context.Context) ([]*courier.Courier, error) {
	var dtos []CourierDTO

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}
	result := tx.
		Preload(clause.Associations).
		Where("status = ?", courier.StatusFree).
		Find(&dtos)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errs.NewObjectNotFoundError("Free couriers", nil)
	}

	aggregates := make([]*courier.Courier, len(dtos))
	for i, dto := range dtos {
		aggregates[i] = DtoToDomain(dto)
	}

	return aggregates, nil
}
