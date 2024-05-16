package main

import (
	"github.com/wojcikp/go-weather-go/weather-feed/config"
	apiclient "github.com/wojcikp/go-weather-go/weather-feed/internal/api_client"
	"github.com/wojcikp/go-weather-go/weather-feed/internal/app"
	citiesreader "github.com/wojcikp/go-weather-go/weather-feed/internal/cities_reader"
)

func main() {
	config := config.GetConfig()
	reader := citiesreader.NewReaderMock()
	client := apiclient.NewApiClient(config.BaseUrl, config.LookBackwardInMonths)

	app := app.NewApp(config, reader, client)
	app.Run()
}
