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

func GetScoresInfo[T ScoreValue](scores []IWeatherScore[T], scoresInfo *bytes.Buffer, dbClient IDbClient) []error {
	var errors []error
	for _, score := range scores {
		id := score.GetId()
		name := score.GetName()
		scoresInfo.WriteString(fmt.Sprintf("Id: %d\n", id))
		scoresInfo.WriteString(fmt.Sprintf("Name: %s\n", name))
		scoresInfo.WriteString("Value: ")
		value, err := score.GetScore(dbClient)
		if err != nil {
			errors = append(errors, fmt.Errorf("ERROR: Score ID: %d, name: %s\nError: %v", id, name, err))
			scoresInfo.WriteString("Error occured")
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
	return errors
}
