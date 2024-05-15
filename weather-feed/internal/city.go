package internal

type CityWeatherData struct {
	Name         string
	Temperatures []float64
	WeatherCodes []int
}

type consumedCityData struct {
	avgTemperature float64
	weatherCodes   []int
}

type BaseCityInfo struct {
	Name       string
	Latitude   string
	Longtitude string
}
