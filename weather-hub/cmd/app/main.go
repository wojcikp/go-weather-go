package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/wojcikp/go-weather-go/weather-hub/internal/app"
	chclient "github.com/wojcikp/go-weather-go/weather-hub/internal/ch_client"
	rabbitmqclient "github.com/wojcikp/go-weather-go/weather-hub/internal/rabbitmq_client"
	scorereader "github.com/wojcikp/go-weather-go/weather-hub/internal/score_reader"
	weatherfeedconsumer "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_feed_consumer"
	weatherfeedreceiver "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_feed_receiver"
	webserver "github.com/wojcikp/go-weather-go/weather-hub/internal/web_server"
)

var rabbitUser, rabbitPass, rabbitHost, rabbitPort, rabbitQueue, clickhouseDb, clickhouseTable string

func main() {
	app, err := initializeApp()
	if err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
	app.Run()
}

func initializeApp() (*app.App, error) {
	// err := godotenv.Load("../../.env")
	// if err != nil {
	// 	return &app.App{}, fmt.Errorf("error loading .env file, err: %w", err)
	// }
	err := setEnvs()
	if err != nil {
		return &app.App{}, fmt.Errorf("setting env variables error: %w", err)
	}
	weatherFeed := make(chan []byte)
	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUser, rabbitPass, rabbitHost, rabbitPort)
	rabbit := rabbitmqclient.NewRabbitClient(rabbitQueue, rabbitUrl, weatherFeed)
	clickhouse := chclient.NewClickhouseClient(clickhouseDb, clickhouseTable)
	receiver := weatherfeedreceiver.NewFeedReceiver(rabbit)
	consumer := weatherfeedconsumer.NewWeatherFeedConsumer(weatherFeed)
	reader := scorereader.NewConsoleScoreReader()
	server := webserver.NewScoresServer()
	return app.NewApp(clickhouse, receiver, consumer, reader, server), nil
}

func setEnvs() error {
	rabbitUser = os.Getenv("RABBITMQ_DEFAULT_USER")
	if rabbitUser == "" {
		return errors.New("env 'RABBITMQ_DEFAULT_USER' was empty")
	}
	rabbitPass = os.Getenv("RABBITMQ_DEFAULT_PASS")
	if rabbitPass == "" {
		return errors.New("env 'RABBITMQ_DEFAULT_PASS' was empty")
	}
	rabbitHost = os.Getenv("RABBITMQ_HOST")
	if rabbitHost == "" {
		return errors.New("env 'RABBITMQ_HOST' was empty")
	}
	rabbitPort = os.Getenv("RABBITMQ_PORT")
	if rabbitPort == "" {
		return errors.New("env 'RABBITMQ_PORT' was empty")
	}
	rabbitQueue = os.Getenv("RABBITMQ_QUEUE")
	if rabbitQueue == "" {
		return errors.New("env 'RABBITMQ_QUEUE' was empty")
	}
	clickhouseDb = os.Getenv("CLICKHOUSE_DB")
	if clickhouseDb == "" {
		return errors.New("env 'CLICKHOUSE_DB' was empty")
	}
	clickhouseTable = os.Getenv("CLICKHOUSE_TABLE")
	if clickhouseTable == "" {
		return errors.New("env 'CLICKHOUSE_TABLE' was empty")
	}
	return nil
}
