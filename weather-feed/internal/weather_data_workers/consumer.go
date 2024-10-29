package weatherdataworkers

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/wojcikp/go-weather-go/weather-feed/internal"
)

type IWeatherDataConsumer interface {
	ProcessWeatherData(data []byte) error
}

type Consumer struct {
	cityData chan internal.CityWeatherData
}

func NewWeatherDataConsumer(cityData chan internal.CityWeatherData) *Consumer {
	return &Consumer{cityData}
}

func (c Consumer) Work(wg *sync.WaitGroup, wdc IWeatherDataConsumer) {
	defer wg.Done()
	for msg := range c.cityData {
		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("ERROR: Marshalling data to json format for city: %s failed due to following error: %v", msg.Name, err)
		}
		if err = wdc.ProcessWeatherData(data); err != nil {
			log.Printf("ERROR: Data for city: %s was not processed due to following error: %v", msg.Name, err)
		}
	}
}
