package client

import (
	"strings"

	"codeberg.org/snonux/gos/internal/config"
)

type ClientConfig struct {
	Server string `json:"Partner,omitempty"`
	APIKey string `json:"APIKey,omitempty"`
	Editor string `json:"Editor,omitempty"`
}

func New(configFile string) (ClientConfig, error) {
	conf, _ := config.FromFile[ClientConfig](configFile)
	// TODO: Refactor
	conf.Server = config.FromENV("GOS_SERVERS", conf.Server)
	conf.APIKey = config.FromENV("GOS_API_KEY", conf.APIKey)
	conf.Editor = config.FromENV("GOS_EDITOR", "EDITOR", conf.Editor, "vi")

	return conf, nil
}

func (conf ClientConfig) Servers() []string {
	return strings.Split(conf.Server, ",")
}
