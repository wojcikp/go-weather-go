package internal

import (
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

type CityWeatherDataSingle struct {
	Name        string
	Time        CustomTime
	Temperature float64
	WindSpeed   float64
	WeatherCode int
	ErrorMsg    string
}

type CityWeatherData struct {
	Name         string
	Time         []CustomTime
	Temperatures []float64
	WindSpeed    []float64
	WeatherCodes []int
	ErrorMsg     string
}

type BaseCityInfo struct {
	Name       string
	Latitude   string
	Longtitude string
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse("2006-01-02T15:04", s)
	return
}
