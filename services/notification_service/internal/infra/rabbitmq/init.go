package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/streadway/amqp"
)

type Consumer struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	queue   string
}

func NewConsumer(url, queue, exchange string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect consuConsumer: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		_ = ch.Close()

		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		conn.Close()
		ch.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	q, err := ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	if err := ch.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	); err != nil {
		conn.Close()
		ch.Close()
		return nil, fmt.Errorf("queue bind failed: %w", err)
	}

	return &Consumer{
		Conn:    conn,
		Channel: ch,
		queue:   queue,
	}, nil
}

func (r *Consumer) Start(handler func(event map[string]any) error) error {
	msgs, err := r.Channel.Consume(
		r.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		for d := range msgs {
			var event map[string]any
			if err := json.Unmarshal(d.Body, &event); err != nil {
				slog.Error("Invalid JSON: %v", "error", err)
				_ = d.Nack(false, false)
				continue
			}

			if err := handler(event); err != nil {
				slog.Error("Handler error: %w", "error", err)
				_ = d.Nack(false, true)
				continue
			}

			_ = d.Ack(false)
		}
	}()

	return nil
}
