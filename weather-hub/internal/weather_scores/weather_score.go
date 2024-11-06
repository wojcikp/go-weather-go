package weatherscores

import (
	"fmt"

	"github.com/wojcikp/go-weather-go/weather-hub/internal"
)

type IDbClient interface {
	QueryDb(query string) (any, error)
}

type ScoreValue interface {
	string | float64
}

type IWeatherScore[T ScoreValue] interface {
	GetId() int
	GetName() string
	GetQuery() string
	GetScore(dbClient IDbClient) (T, error)
}

type BaseWeatherScore struct {
	Id int
}

func (ws *BaseWeatherScore) GetId() int {
	return ws.Id
}

func GetScoresInfo[T ScoreValue](scores []IWeatherScore[T], dbClient IDbClient) ([]internal.ScoreInfo, []error) {
	errors := []error{}
	infos := []internal.ScoreInfo{}
	for _, score := range scores {
		id := score.GetId()
		name := score.GetName()
		scoreInfo := internal.ScoreInfo{
			Id:   id,
			Name: name,
		}
		value, err := score.GetScore(dbClient)
		if err != nil {
			errors = append(errors, fmt.Errorf("ERROR: Score ID: %d, name: %s\nError: %v", id, name, err))
		}
		switch v := any(value).(type) {
		case string:
			scoreInfo.Value = v
		case float64:
			scoreInfo.Value = fmt.Sprintf("%f", v)
		default:
			scoreInfo.Value = "Unsupported score type\n"
		}
		infos = append(infos, scoreInfo)
	}
	return infos, errors
}
