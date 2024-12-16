package rabbitmqclient

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wojcikp/go-weather-go/weather-hub/internal"
)

type RabbitClient struct {
	url         string
	queue       amqp.Queue
	channel     *amqp.Channel
	conn        *amqp.Connection
	weatherFeed chan internal.FeedStream
}

func NewRabbitClient(queueName, url string, weatherFeed chan internal.FeedStream) (*RabbitClient, error) {
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
	return &RabbitClient{url, q, ch, conn, weatherFeed}, nil
}

func (r RabbitClient) ReceiveMessages() {
	msgs, err := r.channel.Consume(
		r.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		r.weatherFeed <- internal.FeedStream{
			Data: []byte{},
			Err:  fmt.Errorf("failed to register a  rabbit consumer, err: %w", err)}
	}

	forever := make(chan struct{})

	go func() {
		for data := range msgs {
			r.weatherFeed <- internal.FeedStream{Data: data.Body, Err: nil}
		}
	}()

	<-forever
}
