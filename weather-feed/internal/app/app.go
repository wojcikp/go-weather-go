package app

import (
	"context"
	"log"
	"sync"

	"github.com/wojcikp/go-weather-go/weather-feed/config"
	"github.com/wojcikp/go-weather-go/weather-feed/internal"
	citiesreader "github.com/wojcikp/go-weather-go/weather-feed/internal/cities_reader"
	rabbitmqclient "github.com/wojcikp/go-weather-go/weather-feed/internal/rabbitmq_publisher"
	weatherdataworkers "github.com/wojcikp/go-weather-go/weather-feed/internal/weather_data_workers"
	"golang.org/x/sync/semaphore"
)

type IApp interface {
	Run()
}

type App struct {
	config       config.Configuration
	reader       citiesreader.ICityReader
	rabbitClient *rabbitmqclient.RabbitClient
	producer     *weatherdataworkers.ApiDataProducer
	consumer     *weatherdataworkers.Consumer
}

func NewApp(
	config config.Configuration,
	reader citiesreader.ICityReader,
	rabbitClient *rabbitmqclient.RabbitClient,
	producer *weatherdataworkers.ApiDataProducer,
	consumer *weatherdataworkers.Consumer,
) *App {
	return &App{config, reader, rabbitClient, producer, consumer}
}

func (app App) Run() {
	ctx := context.Background()
	cities, err := citiesreader.GetCitiesInput(app.reader)
	if err != nil {
		log.Fatal(err)
	}

	wgp := &sync.WaitGroup{}
	sem := semaphore.NewWeighted(5)
	for _, city := range cities {
		wgp.Add(1)
		go func(city internal.BaseCityInfo) {
			sem.Acquire(ctx, 1)
			app.producer.Work(ctx, city, wgp, sem)
		}(city)
	}

	wgc := &sync.WaitGroup{}
	for i := 0; i < app.config.ConsumerCount; i++ {
		wgc.Add(1)
		go app.consumer.Work(wgc, app.rabbitClient)
	}

	wgp.Wait()
	close(app.producer.CityData)
	wgc.Wait()
	log.Print("DONE")
}
