package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/courier"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/order"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/services"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

func TestAssignOrdersCommandHandler_Handle(t *testing.T) {
	testCases := []struct {
		name                 string
		command              AssignOrdersCommand
		setupStubs           func() (*stubUnitOfWork, *stubOrderRepository, *stubCourierRepository)
		expectedError        error
		checkStateAfterError func(*testing.T, *stubUnitOfWork, *stubOrderRepository, *stubCourierRepository)
	}{
		{
			name:    "Empty command should return error",
			command: AssignOrdersCommand{}, // Empty command
			setupStubs: func() (*stubUnitOfWork, *stubOrderRepository, *stubCourierRepository) {
				return &stubUnitOfWork{}, &stubOrderRepository{}, &stubCourierRepository{}
			},
			expectedError: errs.NewValueIsRequiredError("add address command"),
			checkStateAfterError: func(t *testing.T, uow *stubUnitOfWork, orderRepo *stubOrderRepository, courierRepo *stubCourierRepository) {
				if uow.beginCalled || uow.commitCalled || orderRepo.updateCalled || courierRepo.updateCalled {
					t.Error("No operations should be performed for empty command")
				}
			},
		},
		{
			name: "No orders available should return NotAvailableOrders",
			command: func() AssignOrdersCommand {
				cmd, _ := NewAssignOrdersCommand()
				return cmd
			}(),
			setupStubs: func() (*stubUnitOfWork, *stubOrderRepository, *stubCourierRepository) {
				return &stubUnitOfWork{}, &stubOrderRepository{order: nil}, &stubCourierRepository{}
			},
			expectedError: NotAvailableOrders,
			checkStateAfterError: func(t *testing.T, uow *stubUnitOfWork, orderRepo *stubOrderRepository, courierRepo *stubCourierRepository) {
				if uow.beginCalled || uow.commitCalled || orderRepo.updateCalled || courierRepo.updateCalled {
					t.Error("No operations should be performed when no orders are available")
				}
			},
		},
		{
			name: "GetFirstInCreatedStatus error should be returned",
			command: func() AssignOrdersCommand {
				cmd, _ := NewAssignOrdersCommand()
				return cmd
			}(),
			setupStubs: func() (*stubUnitOfWork, *stubOrderRepository, *stubCourierRepository) {
				return &stubUnitOfWork{},
					&stubOrderRepository{getFirstError: errors.New("database error")},
					&stubCourierRepository{}
			},
			expectedError: errors.New("database error"),
			checkStateAfterError: func(t *testing.T, uow *stubUnitOfWork, orderRepo *stubOrderRepository, courierRepo *stubCourierRepository) {
				if uow.beginCalled || uow.commitCalled || orderRepo.updateCalled || courierRepo.updateCalled {
					t.Error("No operations should be performed when GetFirstInCreatedStatus returns an error")
				}
			},
		},
		{
			name: "No couriers available should return NotAvailableCouriers",
			command: func() AssignOrdersCommand {
				cmd, _ := NewAssignOrdersCommand()
				return cmd
			}(),
			setupStubs: func() (*stubUnitOfWork, *stubOrderRepository, *stubCourierRepository) {
				return &stubUnitOfWork{},
					&stubOrderRepository{order: &order.Order{}},
					&stubCourierRepository{couriers: []*courier.Courier{}}
			},
			expectedError: NotAvailableCouriers,
			checkStateAfterError: func(t *testing.T, uow *stubUnitOfWork, orderRepo *stubOrderRepository, courierRepo *stubCourierRepository) {
				if uow.beginCalled || uow.commitCalled || orderRepo.updateCalled || courierRepo.updateCalled {
					t.Error("No operations should be performed when no couriers are available")
				}
			},
		},
		{
			name: "GetAllInFreeStatus error should be returned",
			command: func() AssignOrdersCommand {
				cmd, _ := NewAssignOrdersCommand()
				return cmd
			}(),
			setupStubs: func() (*stubUnitOfWork, *stubOrderRepository, *stubCourierRepository) {
				return &stubUnitOfWork{},
					&stubOrderRepository{order: &order.Order{}},
					&stubCourierRepository{getAllError: errors.New("database error")}
			},
			expectedError: errors.New("database error"),
			checkStateAfterError: func(t *testing.T, uow *stubUnitOfWork, orderRepo *stubOrderRepository, courierRepo *stubCourierRepository) {
				if uow.beginCalled || uow.commitCalled || orderRepo.updateCalled || courierRepo.updateCalled {
					t.Error("No operations should be performed when GetAllInFreeStatus returns an error")
				}
			},
		},
		{
			name: "Order update error should be returned and transaction rolled back",
			command: func() AssignOrdersCommand {
				cmd, _ := NewAssignOrdersCommand()
				return cmd
			}(),
			setupStubs: func() (*stubUnitOfWork, *stubOrderRepository, *stubCourierRepository) {
				testOrder := order.MustNewOrder(uuid.New(), kernel.CreateRandomLocation())
				testCourier := courier.MustNewCourier("courier-1", "transport", 1, kernel.CreateRandomLocation())
				return &stubUnitOfWork{},
					&stubOrderRepository{
						order:       testOrder,
						updateError: errors.New("order update error"),
					},
					&stubCourierRepository{couriers: []*courier.Courier{testCourier}}
			},
			expectedError: errors.New("order update error"),
			checkStateAfterError: func(t *testing.T, uow *stubUnitOfWork, orderRepo *stubOrderRepository, courierRepo *stubCourierRepository) {
				if !uow.beginCalled || uow.commitCalled || !orderRepo.updateCalled || courierRepo.updateCalled {
					t.Error("Expected begin and order update to be called, but commit and courier update should not be called")
				}
			},
		},
		{
			name: "Courier update error should be returned and transaction rolled back",
			command: func() AssignOrdersCommand {
				cmd, _ := NewAssignOrdersCommand()
				return cmd
			}(),
			setupStubs: func() (*stubUnitOfWork, *stubOrderRepository, *stubCourierRepository) {
				testOrder := order.MustNewOrder(uuid.New(), kernel.CreateRandomLocation())
				testCourier := courier.MustNewCourier("courier-1", "transport", 1, kernel.CreateRandomLocation())
				return &stubUnitOfWork{},
					&stubOrderRepository{order: testOrder},
					&stubCourierRepository{
						couriers:    []*courier.Courier{testCourier},
						updateError: errors.New("courier update error"),
					}
			},
			expectedError: errors.New("courier update error"),
			checkStateAfterError: func(t *testing.T, uow *stubUnitOfWork, orderRepo *stubOrderRepository, courierRepo *stubCourierRepository) {
				if !uow.beginCalled || uow.commitCalled || !orderRepo.updateCalled || !courierRepo.updateCalled {
					t.Error("Expected begin, order update and courier update to be called, but commit should not be called")
				}
			},
		},
		{
			name: "Commit error should be returned",
			command: func() AssignOrdersCommand {
				cmd, _ := NewAssignOrdersCommand()
				return cmd
			}(),
			setupStubs: func() (*stubUnitOfWork, *stubOrderRepository, *stubCourierRepository) {
				testOrder := order.MustNewOrder(uuid.New(), kernel.CreateRandomLocation())
				testCourier := courier.MustNewCourier("courier-1", "transport", 1, kernel.CreateRandomLocation())
				return &stubUnitOfWork{commitError: errors.New("commit error")},
					&stubOrderRepository{order: testOrder},
					&stubCourierRepository{couriers: []*courier.Courier{testCourier}}
			},
			expectedError: errors.New("commit error"),
			checkStateAfterError: func(t *testing.T, uow *stubUnitOfWork, orderRepo *stubOrderRepository, courierRepo *stubCourierRepository) {
				if !uow.beginCalled || !uow.commitCalled || !orderRepo.updateCalled || !courierRepo.updateCalled {
					t.Error("Expected all operations to be called")
				}
			},
		},
		{
			name: "Successful execution",
			command: func() AssignOrdersCommand {
				cmd, _ := NewAssignOrdersCommand()
				return cmd
			}(),
			setupStubs: func() (*stubUnitOfWork, *stubOrderRepository, *stubCourierRepository) {
				testOrder := order.MustNewOrder(uuid.New(), kernel.CreateRandomLocation())
				testCourier := courier.MustNewCourier("courier-1", "transport", 1, kernel.CreateRandomLocation())
				return &stubUnitOfWork{},
					&stubOrderRepository{order: testOrder},
					&stubCourierRepository{couriers: []*courier.Courier{testCourier}}
			},
			expectedError: nil,
			checkStateAfterError: func(t *testing.T, uow *stubUnitOfWork, orderRepo *stubOrderRepository, courierRepo *stubCourierRepository) {
				if !uow.beginCalled || !uow.commitCalled || !orderRepo.updateCalled || !courierRepo.updateCalled {
					t.Error("Expected all operations to be called for successful execution")
				}
				if courierRepo.updatedCourier == nil || courierRepo.updatedCourier.Name() != "courier-1" {
					t.Error("Expected correct courier to be updated")
				}
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup stubs
			uowStub, orderRepoStub, courierRepoStub := tc.setupStubs()

			// Create the handler with stubs
			handler, err := NewAssignOrdersCommandHandler(
				uowStub,
				orderRepoStub,
				courierRepoStub,
				services.NewOrderDispatcher(),
			)
			if err != nil {
				t.Fatalf("Failed to create handler: %v", err)
			}

			// Call the method being tested
			err = handler.Handle(context.Background(), tc.command)

			// Check results
			if (tc.expectedError == nil && err != nil) ||
				(tc.expectedError != nil && (err == nil || err.Error() != tc.expectedError.Error())) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}

			// Additional checks on stub state
			if tc.checkStateAfterError != nil {
				tc.checkStateAfterError(t, uowStub, orderRepoStub, courierRepoStub)
			}
		})
	}
}

