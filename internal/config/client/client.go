package client

import (
	"strings"

	"codeberg.org/snonux/gos/internal/config"
)

type ClientConfig struct {
	Server string `json:"Partner,omitempty"`
	ApiKey string `json:"ApiKey,omitempty"`
}

func New(configFile string) (ClientConfig, error) {
	conf, _ := config.FromFile[ClientConfig](configFile)
	conf.Server = config.FromENV("Servers", conf.Server)
	conf.ApiKey = config.FromENV("ApiKey", conf.ApiKey)

	return conf, nil
}

func (conf ClientConfig) Servers() []string {
	return strings.Split(conf.Server, ",")
}
