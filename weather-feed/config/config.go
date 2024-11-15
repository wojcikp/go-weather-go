package config

import (
	"fmt"
	"os"
	"strconv"

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
	prod, err := strconv.ParseBool(os.Getenv("PRODUCTION"))
	if err != nil {
		prod = false
	}
	if prod {
		return "/app/config/config.json"
	} else {
		return "../../config/config.json"
	}
}