// Stubs for dependencies
type stubUnitOfWork struct {
	beginCalled    bool
	commitCalled   bool
	rollbackCalled bool
	commitError    error
}

func (s *stubUnitOfWork) Begin(ctx context.Context) context.Context {
	s.beginCalled = true
	return ctx
}

func (s *stubUnitOfWork) Commit(ctx context.Context) error {
	s.commitCalled = true
	return s.commitError
}

func (s *stubUnitOfWork) Rollback(ctx context.Context) error {
	s.rollbackCalled = true
	return nil
}

type stubOrderRepository struct {
	order         *order.Order
	getFirstError error
	updateCalled  bool
	updateError   error
}

func (s *stubOrderRepository) Add(ctx context.Context, aggregate *order.Order) error {
	return nil
}

func (s *stubOrderRepository) Get(ctx context.Context, ID uuid.UUID) (*order.Order, error) {
	return s.order, nil
}

func (s *stubOrderRepository) GetAllInAssignedStatus(ctx context.Context) ([]*order.Order, error) {
	return []*order.Order{s.order}, nil
}

func (s *stubOrderRepository) GetFirstInCreatedStatus(ctx context.Context) (*order.Order, error) {
	return s.order, s.getFirstError
}

func (s *stubOrderRepository) Update(ctx context.Context, order *order.Order) error {
	s.updateCalled = true
	return s.updateError
}

type stubCourierRepository struct {
	couriers       []*courier.Courier
	getAllError    error
	updateCalled   bool
	updateError    error
	updatedCourier *courier.Courier
}

func (s *stubCourierRepository) Add(ctx context.Context, aggregate *courier.Courier) error {
	return nil
}

func (s *stubCourierRepository) Get(ctx context.Context, ID uuid.UUID) (*courier.Courier, error) {
	return s.couriers[0], nil
}

func (s *stubCourierRepository) GetAllInFreeStatus(ctx context.Context) ([]*courier.Courier, error) {
	return s.couriers, s.getAllError
}

func (s *stubCourierRepository) Update(ctx context.Context, courier *courier.Courier) error {
	s.updateCalled = true
	s.updatedCourier = courier
	return s.updateError
}
