package main

import (
	"log"

	"github.com/lpernett/godotenv"
	"github.com/wojcikp/go-weather-go/weather-hub/internal/app"
	chclient "github.com/wojcikp/go-weather-go/weather-hub/internal/ch_client"
	rabbitmqclient "github.com/wojcikp/go-weather-go/weather-hub/internal/rabbitmq_client"
	scorereader "github.com/wojcikp/go-weather-go/weather-hub/internal/score_reader"
	weatherfeedconsumer "github.com/wojcikp/go-weather-go/weather-hub/internal/weather_feed_consumer"
)

const queueName = "queue1"

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file, err: %v", err)
	}
	weatherFeed := make(chan []byte)
	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT"))
	rabbit := rabbitmqclient.NewRabbitClient(os.Getenv("RABBITMQ_QUEUE"), rabbitUrl, weatherFeed)
	clickhouse := chclient.NewClickhouseClient()
	consumer := weatherfeedconsumer.NewWeatherFeedConsumer(weatherFeed)
	reader := scorereader.NewConsoleScoreReader()
	app := app.NewApp(rabbit, clickhouse, consumer, reader)
	app.Run()
}
