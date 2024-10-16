package weatherscores

import (
	"bytes"
	"fmt"
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

func GetScoresInfo[T ScoreValue](scores []IWeatherScore[T], scoresInfo *bytes.Buffer, dbClient IDbClient) {
	for _, score := range scores {
		scoresInfo.WriteString(fmt.Sprintf("Id: %d\n", score.GetId()))
		scoresInfo.WriteString(fmt.Sprintf("Name: %s\n", score.GetName()))
		scoresInfo.WriteString("Value: ")
		value, err := score.GetScore(dbClient)
		if err != nil {
			scoresInfo.WriteString(err.Error())
		}
		switch v := any(value).(type) {
		case string:
			scoresInfo.WriteString(fmt.Sprintf("%s\n", v))
		case float64:
			scoresInfo.WriteString(fmt.Sprintf("%f\n", v))
		default:
			scoresInfo.WriteString("Unsupported score type\n")
		}
		scoresInfo.WriteString("-----------------------------\n")
	}
}
