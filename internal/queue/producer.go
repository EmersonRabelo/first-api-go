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

func (p *ReportProducer) Publish(ctx context.Context, message *contract.CreateReportMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.channel.PublishWithContext(
		ctx,
		p.exchange,
		p.routingKey,
		false, // mandatory (considere true se quiser detectar unroutable)
		false, // immediate (normalmente false mesmo)
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // opcional
		},
	)
}
