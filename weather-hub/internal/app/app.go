package app

import (
	"bytes"
	"log"

	"github.com/wojcikp/go-weather-go/weather-hub/internal"
	chclient "github.com/wojcikp/go-weather-go/weather-hub/internal/ch_client"
	scorereader "github.com/wojcikp/go-weather-go/weather-hub/internal/score_reader"
	weatherfeedconsumer "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_feed_consumer"
	weatherfeedreceiver "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_feed_receiver"
	weatherscores "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_scores"
	webserver "github.com/wojcikp/go-weather-go/weather-hub/internal/web_server"
)

type App struct {
	clickhouseClient *chclient.ClickhouseClient
	feedReceiver     *weatherfeedreceiver.FeedReceiver
	feedConsumer     *weatherfeedconsumer.Consumer
	reader           scorereader.IScoreReader
	server           *webserver.ScoresServer
}

func NewApp(
	clickhouseClient *chclient.ClickhouseClient,
	feedReceiver *weatherfeedreceiver.FeedReceiver,
	feedConsumer *weatherfeedconsumer.Consumer,
	reader scorereader.IScoreReader,
	server *webserver.ScoresServer,
) *App {
	return &App{clickhouseClient, feedReceiver, feedConsumer, reader, server}
}

func (app App) Run() {
	log.Print("weather hub app run")
	done := make(chan struct{})
	feedCounter := 0
	go app.clickhouseClient.CreateWeatherTable()
	for i := 0; i < 1; i++ {
		go app.feedReceiver.HandleReceiveMessages(&feedCounter)
	}
	for i := 0; i < 5; i++ {
		go app.feedConsumer.Work(done, app.clickhouseClient)
	}
	go processScores(app.reader, app.clickhouseClient, done, &feedCounter, app.server)
	go app.server.RunWeatherScoresServer()
	forever := make(chan struct{})
	<-forever
}

func processScores(
	reader scorereader.IScoreReader,
	clickhouseClient *chclient.ClickhouseClient,
	done chan struct{},
	feedCounter *int,
	server *webserver.ScoresServer,
) {
	stringScores := weatherscores.GetScoresList[string]()
	floatScores := weatherscores.GetScoresList[float64]()
	var scoresInfo bytes.Buffer
	errors := []error{}
	for {
		<-done
		for i := 0; i < *feedCounter-1; i++ {
			<-done
		}
		stringScores, stringErrors := weatherscores.GetScoresInfo(stringScores, &scoresInfo, clickhouseClient)
		floatScores, floatErrors := weatherscores.GetScoresInfo(floatScores, &scoresInfo, clickhouseClient)
		errors = append(errors, stringErrors...)
		errors = append(errors, floatErrors...)
		responseInfo := []internal.ScoreInfo{}
		responseInfo = append(responseInfo, stringScores...)
		responseInfo = append(responseInfo, floatScores...)
		log.Print("responseInfo", responseInfo)
		server.SetScoresInfo(responseInfo)
		if len(errors) > 0 {
			log.Print("ERROR: Some errors occured while reading scores info:")
			for _, err := range errors {
				log.Print(err)
			}
		}
		reader.ReadScores(&scoresInfo)
		*feedCounter = 0
	}
}
