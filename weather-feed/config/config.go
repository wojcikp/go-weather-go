package config

import (
	"fmt"
	"os"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	BaseUrl              string
	LookBackwardInMonths int
	ConsumerCount        int
	MockCityInput        bool
}

func GetConfig() (Configuration, error) {
	configuration := Configuration{}
	if err := gonfig.GetConf(getConfigPath(), &configuration); err != nil {
		return Configuration{}, fmt.Errorf("could not read config file, err: %w", err)
	}
	return configuration, nil
}

func getConfigPath() string {
	if os.Getenv("PRODUCTION") == "1" {
		return "/app/config/config.json"
	} else {
		return "../../config/config.json"
	}
}
