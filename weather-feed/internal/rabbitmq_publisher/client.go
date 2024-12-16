package rabbitmqpublisher

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	queue   amqp.Queue
	channel *amqp.Channel
	conn    *amqp.Connection
	url     string
}

func NewRabbitPublisher(queueName, url string) (*RabbitPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ, err: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel, err: %w", err)
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue, err: %w", err)
	}
	return &RabbitPublisher{q, ch, conn, url}, nil
}

func (r RabbitPublisher) ProcessWeatherData(msg []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.channel.PublishWithContext(
		ctx,
		"",           // exchange
		r.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message, err: %w", err)
	}

	return nil
}

func (r *RabbitPublisher) Close() error {
	if err := r.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel, err: %w", err)
	}

	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection, err: %w", err)
	}

	return nil
}
