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

func (r RabbitPublisher) ProcessWeatherData(data []byte) {
	r.putMsgOnQueue(data)
}

func (r RabbitPublisher) putMsgOnQueue(msg []byte) {
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASS")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ, err: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel, err: %v", err)
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
		log.Fatalf("Failed to declare a queue, err: %v", err)
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
		log.Fatalf("Failed to publish a message, err: %v", err)
	}
}
