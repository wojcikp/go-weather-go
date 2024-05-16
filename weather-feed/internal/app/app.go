package app

import (
	"context"
	"log"

	"github.com/wojcikp/go-weather-go/weather-feed/config"
	apiclient "github.com/wojcikp/go-weather-go/weather-feed/internal/api_client"
	citiesreader "github.com/wojcikp/go-weather-go/weather-feed/internal/cities_reader"
)

type IApp interface {
	Run()
}

type App struct {
	config    config.Configuration
	reader    citiesreader.ICityReader
	apiClient *apiclient.WeatherApiClient
}

func NewApp(config config.Configuration, reader citiesreader.ICityReader, apiClient *apiclient.WeatherApiClient) *App {
	return &App{config, reader, apiClient}
}

func (app App) Run() {
	ctx := context.Background()
	cities, err := citiesreader.GetCitiesInput(app.reader)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(cities)
	data, err := app.apiClient.FetchData(ctx, cities[0])
	if err != nil {
		log.Fatal(err)
	}
	log.Print(data)
}
