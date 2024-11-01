package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lpernett/godotenv"
	"github.com/wojcikp/go-weather-go/weather-hub/internal/app"
	chclient "github.com/wojcikp/go-weather-go/weather-hub/internal/ch_client"
	rabbitmqclient "github.com/wojcikp/go-weather-go/weather-hub/internal/rabbitmq_client"
	scorereader "github.com/wojcikp/go-weather-go/weather-hub/internal/score_reader"
	weatherfeedconsumer "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_feed_consumer"
	weatherfeedreceiver "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_feed_receiver"
)

func main() {
	app, err := initializeApp()
	if err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
	app.Run()
}

func initializeApp() (*app.App, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return &app.App{}, fmt.Errorf("error loading .env file, err: %w", err)
	}
	weatherFeed := make(chan []byte)
	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT"))
	rabbit := rabbitmqclient.NewRabbitClient(os.Getenv("RABBITMQ_QUEUE"), rabbitUrl, weatherFeed)
	clickhouse := chclient.NewClickhouseClient()
	receiver := weatherfeedreceiver.NewFeedReceiver(rabbit)
	consumer := weatherfeedconsumer.NewWeatherFeedConsumer(weatherFeed)
	reader := scorereader.NewConsoleScoreReader()
	return app.NewApp(clickhouse, receiver, consumer, reader), nil
}
