package main

import (
	"log"

	"github.com/lpernett/godotenv"
	"github.com/wojcikp/go-weather-go/weather-feed/config"
	"github.com/wojcikp/go-weather-go/weather-feed/internal"
	apiclient "github.com/wojcikp/go-weather-go/weather-feed/internal/api_client"
	apiworker "github.com/wojcikp/go-weather-go/weather-feed/internal/api_worker"
	"github.com/wojcikp/go-weather-go/weather-feed/internal/app"
	citiesreader "github.com/wojcikp/go-weather-go/weather-feed/internal/cities_reader"
	rabbitmqclient "github.com/wojcikp/go-weather-go/weather-feed/internal/rabbitmq_client"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file, err: %v", err)
	}
	cityData := make(chan internal.CityWeatherData)

	config := config.GetConfig()
	reader := citiesreader.NewReaderMock()
	apiClient := apiclient.NewApiClient(config.BaseUrl, config.LookBackwardInMonths)
	rabbitClient := rabbitmqclient.NewRabbitClient("queue1")
	consumer := weatherdataworkers.NewWeatherDataConsumer(cityData)

	app := app.NewApp(config, reader, worker)
	app.Run()
}
