package apiworker

import (
	"context"
	"log"
	"sync"

	"github.com/wojcikp/go-weather-go/weather-feed/internal"
	apiclient "github.com/wojcikp/go-weather-go/weather-feed/internal/api_client"
	"golang.org/x/sync/semaphore"
)

type ApiDataWorker struct {
	apiClient apiclient.WeatherApiClient
	CityData  chan internal.CityWeatherData
}

func NewApiDataWorker(
	apiClient apiclient.WeatherApiClient,
	CityData chan internal.CityWeatherData,
) *ApiDataWorker {
	return &ApiDataWorker{apiClient, CityData}
}

func (w ApiDataWorker) Work(ctx context.Context, apiJob internal.BaseCityInfo, wg *sync.WaitGroup, sem *semaphore.Weighted) {
	defer wg.Done()
	defer sem.Release(1)
	weatherData, err := w.apiClient.FetchData(ctx, apiJob)
	if err != nil {
		log.Fatal(err)
		log.Fatalf("Data for city: %v not fetched", apiJob.Name)
		w.CityData <- internal.CityWeatherData{
			Name:         apiJob.Name,
			Temperatures: []float64{},
			WeatherCodes: []int{},
			Error:        err,
		}
	}
	w.CityData <- internal.CityWeatherData{
		Name:         apiJob.Name,
		Temperatures: weatherData.Hourly.Temperature2m,
		WeatherCodes: weatherData.Hourly.WeatherCode,
		Error:        nil,
	}
}
