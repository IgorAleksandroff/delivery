package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"

	"github.com/IgorAleksandroff/delivery/internal/core/application/usecases/commands"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
	"github.com/IgorAleksandroff/delivery/pkg/clients/queues/queues/basketconfirmedpb"
)

type BasketConfirmedConsumer struct {
	topic                     string
	consumer                  *kafka.Consumer
	createOrderCommandHandler *commands.CreateOrderCommandHandler
}

func NewBasketConfirmedConsumer(host string, group string, topic string,
	handler *commands.CreateOrderCommandHandler) (*BasketConfirmedConsumer, error) {
	if host == "" {
		return nil, errs.NewValueIsRequiredError("host")
	}
	if group == "" {
		return nil, errs.NewValueIsRequiredError("group")
	}
	if topic == "" {
		return nil, errs.NewValueIsRequiredError("topic")
	}
	if handler == nil {
		return nil, errs.NewValueIsRequiredError("createOrderCommandHandler")
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  host,
		"group.id":           group,
		"enable.auto.commit": false,
		"auto.offset.reset":  "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	return &BasketConfirmedConsumer{
		topic:                     topic,
		consumer:                  consumer,
		createOrderCommandHandler: handler,
	}, err
}

func (c *BasketConfirmedConsumer) Close() error {
	return c.consumer.Close()
}

func (c *BasketConfirmedConsumer) Consume() error {
	err := c.consumer.Subscribe(c.topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	for {
		c.consume()
	}
}

func (c *BasketConfirmedConsumer) consume() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg, err := c.consumer.ReadMessage(-1)
	if err != nil {
		fmt.Printf("Consumer error: %v (%v)\n", err, msg)
	}

	// Обрабатываем сообщение
	fmt.Printf("Received: %s => %s\n", msg.TopicPartition, string(msg.Value))
	var event basketconfirmedpb.BasketConfirmedIntegrationEvent
	err = json.Unmarshal(msg.Value, &event)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
	}

	// Отправляем команду
	createOrderCommand, err := commands.NewCreateOrderCommand(
		createOrderID(event.BasketId), event.GetAddress().GetStreet())
	if err != nil {
		log.Printf("Failed to create changeStocks command: %v", err)
	}
	err = c.createOrderCommandHandler.Handle(ctx, createOrderCommand)
	if err != nil {
		log.Printf("Failed to handle createOrder command: %v", err)
	}

	// Подтверждаем обработку сообщения
	_, err = c.consumer.CommitMessage(msg)
	if err != nil {
		log.Printf("Commit failed: %v", err)
	}
}

func createOrderID(basketID string) uuid.UUID {
	// TODO: orderID == basketID???
	return uuid.MustParse(basketID)
}
