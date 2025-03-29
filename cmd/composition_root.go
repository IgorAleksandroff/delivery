package cmd

import (
	"context"
	"log"

	"gorm.io/gorm"
	
	"github.com/IgorAleksandroff/delivery/internal/adapters/postgres/courierrepo"
	"github.com/IgorAleksandroff/delivery/internal/adapters/postgres/orderrepo"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/services"
	"github.com/IgorAleksandroff/delivery/internal/core/ports"
	"github.com/IgorAleksandroff/delivery/internal/pkg/uow"
)

type CompositionRoot struct {
	DomainServices DomainServices
	Repositories   Repositories
}

type DomainServices struct {
	OrderDispatcher *services.Dispatcher
}

type Repositories struct {
	UnitOfWork        uow.UnitOfWork
	OrderRepository   ports.OrderRepository
	CourierRepository ports.CourierRepository
}

func NewCompositionRoot(ctx context.Context, gormDb *gorm.DB) CompositionRoot {
	// Domain Services
	orderDispatcher := services.NewOrderDispatcher()

	// Repositories
	//unitOfWork, err := postgres.NewUnitOfWork(gormDb)
	//if err != nil {
	//	log.Fatalf("run application error: %s", err)
	//}
	//
	//ctx = unitOfWork.Begin(ctx)
	//defer func() {
	//	err := unitOfWork.Rollback(ctx)
	//	if err != nil {
	//		log.Println("Rollback error:", err)
	//	}
	//}()

	orderRepository, err := orderrepo.NewRepository(gormDb)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	courierRepository, err := courierrepo.NewRepository(gormDb)
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
	}

	return compositionRoot
}
