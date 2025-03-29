package queries

import (
	"gorm.io/gorm"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

type GetAllCouriersQueryHandler struct {
	db *gorm.DB
}

func NewGetAllCouriersQueryHandler(db *gorm.DB) (*GetAllCouriersQueryHandler, error) {
	if db == nil {
		return &GetAllCouriersQueryHandler{}, errs.NewValueIsRequiredError("db")
	}
	return &GetAllCouriersQueryHandler{db: db}, nil
}

func (q *GetAllCouriersQueryHandler) Handle(query GetAllCouriersQuery) (GetAllCouriersResponse, error) {
	if query.isEmpty() {
		return GetAllCouriersResponse{}, errs.NewValueIsRequiredError("query")
	}

	var couriers []CourierResponse
	result := q.db.Raw("SELECT id,name, location_x, location_y FROM couriers").Scan(&couriers)

	if result.Error != nil {
		return GetAllCouriersResponse{}, result.Error
	}

	return GetAllCouriersResponse{Couriers: couriers}, nil
}

type GetAllCouriersQuery struct {
	isSet bool
}

func NewGetAllCouriersQuery() (GetAllCouriersQuery, error) {
	return GetAllCouriersQuery{isSet: true}, nil
}
func (q GetAllCouriersQuery) isEmpty() bool {
	return !q.isSet
}

type GetAllCouriersResponse struct {
	Couriers []CourierResponse
}

type CourierResponse struct {
	ID       uuid.UUID
	Name     string
	Location LocationResponse
}

type LocationResponse struct {
	X int
	Y int
}
