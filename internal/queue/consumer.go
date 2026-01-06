package queue

import (
	"errors"
	"fmt"
	"log"

	"github.com/EmersonRabelo/first-api-go/internal/handler"
	service "github.com/EmersonRabelo/first-api-go/internal/service/consumer"
	amqp "github.com/rabbitmq/amqp091-go"
)

var exchangeType string = "topic"

type ReportConsumer struct {
	ch         *amqp.Channel
	exchange   string
	routingKey string
	queueName  string
	handler    DeliveryHandler
}

type DeliveryHandler interface {
	Handler(delivery amqp.Delivery) error
}

// Construtor
func NewReportConsumer(ch *amqp.Channel, exchange, routingKey, queueName string, handler DeliveryHandler) *ReportConsumer {
	return &ReportConsumer{
		ch:         ch,
		exchange:   exchange,
		routingKey: routingKey,
		queueName:  queueName,
		handler:    handler,
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func (consumer *ReportConsumer) Start() error {

	err := consumer.ch.ExchangeDeclare(
		consumer.exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare exchange")

	queue, err := consumer.ch.QueueDeclare(
		consumer.queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare queue")

	err = consumer.ch.QueueBind(
		queue.Name,
		consumer.routingKey,
		consumer.exchange,
		false,
		nil,
	)

	failOnError(err, "Failed to bind queue")

	deliveries, err := consumer.ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to register consumer")

	log.Printf("Consumidor ativo! Binding: [%s]", consumer.routingKey)

	for delivey := range deliveries {
		err := consumer.handler.Handler(delivey)
		if err != nil {
			if errors.Is(err, handler.ErrPermanent) || errors.Is(err, service.ErrInvalidMessage) {
				_ = delivey.Ack(false)
			} else {
				fmt.Println(err)
				_ = delivey.Nack(false, true) // TOMAR CUIDADO!!! CAI EM UM ERRO DE RETRY DA MENSAGENS, CASO OCORRA ALGUM ERRO NO HANDLER, SERIA BOM UM RATE LIMIT
			}
			continue
		}

		_ = delivey.Ack(false)
	}

	return nil

}
