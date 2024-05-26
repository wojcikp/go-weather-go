package main

import (
	"log"

	"github.com/lpernett/godotenv"
	"github.com/wojcikp/go-weather-go/weather-feed/config"
	"github.com/wojcikp/go-weather-go/weather-feed/internal"
	apiclient "github.com/wojcikp/go-weather-go/weather-feed/internal/api_client"
	"github.com/wojcikp/go-weather-go/weather-feed/internal/app"
	citiesreader "github.com/wojcikp/go-weather-go/weather-feed/internal/cities_reader"
	rabbitmqclient "github.com/wojcikp/go-weather-go/weather-feed/internal/rabbitmq_publisher"
	weatherdataworkers "github.com/wojcikp/go-weather-go/weather-feed/internal/weather_data_workers"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file, err: %v", err)
	}

	cityData := make(chan internal.CityWeatherData)
	config := config.GetConfig()

	var reader citiesreader.ICityReader
	if config.MockCityInput {
		reader = citiesreader.NewReaderMock()
	} else {
		reader = citiesreader.NewReader()
	}

	apiClient := apiclient.NewApiClient(config.BaseUrl, config.LookBackwardInMonths)
	rabbitClient := rabbitmqclient.NewRabbitClient("queue1")
	producer := weatherdataworkers.NewApiDataProducer(*apiClient, cityData)
	consumer := weatherdataworkers.NewWeatherDataConsumer(cityData)

	app := app.NewApp(config, reader, rabbitClient, producer, consumer)
	app.Run()
}
