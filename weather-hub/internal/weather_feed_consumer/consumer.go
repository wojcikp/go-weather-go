package weatherfeedconsumer

import (
	"encoding/json"
	"log"

	"github.com/wojcikp/go-weather-go/weather-hub/internal"
)

type IWeatherFeedConsumer interface {
	ProcessWeatherFeed(code *int, data internal.CityWeatherData)
}

type Consumer struct {
	weatherFeed chan []byte
}

func NewWeatherFeedConsumer(weatherFeed chan []byte) *Consumer {
	return &Consumer{weatherFeed}
}

func (c Consumer) Work(code *int, done chan struct{}, wfc IWeatherFeedConsumer) {
	for msg := range c.weatherFeed {
		data := internal.CityWeatherData{}
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Printf("ERROR: Unmarshalling data failed. Feed has not beed saved to db. error: %v", err)
		}
		wfc.ProcessWeatherFeed(code, data)
		done <- struct{}{}
	}
}
