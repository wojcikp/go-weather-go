package app

import (
	"log"
	"os"
	"path"

	"github.com/wojcikp/go-weather-go/weather-feed/config"
	citiesreader "github.com/wojcikp/go-weather-go/weather-feed/internal/cities_reader"
)

type IApp interface {
	Run()
}

type App struct {
	config config.Configuration
	reader *citiesreader.CitiesReader
}

func NewApp(config config.Configuration, reader *citiesreader.CitiesReader) *App {
	return &App{config, reader}
}

func (app App) Run() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	app.reader.FilePath = path.Join(dir, "..", "..", "assets", "pl172.json")
	cities, err := citiesreader.GetCitiesInput(app.reader)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(cities)
}
