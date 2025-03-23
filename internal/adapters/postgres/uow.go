package postgres

import (
	"context"

	"gorm.io/gorm"

	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
	"github.com/IgorAleksandroff/delivery/internal/pkg/uow"
)

var _ uow.UnitOfWork = &GormUnitOfWork{}

type GormUnitOfWork struct {
	db *gorm.DB
}

func NewGormUnitOfWork(db *gorm.DB) (*GormUnitOfWork, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}
	return &GormUnitOfWork{db: db}, nil
}

func (u *GormUnitOfWork) Begin(ctx context.Context) context.Context {
	tx := u.db.Begin()
	return NewContextWithTx(ctx, tx)
}

func (u *GormUnitOfWork) Commit(ctx context.Context) error {
	tx := GetTxFromContext(ctx, nil)
	if tx != nil {
		return tx.Commit().Error
	}
	return nil
}

func (u *GormUnitOfWork) Rollback(ctx context.Context) error {
	tx := GetTxFromContext(ctx, nil)
	if tx != nil {
		return tx.Rollback().Error
	}
	return nil
}

type txKey struct{}

func NewContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func GetTxFromContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if ok {
		return tx // 🔥 Используем транзакцию, если она есть
	}
	return db // Если нет транзакции, работаем с `db`
}
