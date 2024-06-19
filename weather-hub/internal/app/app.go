package app

import (
	"log"
	"sync"

	chclient "github.com/wojcikp/go-weather-go/weather-hub/internal/ch_client"
	rabbitmqclient "github.com/wojcikp/go-weather-go/weather-hub/internal/rabbitmq_client"
	weatherfeedconsumer "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_feed_consumer"
	weatherscores "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_scores"
)

type App struct {
	rabbitClient     *rabbitmqclient.RabbitClient
	clickhouseClient *chclient.ClickhouseClient
	feedConsumer     *weatherfeedconsumer.Consumer
}

func NewApp(
	rabbitClient *rabbitmqclient.RabbitClient,
	clickhouseClient *chclient.ClickhouseClient,
	feedConsumer *weatherfeedconsumer.Consumer,
) *App {
	return &App{rabbitClient, clickhouseClient, feedConsumer}
}

func (app App) Run() {
	log.Print("weather hub app run")
	for _, score := range weatherscores.GetScoresList[string]() {
		log.Print("Id: ", score.GetId())
		log.Print("Name: ", score.GetName())
		value, err := score.GetScore(app.clickhouseClient)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Value: ", value)
	}
	wg := &sync.WaitGroup{}
	app.clickhouseClient.CreateWeatherTable()
	wg.Add(1)
	go app.rabbitClient.ReceiveMessages(wg)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go app.feedConsumer.Work(wg, app.clickhouseClient)
	}
	wg.Wait()
}
