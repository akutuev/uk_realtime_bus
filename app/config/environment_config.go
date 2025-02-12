package config

import (
	"github.com/caarlos0/env/v6"
)

type EnvSettings struct {
	BusDataApiKey    string   `env:"BUS_DATA_API_KEY,required"`
	BusDataHost      string   `env:"BUS_DATA_HOST,required"`
	BusOperatorRef   string   `env:"BUS_OPERATOR_REF,required"`
	BusesToTrackList []string `env:"BUSES_TO_TRACK_LIST,required"`
}

func NewEnvSettings() *EnvSettings {
	cfg := &EnvSettings{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	return cfg
}
