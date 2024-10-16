package app

import (
	"bytes"
	"log"

	chclient "github.com/wojcikp/go-weather-go/weather-hub/internal/ch_client"
	rabbitmqclient "github.com/wojcikp/go-weather-go/weather-hub/internal/rabbitmq_client"
	scorereader "github.com/wojcikp/go-weather-go/weather-hub/internal/score_reader"
	weatherfeedconsumer "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_feed_consumer"
	weatherscores "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_scores"
)

type App struct {
	rabbitClient     *rabbitmqclient.RabbitClient
	clickhouseClient *chclient.ClickhouseClient
	feedConsumer     *weatherfeedconsumer.Consumer
	reader           scorereader.IScoreReader
}

func NewApp(
	rabbitClient *rabbitmqclient.RabbitClient,
	clickhouseClient *chclient.ClickhouseClient,
	feedConsumer *weatherfeedconsumer.Consumer,
	reader scorereader.IScoreReader,
) *App {
	return &App{rabbitClient, clickhouseClient, feedConsumer, reader}
}

func (app App) Run() {
	log.Print("weather hub app run")
	done := make(chan struct{})
	feedCounter := 0
	go app.clickhouseClient.CreateWeatherTable()
	go app.rabbitClient.ReceiveMessages(&feedCounter)
	for i := 0; i < 5; i++ {
		go app.feedConsumer.Work(done, app.clickhouseClient)
	}
	go readScores(app.reader, app.clickhouseClient, done, &feedCounter)
	forever := make(chan struct{})
	<-forever
}

func readScores(
	reader scorereader.IScoreReader,
	clickhouseClient *chclient.ClickhouseClient,
	done chan struct{},
	feedCounter *int,
) {
	stringScores := weatherscores.GetScoresList[string]()
	floatScores := weatherscores.GetScoresList[float64]()
	var scoresInfo bytes.Buffer
	for {
		<-done
		for i := 0; i < *feedCounter-1; i++ {
			<-done
		}
		weatherscores.GetScoresInfo(stringScores, &scoresInfo, clickhouseClient)
		weatherscores.GetScoresInfo(floatScores, &scoresInfo, clickhouseClient)
		scorereader.ReadScores(reader, &scoresInfo)
		*feedCounter = 0
	}
}
