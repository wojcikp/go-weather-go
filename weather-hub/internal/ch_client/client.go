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

type ClickhouseClient struct{}

func NewClickhouseClient() *ClickhouseClient {
	return &ClickhouseClient{}
}

func (c ClickhouseClient) CreateWeatherTable() {
	q := `CREATE TABLE IF NOT EXISTS weather_database.testowa_dwa
	(
		city String NOT NULL,
		time DateTime NOT NULL,
		temperature Float64,
		wind_speed Float64,
		weather_code Int64
	)
	ENGINE = ReplacingMergeTree
	PRIMARY KEY (city, time)`
	err := c.ExecQueryDb(q)
	if err != nil {
		log.Fatal(err)
	}
}

func (c ClickhouseClient) ProcessWeatherFeed(data internal.CityWeatherData) {
	if len(data.ErrorMsg) > 0 {
		log.Printf(
			"Weather data feed for following city: %s contains error: %s\nSkipping clickhouse insert based on this feed",
			data.Name,
			data.ErrorMsg,
		)
	} else {
		for i := 0; i < len(data.Time); i++ {
			q := fmt.Sprintf(
				"INSERT INTO weather_database.testowa_dwa (city, time, temperature, wind_speed, weather_code) VALUES ('%s', '%s', %f, %f, %d)",
				data.Name,
				data.Time[i].Format(time.DateTime),
				data.Temperatures[i],
				data.WindSpeed[i],
				data.WeatherCodes[i],
			)
			c.ExecQueryDb(q)
		}
		log.Printf("Processed data feed for city: %s", data.Name)
	}
}

func (c ClickhouseClient) ExecQueryDb(query string) error {
	conn, err := connect()
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
	conn, err := connect()
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

func connect() (driver.Conn, error) {
	host := os.Getenv("CLICKHOUSE_HOST")
	port := os.Getenv("CLICKHOUSE_PORT")
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", host, port)},
		Auth: clickhouse.Auth{
			Database: "weather_database",
			Username: os.Getenv("CLICKHOUSE_USER"),
			Password: os.Getenv("CLICKHOUSE_PASS"),
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
		return nil, err
	}
	return conn, err
}
