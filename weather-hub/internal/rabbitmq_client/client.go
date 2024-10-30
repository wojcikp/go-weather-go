package rabbitmqclient

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	queue       string
	url         string
	weatherFeed chan []byte
}

func NewRabbitClient(queue, url string, weatherFeed chan []byte) *RabbitClient {
	return &RabbitClient{queue, url, weatherFeed}
}

	conn, err := amqp.Dial(r.url)
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
		c.queue, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue, err: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer, err: %v", err)
	}

	forever := make(chan struct{})

	go func() {
		for data := range msgs {
			c.weatherFeed <- data.Body
			*feedCounter++
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
