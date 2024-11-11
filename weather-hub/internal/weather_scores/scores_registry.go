package weatherscores

import "log"

var idCounter int
var nextId = func() int {
	idCounter++
	return idCounter
}

func GetScoresList[T ScoreValue]() []IWeatherScore[T] {
	var dummy T
	switch v := any(dummy).(type) {
	case string:
		return []IWeatherScore[T]{
			&HighestAvgTempCity7d[T]{BaseWeatherScore{Id: nextId()}},
			&MostRainyCity7d[T]{BaseWeatherScore{Id: nextId()}},
			&MostRainyCity31d[T]{BaseWeatherScore{Id: nextId()}},
		}
	case float64:
		return []IWeatherScore[T]{
			&AvgTemperatureWarsaw14d[T]{BaseWeatherScore{Id: nextId()}},
		}
	default:
		log.Printf("ERROR: Unsupported score type: %v", v)
		return []IWeatherScore[T]{}
	}
}
