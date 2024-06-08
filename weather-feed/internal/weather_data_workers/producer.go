package weatherdataworkers

import (
	"context"
	"log"
	"sync"

	"github.com/wojcikp/go-weather-go/weather-feed/internal"
	apiclient "github.com/wojcikp/go-weather-go/weather-feed/internal/api_client"
	"golang.org/x/sync/semaphore"
)

type ApiDataProducer struct {
	apiClient apiclient.WeatherApiClient
	CityData  chan internal.CityWeatherData
}

func NewApiDataProducer(
	apiClient apiclient.WeatherApiClient,
	CityData chan internal.CityWeatherData,
) *ApiDataProducer {
	return &ApiDataProducer{apiClient, CityData}
}

func (w ApiDataProducer) Work(ctx context.Context, apiJob internal.BaseCityInfo, wg *sync.WaitGroup, sem *semaphore.Weighted) {
	defer wg.Done()
	defer sem.Release(1)
	weatherData, err := w.apiClient.FetchData(ctx, apiJob)
	if err != nil {
		log.Fatalf("Data for city: %s not fetched, err: %v", apiJob.Name, err)
		w.CityData <- internal.CityWeatherData{
			Name:         apiJob.Name,
			Time:         []internal.CustomTime{},
			Temperatures: []float64{},
			WindSpeed:    []float64{},
			WeatherCodes: []int{},
			Error:        err,
		}
	}
	w.CityData <- internal.CityWeatherData{
		Name:         apiJob.Name,
		Time:         weatherData.Hourly.Time,
		Temperatures: weatherData.Hourly.Temperature2m,
		WindSpeed:    weatherData.Hourly.WindSpeed,
		WeatherCodes: weatherData.Hourly.WeatherCode,
		Error:        nil,
	}
}
