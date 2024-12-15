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
	weatherFeed chan internal.FeedStream
}

func NewWeatherFeedConsumer(weatherFeed chan internal.FeedStream) *Consumer {
	return &Consumer{weatherFeed}
}

func (c Consumer) Work(done chan struct{}, wfc IWeatherFeedConsumer) {
	for msg := range c.weatherFeed {
		if msg.Err != nil {
			log.Printf("ERROR: Weather data feed consuming error: %v", msg.Err)
		}
		data := internal.CityWeatherData{}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("ERROR: Unmarshalling data failed. Feed has not beed saved to db. error: %v", err)
		}
		wfc.ProcessWeatherFeed(data)
		done <- struct{}{}
	}
}
