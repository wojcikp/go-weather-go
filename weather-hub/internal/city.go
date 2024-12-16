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

type ScoreInfo struct {
	Id    int
	Name  string
	Value string
}

type FeedStream struct {
	Data []byte
	Err  error
}

type IFeedReceiver interface {
	ReceiveMessages()
}
