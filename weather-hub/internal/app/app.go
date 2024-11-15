package app

import (
	"bytes"
	"fmt"
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
	log.Print("Weather hub app run")
	log.Print("Waiting for messages. To exit press CTRL+C")
	done := make(chan struct{})
	var code int
	go app.server.RunWeatherScoresServer()
	for i := 0; i < 100; i++ {
		go app.feedReceiver.HandleReceiveMessages()
	}
	for i := 0; i < 1; i++ {
		go app.feedConsumer.Work(&code, done, app.clickhouseClient)
	}
	go processScores(app.server, app.reader, app.clickhouseClient, done, &code)
	forever := make(chan struct{})
	<-forever
}

func processScores(
	server *webserver.ScoresServer,
	reader scorereader.IScoreReader,
	clickhouseClient *chclient.ClickhouseClient,
	done chan struct{},
	code *int,
) {
	stringScores := weatherscores.GetScoresList[string]()
	floatScores := weatherscores.GetScoresList[float64]()
	const feedLength = 172
	for {
		<-done
		for i := 0; i < feedLength-1; i++ {
			<-done
		}
		log.Print("Processing scores...")
		responseScoresInfo := []internal.ScoreInfo{}
		errors := []error{}
		stringScoresInfo, stringErrors := weatherscores.GetScoresInfo(stringScores, clickhouseClient)
		floatScoresInfo, floatErrors := weatherscores.GetScoresInfo(floatScores, clickhouseClient)
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
		fmt.Print(*code)
		*code = 0
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
