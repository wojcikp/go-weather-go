package weatherfeedconsumer

import (
	"encoding/json"
	"log"

	"github.com/wojcikp/go-weather-go/weather-hub/internal"
)

type IWeatherFeedConsumer interface {
	ProcessWeatherFeed(data internal.CityWeatherDataSingle)
}

type Consumer struct {
	weatherFeed chan []byte
}

func NewWeatherFeedConsumer(weatherFeed chan []byte) *Consumer {
	return &Consumer{weatherFeed}
}

func (c Consumer) Work(done chan struct{}, wfc IWeatherFeedConsumer) {
	counter := 0
	for msg := range c.weatherFeed {
		data := internal.CityWeatherDataSingle{}
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Printf("ERROR: Unmarshalling data failed. Feed has not beed saved to db. error: %v", err)
		}
		// log.Printf("Processing data feed for city: %s", data.Name)
		wfc.ProcessWeatherFeed(data)
		counter++
		if counter%100 == 0 {
			log.Print("Process weather feed finished rows: ", counter)
		}
		done <- struct{}{}
	}
}
