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

func (w ApiDataProducer) Work(ctx context.Context, city internal.BaseCityInfo, wg *sync.WaitGroup, sem *semaphore.Weighted) {
	defer wg.Done()
	defer sem.Release(1)
	weatherData, err := w.apiClient.FetchData(ctx, city)
	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		if err == nil {
			break
		}
		log.Printf("ERROR: Data for city: %s not fetched, err: %v\nRetrying...", city.Name, err)
		time.Sleep(time.Second * time.Duration(i+1))
		weatherData, err = w.apiClient.FetchData(ctx, city)
		if i == maxRetries-1 {
			log.Printf("ERROR: Last attempt to fetch data for city: %s failed. Putting on queue empty data for this city.", city.Name)
			w.CityData <- internal.CityWeatherData{
				Name:         city.Name,
				Time:         []internal.CustomTime{},
				Temperatures: []float64{},
				WindSpeed:    []float64{},
				WeatherCodes: []int{},
				ErrorMsg:     err.Error(),
			}
		}
	}
	if err == nil {
		w.CityData <- internal.CityWeatherData{
			Name:         city.Name,
			Time:         weatherData.Hourly.Time,
			Temperatures: weatherData.Hourly.Temperature2m,
			WindSpeed:    weatherData.Hourly.WindSpeed,
			WeatherCodes: weatherData.Hourly.WeatherCode,
			ErrorMsg:     "",
		}
	}
}
