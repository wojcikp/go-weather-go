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
	CityData  chan internal.CityWeatherDataSingle
}

func NewApiDataProducer(
	apiClient apiclient.WeatherApiClient,
	CityData chan internal.CityWeatherDataSingle,
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
			w.CityData <- internal.CityWeatherDataSingle{
				Name:        city.Name,
				Time:        internal.CustomTime{},
				Temperature: 0.0,
				WindSpeed:   0.0,
				WeatherCode: 0,
				ErrorMsg:    err.Error(),
			}
		}
	}
	if err == nil {
		for i := 0; i < len(weatherData.Hourly.Time); i++ {
			w.CityData <- internal.CityWeatherDataSingle{
				Name:        city.Name,
				Time:        weatherData.Hourly.Time[i],
				Temperature: weatherData.Hourly.Temperature2m[i],
				WindSpeed:   weatherData.Hourly.WindSpeed[i],
				WeatherCode: weatherData.Hourly.WeatherCode[i],
				ErrorMsg:    "",
			}

		}
	}
}
