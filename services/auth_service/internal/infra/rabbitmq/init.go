package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		"auth_exchange",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		_ = ch.Close()
		_ = conn.Close()

		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil

}

func (r *RabbitMQ) Publish(event string, body []byte) error {
	return r.channel.Publish(
		"auth_exchange",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		_ = r.channel.Close()
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
}
