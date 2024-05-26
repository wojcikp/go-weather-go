package config

import "github.com/tkanos/gonfig"

type Configuration struct {
	BaseUrl              string
	LookBackwardInMonths int
	ConsumerCount        int
	MockCityInput        bool
}

func GetConfig() Configuration {
	configuration := Configuration{}
	if err := gonfig.GetConf("../../config/config.json", &configuration); err != nil {
		panic(err)
	}
	return configuration
}
