package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/wojcikp/go-weather-go/weather-feed/internal"
)

type Hourly struct {
	Time          []internal.CustomTime `json:"time"`
	Temperature2m []float64             `json:"temperature_2m"`
	WindSpeed     []float64             `json:"wind_speed_10m"`
	WeatherCode   []int                 `json:"weather_code"`
}

type ApiResponse struct {
	Hourly Hourly
}

type WeatherApiClient struct {
	baseUrl              string
	lookBackwardInMonths int
}

func NewApiClient(baseUrl string, lookBackwardInMonths int) *WeatherApiClient {
	return &WeatherApiClient{baseUrl, lookBackwardInMonths}
}

func (c WeatherApiClient) FetchData(ctx context.Context, cityInfo internal.BaseCityInfo) (ApiResponse, error) {
	timeout := 3 * time.Second
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	const twoDaysAgo = -2
	endDate := time.Now().AddDate(0, 0, twoDaysAgo)
	startDate := endDate.AddDate(0, c.lookBackwardInMonths, 0).Format("2006-01-02")
	url := fmt.Sprintf(
		"%s?latitude=%s&longitude=%s&start_date=%s&end_date=%s&hourly=temperature_2m,weather_code,wind_speed_10m",
		c.baseUrl, cityInfo.Latitude, cityInfo.Longtitude, startDate, endDate.Format("2006-01-02"),
	)

	req, err := http.NewRequestWithContext(ctxTimeout, http.MethodGet, url, nil)
	if err != nil {
		return ApiResponse{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ApiResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		errorStr := fmt.Sprintf("response failed with status code: %d and\nbody: %s", res.StatusCode, res.Body)
		return ApiResponse{}, fmt.Errorf(errorStr)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return ApiResponse{}, err
	}

	var responseData ApiResponse
	if err := json.Unmarshal(data, &responseData); err != nil {
		return ApiResponse{}, err
	}

	return responseData, nil
}
