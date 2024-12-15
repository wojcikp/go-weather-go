package app

import (
	"bytes"
	"fmt"
	"log"

	"github.com/wojcikp/go-weather-go/weather-hub/internal"
	chclient "github.com/wojcikp/go-weather-go/weather-hub/internal/ch_client"
	scorereader "github.com/wojcikp/go-weather-go/weather-hub/internal/score_reader"
	weatherfeedconsumer "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_feed_consumer"
	weatherscores "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_scores"
	webserver "github.com/wojcikp/go-weather-go/weather-hub/internal/web_server"
)

type App struct {
	clickhouseClient *chclient.ClickhouseClient
	feedReceiver     internal.IFeedReceiver
	feedConsumer     *weatherfeedconsumer.Consumer
	reader           scorereader.IScoreReader
	server           *webserver.ScoresServer
}

func NewApp(
	clickhouseClient *chclient.ClickhouseClient,
	feedReceiver internal.IFeedReceiver,
	feedConsumer *weatherfeedconsumer.Consumer,
	reader scorereader.IScoreReader,
	server *webserver.ScoresServer,
) *App {
	return &App{clickhouseClient, feedReceiver, feedConsumer, reader, server}
}

func (app App) Run() {
	log.Print("Weather hub app run")
	log.Print("Waiting for messages. To exit press CTRL+C")
	done := make(chan struct{})
	go app.server.RunWeatherScoresServer()
	for i := 0; i < 10; i++ {
		go app.feedReceiver.ReceiveMessages()
	}
	for i := 0; i < 10; i++ {
		go app.feedConsumer.Work(done, app.clickhouseClient)
	}
	go processScores(app.server, app.reader, app.clickhouseClient, done)
	forever := make(chan struct{})
	<-forever
}

func processScores(
	server *webserver.ScoresServer,
	reader scorereader.IScoreReader,
	clickhouseClient *chclient.ClickhouseClient,
	done chan struct{},
) {
	stringScores := weatherscores.GetScoresList[string]()
	floatScores := weatherscores.GetScoresList[float64]()
	const feedLength = 172
	for {
		<-done
		for i := 0; i < feedLength-1; i++ {
			<-done
		}
		log.Print("Actual Scores:")
		stringScoresInfo, stringErrors := weatherscores.GetScoresInfo(stringScores, clickhouseClient)
		floatScoresInfo, floatErrors := weatherscores.GetScoresInfo(floatScores, clickhouseClient)
		responseScoresInfo := make([]internal.ScoreInfo, 0, len(stringScoresInfo)+len(floatScoresInfo))
		errors := []error{}
		responseScoresInfo = append(responseScoresInfo, stringScoresInfo...)
		responseScoresInfo = append(responseScoresInfo, floatScoresInfo...)
		errors = append(errors, stringErrors...)
		errors = append(errors, floatErrors...)
		if len(errors) > 0 {
			log.Print("ERROR: Some errors occured while reading scores info: ")
			for _, err := range errors {
				log.Print(err)
			}
		}
		publishScores(server, reader, responseScoresInfo)
	}
}

func publishScores(
	server *webserver.ScoresServer,
	reader scorereader.IScoreReader,
	scoresInfo []internal.ScoreInfo,
) {
	var scoresInfoMessage bytes.Buffer
	for _, score := range scoresInfo {
		scoresInfoMessage.WriteString(fmt.Sprintf("Id: %d\n", score.Id))
		scoresInfoMessage.WriteString(fmt.Sprintf("Name: %s\n", score.Name))
		scoresInfoMessage.WriteString(fmt.Sprintf("Value: %s\n", score.Value))
		scoresInfoMessage.WriteString("-----------------------------\n")
	}
	server.SetScoresInfo(scoresInfo)
	reader.ReadScores(&scoresInfoMessage)
}
