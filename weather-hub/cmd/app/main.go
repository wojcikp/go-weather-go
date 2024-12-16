package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
	"github.com/wojcikp/go-weather-go/weather-hub/internal"
	"github.com/wojcikp/go-weather-go/weather-hub/internal/app"
	chclient "github.com/wojcikp/go-weather-go/weather-hub/internal/ch_client"
	rabbitmqclient "github.com/wojcikp/go-weather-go/weather-hub/internal/rabbitmq_client"
	scorereader "github.com/wojcikp/go-weather-go/weather-hub/internal/score_reader"
	weatherfeedconsumer "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_feed_consumer"
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
	prod, err := strconv.ParseBool(os.Getenv("PRODUCTION"))
	if err != nil {
		log.Print("os env PRODUCTION not found. running local development mode.")
		prod = false
	}
	if !prod {
		if err := godotenv.Load("../../.env"); err != nil {
			return &app.App{}, fmt.Errorf("error loading .env file, err: %w", err)
		}
	}

	if err = setEnvs(); err != nil {
		return &app.App{}, fmt.Errorf("setting env variables error: %w", err)
	}
	weatherFeed := make(chan internal.FeedStream)
	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUser, rabbitPass, rabbitHost, rabbitPort)
	rabbit, err := rabbitmqclient.NewRabbitClient(rabbitQueue, rabbitUrl, weatherFeed)
	if err != nil {
		return &app.App{}, fmt.Errorf("creating rabbit client error: %w", err)
	}
	clickhouse := chclient.NewClickhouseClient(clickhouseDb, clickhouseTable)
	consumer := weatherfeedconsumer.NewWeatherFeedConsumer(weatherFeed)
	reader := scorereader.NewConsoleScoreReader()
	server := webserver.NewScoresServer()
	return app.NewApp(clickhouse, rabbit, consumer, reader, server), nil
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
