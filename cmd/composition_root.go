package cmd

import (
	"github.com/IgorAleksandroff/delivery/internal/adapters/jobs"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/outbox"
	"github.com/IgorAleksandroff/delivery/internal/core/application/eventhandlers"
	"log"
	"reflect"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"github.com/IgorAleksandroff/delivery/internal/adapters/in/kafka"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/grpc/geo"
	kafkaout "github.com/IgorAleksandroff/delivery/internal/adapters/out/kafka"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/postgres"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/postgres/courierrepo"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/postgres/orderrepo"
	"github.com/IgorAleksandroff/delivery/internal/core/application/usecases/commands"
	"github.com/IgorAleksandroff/delivery/internal/core/application/usecases/queries"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/orders"
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
	Consumers       Consumers
	Producers       Producers
	EventRegistry   outbox.EventRegistry

	closeFns []func() error
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
	OutboxJob       cron.Job
}

type Consumers struct {
	BasketConfirmedConsumer *kafka.BasketConfirmedConsumer
}

type Producers struct {
	OrderChangedProducer ports.OrderProducer
}

func NewCompositionRoot(gormDb *gorm.DB, cfg Config) CompositionRoot {
	// Domain Services
	orderDispatcher := services.NewOrderDispatcher()

	// Message Registry
	eventRegistry, err := outbox.NewEventRegistry()
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}
	err = eventRegistry.RegisterDomainEvent(reflect.TypeOf(orders.CompletedDomainEvent{}))
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

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

	outboxRepository, err := outbox.NewRepository(gormDb)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Grpc Clients
	geoClient, err := geo.NewClient(cfg.GeoServiceGrpcHost)
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

	outboxJob, err := jobs.NewOutboxJob(outboxRepository, eventRegistry)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Kafka Consumers
	basketConfirmedConsumer, err := kafka.NewBasketConfirmedConsumer(cfg.KafkaHost, cfg.KafkaConsumerGroup,
		cfg.KafkaBasketConfirmedTopic, createOrderCommandHandler)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Kafka Producers
	brokers := []string{cfg.KafkaHost}
	orderKafkaProducer, err := kafkaout.NewOrderProducer(brokers, cfg.KafkaOrderChangedTopic)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Domain Event Handlers
	orderDomainEventHandler, err := eventhandlers.NewOrderCompleted(orderKafkaProducer)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Mediatr Subscribes
	err = mediatr.RegisterNotificationHandlers[*orders.CompletedDomainEvent](orderDomainEventHandler)
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
			OutboxJob:       outboxJob,
		},
		Consumers: Consumers{
			BasketConfirmedConsumer: basketConfirmedConsumer,
		},
		Producers: Producers{
			OrderChangedProducer: orderKafkaProducer,
		},
		EventRegistry: eventRegistry,
	}

	// Close
	compositionRoot.closeFns = append(compositionRoot.closeFns, geoClient.Close)
	compositionRoot.closeFns = append(compositionRoot.closeFns, basketConfirmedConsumer.Close)
	compositionRoot.closeFns = append(compositionRoot.closeFns, orderKafkaProducer.Close)
	return compositionRoot
}

func (cr *CompositionRoot) Close() {
	for _, fn := range cr.closeFns {
		if err := fn(); err != nil {
			log.Printf("ошибка при закрытии зависимости: %v", err)
		}
	}
}
