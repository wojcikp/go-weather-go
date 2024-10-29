package rabbitmqpublisher

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	queue string
}

func NewRabbitPublisher(queue string) *RabbitPublisher {
	return &RabbitPublisher{queue}
}

func (r RabbitPublisher) ProcessWeatherData(data []byte) error {
	if err := r.putMsgOnQueue(data); err != nil {
		return err
	}
	return nil
}

func (r RabbitPublisher) putMsgOnQueue(msg []byte) error {
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ, err: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel, err: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		r.queue, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue, err: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
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
