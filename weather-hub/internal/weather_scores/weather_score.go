package weatherscores

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
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
	GetQuery() (string, error)
	GetScore(dbClient IDbClient) (T, error)
}

type BaseWeatherScore struct {
	Id int
}

func (ws *BaseWeatherScore) GetId() int {
	return ws.Id
}

func (ws *BaseWeatherScore) GetQueryResults(dbClient IDbClient, query string) ([][]interface{}, error) {
	var results [][]interface{}

	data, err := dbClient.QueryDb(query)
	if err != nil {
		return nil, err
	}

	rows, ok := data.(driver.Rows)
	if !ok {
		return nil, fmt.Errorf("return data is not clickhouse rows type, err: %w", err)
	}
	defer rows.Close()
	var (
		columnTypes = rows.ColumnTypes()
		vars        = make([]interface{}, len(columnTypes))
	)
	for i := range columnTypes {
		vars[i] = reflect.New(columnTypes[i].ScanType()).Interface()
	}
	for rows.Next() {
		if err := rows.Scan(vars...); err != nil {
			return nil, err
		}
		row := make([]interface{}, len(vars))
		for i, v := range vars {
			switch v := v.(type) {
			case *string:
				row[i] = *v
			case *uint64:
				row[i] = *v
			case *float64:
				row[i] = *v
			case *time.Time:
				row[i] = *v
			}
		}
		results = append(results, row)
	}
	return results, nil
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
