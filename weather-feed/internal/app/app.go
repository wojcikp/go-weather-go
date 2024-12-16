package app

import (
	"context"
	"log"
	"os/signal"
	"syscall"
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cities, err := citiesreader.GetCitiesInput(app.reader)
	if err != nil {
		log.Fatal(err)
	}

	const publishFeedInterval = 3 * time.Minute
	go runProducers(ctx, app, cities, publishFeedInterval)

	done := make(chan struct{})
	for i := 0; i < app.config.ConsumerCount; i++ {
		go app.consumer.Work(done, app.rabbitPublisher)
	}

	go func() {
		for {
			for i := 0; i < len(cities); i++ {
				<-done
			}
			log.Print("Published weather feed.")
		}
	}()

	<-ctx.Done()
	if err := app.rabbitPublisher.Close(); err != nil {
		log.Print("rabbit publisher close error: ", err)
	}
	log.Print("Weather feed shutdown signal received, job finished.")
}

func runProducers(ctx context.Context, app App, cities []internal.BaseCityInfo, publishFeedInterval time.Duration) {
	sem := semaphore.NewWeighted(5)
	for {
		for _, city := range cities {
			go func(city internal.BaseCityInfo) {
				sem.Acquire(ctx, 1)
				app.producer.Work(ctx, city)
				sem.Release(1)
			}(city)
		}
		<-time.After(publishFeedInterval)
	}
}
