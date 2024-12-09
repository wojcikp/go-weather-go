package app

import (
	"context"
	"log"
	"time"

	"github.com/wojcikp/go-weather-go/weather-feed/config"
	"github.com/wojcikp/go-weather-go/weather-feed/internal"
	citiesreader "github.com/wojcikp/go-weather-go/weather-feed/internal/cities_reader"
	rabbitmqpublisher "github.com/wojcikp/go-weather-go/weather-feed/internal/rabbitmq_publisher"
	weatherdataworkers "github.com/wojcikp/go-weather-go/weather-feed/internal/weather_data_workers"
	"golang.org/x/sync/semaphore"
)

type App struct {
	config          config.Configuration
	reader          citiesreader.ICityReader
	rabbitPublisher *rabbitmqpublisher.RabbitPublisher
	producer        *weatherdataworkers.ApiDataProducer
	consumer        *weatherdataworkers.Consumer
}

func NewApp(
	config config.Configuration,
	reader citiesreader.ICityReader,
	rabbitPublisher *rabbitmqpublisher.RabbitPublisher,
	producer *weatherdataworkers.ApiDataProducer,
	consumer *weatherdataworkers.Consumer,
) *App {
	return &App{config, reader, rabbitPublisher, producer, consumer}
}

func (app App) Run() {
	cities, err := citiesreader.GetCitiesInput(app.reader)
	if err != nil {
		log.Fatal(err)
	}

	const publishFeedInterval = 20 * time.Second
	go runProducers(app, cities, publishFeedInterval)

	done := make(chan struct{})
	for i := 0; i < app.config.ConsumerCount; i++ {
		go app.consumer.Work(done, app.rabbitPublisher)
	}

	go func() {
		for {
			<-done
			for i := 0; i < len(cities)-1; i++ {
				<-done
			}
			log.Print("Published weather feed.")
		}
	}()

	forever := make(chan struct{})
	<-forever
}

func runProducers(app App, cities []internal.BaseCityInfo, publishFeedInterval time.Duration) {
	ctx := context.Background()
	sem := semaphore.NewWeighted(5)
	for {
		for _, city := range cities {
			go func(city internal.BaseCityInfo) {
				sem.Acquire(ctx, 1)
				app.producer.Work(ctx, city, sem)
			}(city)
		}
		<-time.After(publishFeedInterval)
	}
}
