package client

import (
	"fmt"
	"os"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
)

type ClientConfig struct {
	Server  string `json:"Partner,omitempty"`
	APIKey  string `json:"APIKey,omitempty"`
	Editor  string `json:"Editor,omitempty"`
	DataDir string `json:"StateDir,omitempty"`
}

func New(configFile string) (ClientConfig, error) {
	conf, _ := config.FromFile[ClientConfig](configFile)

	conf.Server = config.FromENV("GOS_SERVERS", conf.Server)
	conf.APIKey = config.FromENV("GOS_API_KEY", conf.APIKey)
	conf.Editor = config.FromENV("GOS_EDITOR", "EDITOR", conf.Editor, "vi")

	defaultDataDir := fmt.Sprintf("%s/.gos/data", os.Getenv("HOME"))
	conf.DataDir = config.FromENV("GOS_DATA_DIR", conf.DataDir, defaultDataDir)

	return conf, nil
}

func (conf ClientConfig) Servers() []string {
	return strings.Split(conf.Server, ",")
}
