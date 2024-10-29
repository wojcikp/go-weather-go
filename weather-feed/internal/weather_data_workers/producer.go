package weatherdataworkers

import (
	"context"
	"log"
	"sync"
	"time"

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
	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		if err == nil {
			break
		}
		log.Printf("ERROR: Data for city: %s not fetched, err: %v\nRetrying...", apiJob.Name, err)
		time.Sleep(time.Second * time.Duration(i+1))
		weatherData, err = w.apiClient.FetchData(ctx, apiJob)
		if i == maxRetries-1 {
			log.Printf("ERROR: Last attempt to fetch data for city: %s failed. Putting on queue empty data for this city.", apiJob.Name)
			w.CityData <- internal.CityWeatherData{
				Name:         apiJob.Name,
				Time:         []internal.CustomTime{},
				Temperatures: []float64{},
				WindSpeed:    []float64{},
				WeatherCodes: []int{},
				Error:        err,
			}
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
