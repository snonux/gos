package client

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
)

type ClientConfig struct {
	Server      string `json:"Partner,omitempty"`
	APIKey      string `json:"APIKey,omitempty"`
	Editor      string `json:"Editor,omitempty"`
	DataDir     string `json:"StateDir,omitempty"`
	ComposeFile string `json:"ComposeFile,omitempty"`
	LogFile     string `json:"LogFile,omitempty"`
}

func New(configFile string) (ClientConfig, error) {
	conf, _ := config.FromFile[ClientConfig](configFile)

	conf.Server = config.EnvToStr("GOS_SERVERS", conf.Server)
	conf.APIKey = config.EnvToStr("GOS_API_KEY", conf.APIKey)
	conf.Editor = config.EnvToStr("GOS_EDITOR", "EDITOR", conf.Editor, "vi")

	defaultDataDir := fmt.Sprintf("%s/.gos/data", os.Getenv("HOME"))
	conf.DataDir = config.EnvToStr("GOS_DATA_DIR", conf.DataDir, defaultDataDir)
	conf.ComposeFile = config.EnvToStr("GOS_COMPOSE_FILE", conf.ComposeFile, "compose.txt")

	defaultLogFile := fmt.Sprintf("%s/.gos/gos.log", os.Getenv("HOME"))
	conf.LogFile = config.EnvToStr("GOS_LOG_FILE", conf.LogFile, defaultLogFile)

	return conf, nil
}

func (conf ClientConfig) Servers() ([]string, error) {
	var servers []string

	for _, server := range strings.Split(conf.Server, ",") {
		if server != "" {
			servers = append(servers, server)
		}
	}

	if len(servers) == 0 {
		return servers, errors.New("no server(s) configured")
	}

	return servers, nil
}
