package cmd

import (
	"log"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"github.com/IgorAleksandroff/delivery/internal/adapters/in/jobs"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/grpc/geo"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/postgres"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/postgres/courierrepo"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/postgres/orderrepo"
	"github.com/IgorAleksandroff/delivery/internal/core/application/usecases/commands"
	"github.com/IgorAleksandroff/delivery/internal/core/application/usecases/queries"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/services"
	"github.com/IgorAleksandroff/delivery/internal/core/ports"
	"github.com/IgorAleksandroff/delivery/internal/pkg/uow"
)

type CompositionRoot struct {
	DomainServices  DomainServices
	Repositories    Repositories
	CommandHandlers CommandHandlers
	QueryHandlers   QueryHandlers
	Clients         Clients
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
	AssignOrdersCommandHandler *commands.AssignOrdersCommandHandler
	CreateOrderCommandHandler  *commands.CreateOrderCommandHandler
	MoveCouriersCommandHandler *commands.MoveCouriersCommandHandler
}

type QueryHandlers struct {
	GetAllCouriersQueryHandler        *queries.GetAllCouriersQueryHandler
	GetNotCompletedOrdersQueryHandler *queries.GetNotCompletedOrdersQueryHandler
}

type Clients struct {
	GeoClient ports.GeoClient
}

type Jobs struct {
	AssignOrdersJob cron.Job
	MoveCouriersJob cron.Job
}

func NewCompositionRoot(gormDb *gorm.DB, geoServiceGrpcHost string) CompositionRoot {
	// Domain Services
	orderDispatcher := services.NewOrderDispatcher()

	// Repositories
	unitOfWork, err := postgres.NewUnitOfWork(gormDb)
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

	// Grpc Clients
	geoClient, err := geo.NewClient(geoServiceGrpcHost)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Command Handlers
	createOrderCommandHandler, err := commands.NewCreateOrderCommandHandler(orderRepository, geoClient)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	assignOrdersCommandHandler, err := commands.NewAssignOrdersCommandHandler(
		unitOfWork, orderRepository, courierRepository, orderDispatcher)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	moveCouriersCommandHandler, err := commands.NewMoveCouriersCommandHandler(
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
		Clients: Clients{
			GeoClient: geoClient,
		},
		Jobs: Jobs{
			AssignOrdersJob: assignOrdersJob,
			MoveCouriersJob: moveCouriersJob,
		},
	}

	return compositionRoot
}
