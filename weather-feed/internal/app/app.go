package app

import (
	"context"
	"log"
	"sync"

	"github.com/wojcikp/go-weather-go/weather-feed/config"
	"github.com/wojcikp/go-weather-go/weather-feed/internal"
	apiworker "github.com/wojcikp/go-weather-go/weather-feed/internal/api_worker"
	citiesreader "github.com/wojcikp/go-weather-go/weather-feed/internal/cities_reader"
	"golang.org/x/sync/semaphore"
)

type IApp interface {
	Run()
}

type App struct {
	config    config.Configuration
	reader    citiesreader.ICityReader
	apiWorker *apiworker.ApiDataWorker
}

func NewApp(
	config config.Configuration,
	reader citiesreader.ICityReader,
	apiWorker *apiworker.ApiDataWorker,
) *App {
	return &App{config, reader, apiWorker}
}

func (app App) Run() {
	ctx := context.Background()
	cities, err := citiesreader.GetCitiesInput(app.reader)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(cities)

	wg := &sync.WaitGroup{}
	sem := semaphore.NewWeighted(5)

	for _, city := range cities {
		wg.Add(1)
		go func(city internal.BaseCityInfo) {
			sem.Acquire(ctx, 1)
			app.apiWorker.Work(ctx, city, wg, sem)
		}(city)
	}
	log.Print(<-app.apiWorker.CityData)
	log.Print(<-app.apiWorker.CityData)
	log.Print(<-app.apiWorker.CityData)
	log.Print(<-app.apiWorker.CityData)
	log.Print(<-app.apiWorker.CityData)
	wg.Wait()
	close(app.apiWorker.CityData)
	log.Print("DONE")
}
