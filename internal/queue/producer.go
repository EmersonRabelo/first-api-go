package queue

import (
	"context"
	"encoding/json"

	contract "github.com/EmersonRabelo/first-api-go/internal/dtos/report/message"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Estrutura concreta do Producer
type ReportProducer struct {
	channel    *amqp.Channel
	exchange   string
	routingKey string
}

// Construtor
func NewReportProducer(channel *amqp.Channel, exchange, routingKey string) *ReportProducer {
	return &ReportProducer{
		channel:    channel,
		exchange:   exchange,
		routingKey: routingKey,
	}
}

func (producer *ReportProducer) Publish(context context.Context, message *contract.CreateReportMessage) error {
	body, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return producer.channel.PublishWithContext(
		context,
		producer.exchange,
		producer.routingKey,
		false,
		false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
