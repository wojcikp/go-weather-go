package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lpernett/godotenv"
	"github.com/wojcikp/go-weather-go/weather-feed/config"
	"github.com/wojcikp/go-weather-go/weather-feed/internal"
	apiclient "github.com/wojcikp/go-weather-go/weather-feed/internal/api_client"
	"github.com/wojcikp/go-weather-go/weather-feed/internal/app"
	citiesreader "github.com/wojcikp/go-weather-go/weather-feed/internal/cities_reader"
	rabbitmqpublisher "github.com/wojcikp/go-weather-go/weather-feed/internal/rabbitmq_publisher"
	weatherdataworkers "github.com/wojcikp/go-weather-go/weather-feed/internal/weather_data_workers"
)

func main() {
	app, err := initializeApp()
	if err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
	app.Run()
}

func initializeApp() (*app.App, error) {
	if err := godotenv.Load("../../.env"); err != nil {
		return &app.App{}, fmt.Errorf("error loading .env file, err: %w", err)
	}

	config, err := config.GetConfig()
	if err != nil {
		return &app.App{}, nil
	}

	var reader citiesreader.ICityReader
	if config.MockCityInput {
		reader = citiesreader.NewReaderMock()
	} else {
		reader = citiesreader.NewReader()
	}

	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT"))
	cityData := make(chan internal.CityWeatherData)
	apiClient := apiclient.NewApiClient(config.BaseUrl, config.LookBackwardInMonths)
	rabbitClient := rabbitmqpublisher.NewRabbitPublisher(config.QueueName, rabbitUrl)
	producer := weatherdataworkers.NewApiDataProducer(*apiClient, cityData)
	consumer := weatherdataworkers.NewWeatherDataConsumer(cityData)

	return app.NewApp(config, reader, rabbitClient, producer, consumer), nil
}
