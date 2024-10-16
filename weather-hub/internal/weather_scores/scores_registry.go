package weatherscores

import "log"

func GetScoresList[T ScoreValue]() []IWeatherScore[T] {
	var dummy T
	var idCounter int
	nextId := func() int {
		idCounter++
		return idCounter
	}
	switch v := any(dummy).(type) {
	case string:
		return []IWeatherScore[T]{
			&HighestAvgTempCity7d[T]{BaseWeatherScore{Id: nextId()}},
			&MostRainyCity7d[T]{BaseWeatherScore{Id: nextId()}},
		}
	case float64:
		return []IWeatherScore[T]{}
	default:
		log.Fatalf("Unsupported score type: %v", v)
		return []IWeatherScore[T]{}
	}
}
