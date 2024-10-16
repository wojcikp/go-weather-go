package weatherfeedconsumer

import (
	"encoding/json"
	"log"

	"github.com/wojcikp/go-weather-go/weather-hub/internal"
)

type IWeatherFeedConsumer interface {
	ProcessWeatherFeed(data internal.CityWeatherData)
}

type Consumer struct {
	weatherFeed chan []byte
}

func NewWeatherFeedConsumer(weatherFeed chan []byte) *Consumer {
	return &Consumer{weatherFeed}
}

func (c Consumer) Work(done chan struct{}, wfc IWeatherFeedConsumer) {
	for msg := range c.weatherFeed {
		data := internal.CityWeatherData{}
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Fatal(err)
		}
		log.Printf("Processing data feed for city: %s", data.Name)
		wfc.ProcessWeatherFeed(data)
		done <- struct{}{}
	}
}
