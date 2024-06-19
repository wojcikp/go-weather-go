package weatherfeedconsumer

import (
	"encoding/json"
	"log"
	"sync"

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

func (c Consumer) Work(wg *sync.WaitGroup, wfc IWeatherFeedConsumer) {
	defer wg.Done()
	for msg := range c.weatherFeed {
		data := internal.CityWeatherData{}
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Fatal(err)
		}
		wfc.ProcessWeatherFeed(data)
	}
}
