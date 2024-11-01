package weatherscores

import (
	"fmt"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type HighestAvgTempCity7d[T ScoreValue] struct {
	BaseWeatherScore
}

func (wc *HighestAvgTempCity7d[T]) GetName() string {
	return "Highest_Avg_Temp_City_7d"
}

func (wc *HighestAvgTempCity7d[T]) GetQuery() string {
	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf(
		`SELECT
    		city,
    		avg(temperature) as avgTemp
		FROM %s.%s
		WHERE (time >= '%s 00:00:00') AND (time <= '%s 00:00:00')
		GROUP BY city
		ORDER BY avgTemp DESC`,
		os.Getenv("CLICKHOUSE_DB"),
		os.Getenv("CLICKHOUSE_TABLE"),
		startDate,
		today,
	)
}

func (wc *HighestAvgTempCity7d[T]) GetScore(dbClient IDbClient) (T, error) {
	var empty T
	data, err := dbClient.QueryDb(wc.GetQuery())
	if err != nil {
		return empty, err
	}

	rows, ok := data.(driver.Rows)
	if !ok {
		return empty, fmt.Errorf("return data is not clickhouse rows type, err: %w", err)
	}
	defer rows.Close()

	var cities []string
	for rows.Next() {
		var (
			city    string
			avgTemp float64
		)
		if err := rows.Scan(&city, &avgTemp); err != nil {
			return empty, err
		}
		cities = append(cities, city)
	}

	if len(cities) == 0 {
		return empty, err
	}

	score, ok := any(cities[0]).(T)
	if !ok {
		return empty, fmt.Errorf("wrong data type for score %s with id: %d", wc.GetName(), wc.GetId())
	}

	return score, err
}
