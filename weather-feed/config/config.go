package config

import (
	"fmt"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	BaseUrl              string
	LookBackwardInMonths int
	ConsumerCount        int
	MockCityInput        bool
}

const configPath = "/app/config/config.json"

func GetConfig() (Configuration, error) {
	configuration := Configuration{}
	if err := gonfig.GetConf(configPath, &configuration); err != nil {
		return Configuration{}, fmt.Errorf("could not read config file, err: %w", err)
	}
	return configuration, nil
}
