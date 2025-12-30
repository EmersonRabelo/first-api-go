package config

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConnection(amqpURL string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewRabbitMQChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func InitBroker() (*amqp.Connection, *amqp.Channel) {
	params := AppSetting.GetBroker()

	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		params.User, params.Password, params.Host, params.Port,
	)

	conn, err := NewRabbitMQConnection(url)
	if err != nil {
		log.Fatalf("Erro ao conectar no RabbitMQ: %v", err)
	}

	ch, err := NewRabbitMQChannel(conn)
	if err != nil {
		_ = conn.Close()
		log.Fatalf("Erro ao abrir channel no RabbitMQ: %v", err)
	}

	return conn, ch
}
