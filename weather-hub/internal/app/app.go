package app

import (
	"bytes"
	"log"

	chclient "github.com/wojcikp/go-weather-go/weather-hub/internal/ch_client"
	scorereader "github.com/wojcikp/go-weather-go/weather-hub/internal/score_reader"
	weatherfeedconsumer "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_feed_consumer"
	weatherfeedreceiver "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_feed_receiver"
	weatherscores "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_scores"
)

type App struct {
	clickhouseClient *chclient.ClickhouseClient
	feedReceiver     *weatherfeedreceiver.FeedReceiver
	feedConsumer     *weatherfeedconsumer.Consumer
	reader           scorereader.IScoreReader
}

func NewApp(
	clickhouseClient *chclient.ClickhouseClient,
	feedReceiver *weatherfeedreceiver.FeedReceiver,
	feedConsumer *weatherfeedconsumer.Consumer,
	reader scorereader.IScoreReader,
) *App {
	return &App{clickhouseClient, feedReceiver, feedConsumer, reader}
}

func (app App) Run() {
	log.Print("weather hub app run")
	done := make(chan struct{})
	feedCounter := 0
	go app.clickhouseClient.CreateWeatherTable()
	go app.feedReceiver.HandleReceiveMessages(&feedCounter)
	for i := 0; i < 5; i++ {
		go app.feedConsumer.Work(done, app.clickhouseClient)
	}
	go processScores(app.reader, app.clickhouseClient, done, &feedCounter)
	forever := make(chan struct{})
	<-forever
}

func processScores(
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
		reader.ReadScores(&scoresInfo)
		*feedCounter = 0
	}
}
