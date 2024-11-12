package rabbitmqclient

import (
	"fmt"
	"log"
	"sync"

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

func (r RabbitClient) ReceiveMessages(feedCounter *int, mu *sync.Mutex) error {
	conn, err := amqp.Dial(r.url)
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
		return fmt.Errorf("failed to register a consumer, err: %w", err)
	}

	forever := make(chan struct{})

	go func() {
		for data := range msgs {
			r.weatherFeed <- data.Body
			mu.Lock()
			*feedCounter++
			mu.Unlock()
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}
