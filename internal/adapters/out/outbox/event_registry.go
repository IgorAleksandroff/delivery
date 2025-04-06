package outbox

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/IgorAleksandroff/delivery/internal/core/domain"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

type IEventRegistry interface {
	DecodeDomainEvent(event *Message) (domain.Event, error)
}

var _ IEventRegistry = &EventRegistry{}

type EventRegistry struct {
	EventRegistry map[string]reflect.Type
}

func NewEventRegistry() (*EventRegistry, error) {
	return &EventRegistry{
		EventRegistry: make(map[string]reflect.Type),
	}, nil
}

func (r *EventRegistry) RegisterDomainEvent(eventType reflect.Type) error {
	if eventType == nil {
		return errs.NewValueIsRequiredError("eventType")
	}
	eventName := eventType.Name()
	r.EventRegistry[eventName] = eventType
	return nil
}

func EncodeDomainEvent(domainEvent domain.Event) (Message, error) {
	payload, err := json.Marshal(domainEvent)
	if err != nil {
		return Message{}, fmt.Errorf("failed to marshal event: %w", err)
	}

	return Message{
		ID:          domainEvent.GetID(),
		Name:        domainEvent.GetName(),
		Payload:     payload,
		OccurredAt:  time.Now().UTC(),
		ProcessedAt: nil,
	}, nil
}

func EncodeDomainEvents(domainEvent []domain.Event) ([]Message, error) {
	outboxMessages := make([]Message, 0)
	for _, event := range domainEvent {
		event, err := EncodeDomainEvent(event)
		if err != nil {
			return nil, err
		}
		outboxMessages = append(outboxMessages, event)
	}
	return outboxMessages, nil
}

func (r *EventRegistry) DecodeDomainEvent(outboxMessage *Message) (domain.Event, error) {
	t, ok := r.EventRegistry[outboxMessage.Name]
	if !ok {
		return nil, fmt.Errorf("unknown outboxMessage type: %s", outboxMessage.Name)
	}

	// Создаём новый указатель на нужный тип
	eventPtr := reflect.New(t).Interface()

	if err := json.Unmarshal(outboxMessage.Payload, eventPtr); err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	// Приводим к DomainEvent
	domainEvent, ok := eventPtr.(domain.Event)
	if !ok {
		return nil, fmt.Errorf("decoded outboxMessage does not implement DomainEvent")
	}

	return domainEvent, nil
}
