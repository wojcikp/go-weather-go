package internal

import (
	"time"
)

type CityWeatherData struct {
	Name         string
	Time         []time.Time
	Temperatures []float64
	WindSpeed    []float64
	WeatherCodes []int
	ErrorMsg     string
}

type CityWeatherDataSingle struct {
	Name        string
	Time        time.Time
	Temperature float64
	WindSpeed   float64
	WeatherCode int
	ErrorMsg    string
}

type ScoreInfo struct {
	Id    int
	Name  string
	Value string
}
