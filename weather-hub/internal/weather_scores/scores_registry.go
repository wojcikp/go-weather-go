package weatherscores

import "log"

func GetScoresList[T ScoreValue]() []IWeatherScore[T] {
	var dummy T
	switch v := any(dummy).(type) {
	case string:
		return []IWeatherScore[T]{
			&HighestAvgTempCity7d[T]{},
			&MostRainyCity7d[T]{},
		}
	case float64:
		return []IWeatherScore[T]{}
	default:
		log.Fatalf("Unsupported feature type: %v", v)
		return []IWeatherScore[T]{}
	}
}
