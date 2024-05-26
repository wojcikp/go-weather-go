package app

import "log"

type App struct{}

func NewApp() *App {
	return &App{}
}

func (app App) Run() {
	log.Print("weather hub app run")
}
