package main

import (
	"github.com/wojcikp/go-weather-go/weather-feed/config"
	"github.com/wojcikp/go-weather-go/weather-feed/internal/app"
	citiesreader "github.com/wojcikp/go-weather-go/weather-feed/internal/cities_reader"
)

func main() {
	config := config.GetConfig()
	r := citiesreader.NewReader(config.CitiesJsonPath)

	app := app.NewApp(config, r)
	app.Run()
}
