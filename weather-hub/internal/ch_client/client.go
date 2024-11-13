package chclient

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/wojcikp/go-weather-go/weather-hub/internal"
)

type ClickhouseClient struct {
	db    string
	table string
}

func NewClickhouseClient(db, table string) *ClickhouseClient {
	return &ClickhouseClient{db, table}
}

func (c ClickhouseClient) CreateWeatherTable() {
	q := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s
	(
		city String NOT NULL,
		time DateTime NOT NULL,
		temperature Float64,
		wind_speed Float64,
		weather_code Int64
	)
	ENGINE = ReplacingMergeTree
	PRIMARY KEY (city, time)`, c.db, c.table)
	err := c.ExecQueryDb(q)
	if err != nil {
		log.Fatalf("Creating table: %s in clickhouse db: %s failed due to following error: %v.\nexecuted query: %s", c.table, c.db, err, q)
	}
}

func (c ClickhouseClient) ProcessWeatherFeed(code *int, data internal.CityWeatherData) {
	conn, err := connect(c.db)
	if err != nil {
		log.Print("ERROR: ", err)
	}
	defer conn.Close()
	if len(data.ErrorMsg) > 0 {
		log.Printf(
			"Weather data feed for following city: %s contains error: %s\nSkipping clickhouse insert based on this feed",
			data.Name,
			data.ErrorMsg,
		)
	} else {
		var executeQueryErrors []error

		for _, batch := range getBatchedQueries(data) {
			stmt, err := conn.PrepareBatch(
				context.Background(),
				fmt.Sprintf("INSERT INTO %s.%s (city, time, temperature, wind_speed, weather_code) VALUES", c.db, c.table),
			)
			if err != nil {
				executeQueryErrors = append(executeQueryErrors, err)
			}
			for _, query := range batch {
				query, ok := query.([]any)
				if !ok {
					executeQueryErrors = append(executeQueryErrors, err)
				}
				// if 2 == 2 {
				// 	log.Print()
				// }
				if err := stmt.Append(query...); err != nil {
					executeQueryErrors = append(executeQueryErrors, err)
				}
				name, ok := query[0].(string)
				if !ok {
					log.Print("EEEEEEEEEEEEEEEER: no name")
				}
				date, ok := query[1].(string)
				if !ok {
					log.Print("EEEEEEEEEEEEEEEER: no date")
				}
				parsedDate, err := time.Parse("2006-01-02 15:04:05", date)
				if err != nil {
					log.Print("parse date err: ", err)
				}

				// temp, ok := query[2].(float64)
				// if !ok {
				// 	log.Print("EEEEEEEEEEEEEEEER: no temp")
				// }
				// wind, ok := query[3].(float64)
				// if !ok {
				// 	log.Print("EEEEEEEEEEEEEEEER: no wind")
				// }
				c, ok := query[4].(int)
				if !ok {
					log.Print("EEEEEEEEEEEEEEEER: no weather codes")
				}
				if parsedDate.Before(time.Now()) && parsedDate.After(time.Now().AddDate(0, -1, 0)) && name == "Złotów" {
					// fmt.Print(name)
					fmt.Print(date)
					// fmt.Print(temp)
					// fmt.Print(wind)
					fmt.Print(c)
				}
				*code += c
			}
			if err := stmt.Send(); err != nil {
				executeQueryErrors = append(executeQueryErrors, err)
			}
		}
		if len(executeQueryErrors) > 0 {
			log.Printf("ERROR: Processing weather feed for city: %s failed. Errors: %v", data.Name, executeQueryErrors)
		}
		log.Printf("Processed data feed for city: %s", data.Name)
	}
}

func getBatchedQueries(data internal.CityWeatherData) [][]any {
	var batchedQueries [][]any
	batchSize := 100
	maxValue := len(data.Time) - 1

	for i := 0; i <= maxValue; i += batchSize {
		batch := make([]any, 0, batchSize)
		for j := i; j < i+batchSize && j <= maxValue; j++ {

			q := []any{
				data.Name,
				data.Time[j].Format(time.DateTime),
				data.Temperatures[j],
				data.WindSpeed[j],
				data.WeatherCodes[j],
			}
			batch = append(batch, q)
		}
		batchedQueries = append(batchedQueries, batch)
	}
	return batchedQueries
}

func (c ClickhouseClient) ExecQueryDb(query string) error {
	conn, err := connect(c.db)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx := context.Background()
	err = conn.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (c ClickhouseClient) QueryDb(query string) (any, error) {
	conn, err := connect(c.db)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx := context.Background()
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func connect(db string) (driver.Conn, error) {
	host := os.Getenv("CLICKHOUSE_HOST")
	port := os.Getenv("CLICKHOUSE_PORT")
	user := os.Getenv("CLICKHOUSE_USER")
	pass := os.Getenv("CLICKHOUSE_PASS")
	if host == "" || port == "" || user == "" || pass == "" {
		return nil, fmt.Errorf("missing db credentials - host: %s, port: %s, user: %s, pass: %s", host, port, user, pass)
	}
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", host, port)},
		Auth: clickhouse.Auth{
			Database: db,
			Username: user,
			Password: pass,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: false,
		Debugf: func(format string, v ...any) {
			log.Printf(format+"\n", v...)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
	})
	if err != nil {
		return nil, fmt.Errorf(
			"connecting to clickhouse db: %s, host: %s, port %s failed due to following error %w", db, host, port, err)
	}
	return conn, nil
}
