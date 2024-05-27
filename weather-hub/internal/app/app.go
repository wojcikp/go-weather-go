package app

import (
	"log"

	rabbitmqclient "github.com/wojcikp/go-weather-go/weather-hub/internal/rabbitmq_client"
)

type App struct {
	rabbitClient *rabbitmqclient.RabbitClient
}

func NewApp(rabbitClient *rabbitmqclient.RabbitClient) *App {
	return &App{rabbitClient}
}

func (app App) Run() {
	log.Print("weather hub app run")
	app.rabbitClient.ReceiveMessages()
}
