package main

import "github.com/wojcikp/go-weather-go/weather-hub/internal/app"

func main() {
	app := app.NewApp()
	app.Run()
}
