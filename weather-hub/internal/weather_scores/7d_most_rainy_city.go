package weatherscores

import (
	"fmt"
	"log"
	"os"
	"time"
)

type MostRainyCity7d[T ScoreValue] struct {
	BaseWeatherScore
}

func (wc *MostRainyCity7d[ScoreValue]) GetName() string {
	return "Most_Rainy_City_7d"
}

func (wc *MostRainyCity7d[ScoreValue]) GetQuery() (string, error) {
	db := os.Getenv("CLICKHOUSE_DB")
	table := os.Getenv("CLICKHOUSE_TABLE")
	if db == "" || table == "" {
		return "", fmt.Errorf("missing CLICKHOUSE_DB or CLICKHOUSE_TABLE os envs, db: %s, table: %s", db, table)
	}
	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf(
		`SELECT
			city,
			count(*) AS code_count
		FROM %s.%s FINAL
		WHERE ((weather_code = 61) OR (weather_code = 63) OR (weather_code = 65)
		OR (weather_code = 80) OR (weather_code = 81) OR (weather_code = 82))
		AND ((time >= '%s 00:00:00') AND (time <= '%s 00:00:00'))
		GROUP BY city
		ORDER BY code_count DESC`,
		db,
		table,
		startDate,
		today,
	), nil
}

func (wc *MostRainyCity7d[ScoreValue]) GetScore(dbClient IDbClient) (ScoreValue, error) {
	var empty ScoreValue
	query, err := wc.GetQuery()
	if err != nil {
		return empty, err
	}

	results, err := wc.GetQueryResults(dbClient, query)
	if err != nil {
		log.Print("ERROR GetQueryResults: ", err)
	}

	if len(results) == 0 {
		if noRainMessage, ok := any("No rain in the past 7 days.").(ScoreValue); ok {
			empty = noRainMessage
		}
		return empty, err
	}

	score, ok := any(results[0][0]).(ScoreValue)
	if !ok {
		return empty, fmt.Errorf("wrong data type for score %s with id: %d", wc.GetName(), wc.GetId())
	}

	return score, err
}
