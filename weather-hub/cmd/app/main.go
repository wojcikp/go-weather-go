package main

import (
	"log"

	"github.com/lpernett/godotenv"
	"github.com/wojcikp/go-weather-go/weather-hub/internal/app"
	rabbitmqclient "github.com/wojcikp/go-weather-go/weather-hub/internal/rabbitmq_client"
)

const queueName = "queue1"

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file, err: %v", err)
	}
	weatherFeed := make(chan []byte)
	r := rabbitmqclient.GetRabbitClient(queueName, weatherFeed)
	app := app.NewApp(r)
	app.Run()
}
