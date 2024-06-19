package weatherdataworkers

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/wojcikp/go-weather-go/weather-feed/internal"
)

type IWeatherDataConsumer interface {
	ProcessWeatherData(data []byte)
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
			log.Fatal(err)
		}
		wdc.ProcessWeatherData(data)
	}
}
