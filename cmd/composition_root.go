package cmd

import (
	"context"
	"github.com/IgorAleksandroff/delivery/internal/adapters/in/jobs"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/postgres"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/postgres/courierrepo"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/postgres/orderrepo"
	"github.com/IgorAleksandroff/delivery/internal/core/application/usecases"
	"github.com/IgorAleksandroff/delivery/internal/core/application/usecases/queries"
	"github.com/robfig/cron/v3"
	"log"

	"gorm.io/gorm"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/services"
	"github.com/IgorAleksandroff/delivery/internal/core/ports"
	"github.com/IgorAleksandroff/delivery/internal/pkg/uow"
)

type CompositionRoot struct {
	DomainServices  DomainServices
	Repositories    Repositories
	CommandHandlers CommandHandlers
	QueryHandlers   QueryHandlers
	Jobs            Jobs
}

type DomainServices struct {
	OrderDispatcher *services.Dispatcher
}

type Repositories struct {
	UnitOfWork        uow.UnitOfWork
	OrderRepository   ports.OrderRepository
	CourierRepository ports.CourierRepository
}

type CommandHandlers struct {
	AssignOrdersCommandHandler *usecases.AssignOrdersCommandHandler
	CreateOrderCommandHandler  *usecases.CreateOrderCommandHandler
	MoveCouriersCommandHandler *usecases.MoveCouriersCommandHandler
}

type QueryHandlers struct {
	GetAllCouriersQueryHandler        *queries.GetAllCouriersQueryHandler
	GetNotCompletedOrdersQueryHandler *queries.GetNotCompletedOrdersQueryHandler
}

type Jobs struct {
	AssignOrdersJob cron.Job
	MoveCouriersJob cron.Job
}

func NewCompositionRoot(ctx context.Context, gormDb *gorm.DB) CompositionRoot {
	// Domain Services
	orderDispatcher := services.NewOrderDispatcher()

	// Repositories
	unitOfWork, err := postgres.NewGormUnitOfWork(gormDb)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	orderRepository, err := orderrepo.NewRepository(gormDb)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	courierRepository, err := courierrepo.NewRepository(gormDb)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Command Handlers
	createOrderCommandHandler, err := usecases.NewCreateOrderCommandHandler(orderRepository)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	assignOrdersCommandHandler, err := usecases.NewAssignOrdersCommandHandler(
		unitOfWork, orderRepository, courierRepository, orderDispatcher)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	moveCouriersCommandHandler, err := usecases.NewMoveCouriersCommandHandler(
		unitOfWork, orderRepository, courierRepository)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Query Handlers
	getAllCouriersQueryHandler, err := queries.NewGetAllCouriersQueryHandler(gormDb)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	getNotCompletedOrdersQueryHandler, err := queries.NewGetNotCompletedOrdersQueryHandler(gormDb)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Jobs
	assignOrdersJob, err := jobs.NewAssignOrdersJob(assignOrdersCommandHandler)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	moveCouriersJob, err := jobs.NewMoveCouriersJob(moveCouriersCommandHandler)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	compositionRoot := CompositionRoot{
		DomainServices: DomainServices{
			OrderDispatcher: orderDispatcher,
		},

		Repositories: Repositories{
			OrderRepository:   orderRepository,
			CourierRepository: courierRepository,
		},
		CommandHandlers: CommandHandlers{
			AssignOrdersCommandHandler: assignOrdersCommandHandler,
			CreateOrderCommandHandler:  createOrderCommandHandler,
			MoveCouriersCommandHandler: moveCouriersCommandHandler,
		},
		QueryHandlers: QueryHandlers{
			GetAllCouriersQueryHandler:        getAllCouriersQueryHandler,
			GetNotCompletedOrdersQueryHandler: getNotCompletedOrdersQueryHandler,
		},
		Jobs: Jobs{
			AssignOrdersJob: assignOrdersJob,
			MoveCouriersJob: moveCouriersJob,
		},
	}

	return compositionRoot
}
