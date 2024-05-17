package main

import (
	"github.com/wojcikp/go-weather-go/weather-feed/config"
	"github.com/wojcikp/go-weather-go/weather-feed/internal"
	apiclient "github.com/wojcikp/go-weather-go/weather-feed/internal/api_client"
	apiworker "github.com/wojcikp/go-weather-go/weather-feed/internal/api_worker"
	"github.com/wojcikp/go-weather-go/weather-feed/internal/app"
	citiesreader "github.com/wojcikp/go-weather-go/weather-feed/internal/cities_reader"
)

func main() {
	cityData := make(chan internal.CityWeatherData)

	config := config.GetConfig()
	reader := citiesreader.NewReaderMock()
	client := apiclient.NewApiClient(config.BaseUrl, config.LookBackwardInMonths)
	worker := apiworker.NewApiDataWorker(*client, cityData)

	app := app.NewApp(config, reader, worker)
	app.Run()
}
