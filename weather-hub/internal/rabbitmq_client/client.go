package rabbitmqclient

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	queue       string
	weatherFeed chan []byte
}

func GetRabbitClient(queue string, weatherFeed chan []byte) *RabbitClient {
	return &RabbitClient{queue, weatherFeed}
}

func (c RabbitClient) ReceiveMessages() {
	user := os.Getenv("RABBITMQ_DEFAULT_USER")
	pass := os.Getenv("RABBITMQ_DEFAULT_PASS")
	url := fmt.Sprintf("amqp://%s:%s@localhost:5672/", user, pass)
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

	var forever chan struct{}

	go func() {
		for data := range msgs {
			c.weatherFeed <- data.Body
			log.Printf("Received a message: %s", data.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
