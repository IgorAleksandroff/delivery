package jobs

import (
	"context"
	"github.com/IgorAleksandroff/delivery/internal/adapters/out/outbox"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/orders"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
	"github.com/labstack/gommon/log"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/robfig/cron/v3"
	"time"
)

var _ cron.Job = &OutboxJob{}

type OutboxJob struct {
	outboxRepository outbox.OutboxRepository
	eventRegistry    outbox.EventRegistry
}

func NewOutboxJob(outboxRepository outbox.OutboxRepository, eventRegistry outbox.EventRegistry) (*OutboxJob, error) {
	if outboxRepository == nil {
		return nil, errs.NewValueIsRequiredError("outboxRepository")
	}
	if eventRegistry == nil {
		return nil, errs.NewValueIsRequiredError("eventRegistry")
	}

	return &OutboxJob{
		outboxRepository: outboxRepository,
		eventRegistry:    eventRegistry}, nil
}

func (j *OutboxJob) Run() {
	ctx := context.Background()

	// Получаем не отправленные Outbox Events
	outboxMessages, err := j.outboxRepository.GetNotPublishedMessages()
	if err != nil {
		return
	}

	// Перебираем в цикле
	for _, outboxMessage := range outboxMessages {
		// Приводим Outbox Message -> Domain Event
		domainEvent, err := j.eventRegistry.DecodeDomainEvent(outboxMessage)
		if err != nil {
			log.Error(err)
			return
		}

		// Go не поддерживает вызов generic-функций с параметрами T во время выполнения
		// Поэтому делаем Switch
		switch domainEvent.(type) {
		case *orders.CompletedDomainEvent:
			err := mediatr.Publish[*orders.CompletedDomainEvent](ctx, domainEvent.(*orders.CompletedDomainEvent))
			if err != nil {
				log.Error(err)
				continue
			}
		}

		// Если ошибок нет, помечаем Outbox Message как отправленное и сохраняем в БД
		// А если были ошибки, то цикл просто повторяется
		now := time.Now().UTC()
		outboxMessage.ProcessedAt = &now
		err = j.outboxRepository.Update(ctx, outboxMessage)
		if err != nil {
			log.Error(err)
			continue
		}
	}
}
