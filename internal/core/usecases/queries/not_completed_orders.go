package queries

import (
	"gorm.io/gorm"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/order"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

type GetNotCompletedOrdersQueryHandler struct {
	db *gorm.DB
}

func NewGetNotCompletedOrdersQueryHandler(db *gorm.DB) (*GetNotCompletedOrdersQueryHandler, error) {
	if db == nil {
		return &GetNotCompletedOrdersQueryHandler{}, errs.NewValueIsRequiredError("db")
	}
	return &GetNotCompletedOrdersQueryHandler{db: db}, nil
}

func (q *GetNotCompletedOrdersQueryHandler) Handle(query GetNotCompletedOrdersQuery) (GetNotCompletedOrdersResponse, error) {
	if query.isEmpty() {
		return GetNotCompletedOrdersResponse{}, errs.NewValueIsRequiredError("query")
	}

	var orders []OrderResponse
	result := q.db.Raw("SELECT id, courier_id, location_x, location_y, status FROM public.orders where status!=?",
		order.StatusCompleted).Scan(&orders)

	if result.Error != nil {
		return GetNotCompletedOrdersResponse{}, result.Error
	}

	return GetNotCompletedOrdersResponse{Orders: orders}, nil
}

type GetNotCompletedOrdersQuery struct {
	isSet bool
}

func NewGetNotCompletedOrdersQuery() (GetNotCompletedOrdersQuery, error) {
	return GetNotCompletedOrdersQuery{isSet: true}, nil
}
func (q GetNotCompletedOrdersQuery) isEmpty() bool {
	return !q.isSet
}

type GetNotCompletedOrdersResponse struct {
	Orders []OrderResponse
}

type OrderResponse struct {
	ID       uuid.UUID        `gorm:"type:uuid;primaryKey"`
	Location LocationResponse `gorm:"embedded;embeddedPrefix:location_"`
}
