package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/order"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
	"github.com/IgorAleksandroff/delivery/pkg/clients/queues/queues/orderstatuschangedpb"
)

type OrderProducer struct {
	topic  string
	sarama sarama.SyncProducer
}

func NewOrderProducer(brokers []string, topic string) (*OrderProducer, error) {
	version, err := sarama.ParseKafkaVersion("3.4.0")
	if err != nil {
		return nil, fmt.Errorf("parse Kafka version: %w", err)
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.RequiredAcks = sarama.WaitForAll
	saramaCfg.Producer.Retry.Max = 5
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Return.Errors = true
	saramaCfg.Producer.Partitioner = sarama.NewHashPartitioner
	saramaCfg.Version = version

	producer, err := sarama.NewSyncProducer(brokers, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("create async producer: %w", err)
	}

	return &OrderProducer{
		topic:  topic,
		sarama: producer,
	}, nil
}

func (p *OrderProducer) Publish(_ context.Context, domainEvent *order.CompletedDomainEvent) error {
	integrationEvent, err := p.mapDomainEventToIntegrationEvent(domainEvent)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(integrationEvent)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(domainEvent.GetID().String()),
		Value: sarama.ByteEncoder(bytes),
	}

	return p.sarama.SendMessages([]*sarama.ProducerMessage{msg})
}

func (p *OrderProducer) Close() error {
	return p.sarama.Close()
}

func (p *OrderProducer) mapDomainEventToIntegrationEvent(domainEvent *order.CompletedDomainEvent) (*orderstatuschangedpb.OrderStatusChangedIntegrationEvent, error) {
	if domainEvent == nil {
		return nil, errs.NewValueIsInvalidError("CompletedDomainEvent")
	}

	status, ok := orderstatuschangedpb.OrderStatus_value[domainEvent.GetOrderStatus()]
	if !ok {
		return nil, errs.NewValueIsInvalidError("OrderStatus")
	}

	integrationEvent := orderstatuschangedpb.OrderStatusChangedIntegrationEvent{
		OrderId:     domainEvent.GetOrderID().String(),
		OrderStatus: orderstatuschangedpb.OrderStatus(status),
	}
	return &integrationEvent, nil
}
