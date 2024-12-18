package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
	"github.com/wojcikp/go-weather-go/weather-feed/config"
	"github.com/wojcikp/go-weather-go/weather-feed/internal"
	apiclient "github.com/wojcikp/go-weather-go/weather-feed/internal/api_client"
	"github.com/wojcikp/go-weather-go/weather-feed/internal/app"
	citiesreader "github.com/wojcikp/go-weather-go/weather-feed/internal/cities_reader"
	rabbitmqpublisher "github.com/wojcikp/go-weather-go/weather-feed/internal/rabbitmq_publisher"
	weatherdataworkers "github.com/wojcikp/go-weather-go/weather-feed/internal/weather_data_workers"
)

var rabbitUser, rabbitPass, rabbitHost, rabbitPort, rabbitQueue string

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

	err = setEnvs()
	if err != nil {
		return &app.App{}, fmt.Errorf("setting env variables error: %w", err)
	}

	config, err := config.GetConfig()
	if err != nil {
		return &app.App{}, fmt.Errorf("reading config file error: %w", err)
	}

	var reader citiesreader.ICityReader
	if config.MockCityInput {
		reader = citiesreader.NewReaderMock()
	} else {
		reader = citiesreader.NewReader()
	}

	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUser, rabbitPass, rabbitHost, rabbitPort)
	rabbitClient, err := rabbitmqpublisher.NewRabbitPublisher(rabbitQueue, rabbitUrl)

	if err != nil {
		return &app.App{}, fmt.Errorf("creating rabbit publisher error: %w", err)
	}

	cityData := make(chan internal.CityWeatherData)
	apiClient := apiclient.NewApiClient(config.BaseUrl, config.LookBackwardInMonths)
	producer := weatherdataworkers.NewApiDataProducer(*apiClient, cityData)
	consumer := weatherdataworkers.NewWeatherDataConsumer(cityData)

	return app.NewApp(config, reader, rabbitClient, producer, consumer), nil
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
	return nil
}
