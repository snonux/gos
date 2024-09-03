package client

import (
	"fmt"
	"log"
	"os"

	"codeberg.org/snonux/gos/internal/config"
)

type ClientConfig struct {
	Servers     []string `json:"Servers,omitempty"`
	APIKey      string   `json:"APIKey,omitempty"`
	Editor      string   `json:"Editor,omitempty"`
	DataDir     string   `json:"StateDir,omitempty"`
	ComposeFile string   `json:"ComposeFile,omitempty"`
	LogFile     string   `json:"LogFile,omitempty"`
}

func New(configFile string) (ClientConfig, error) {
	conf, err := config.FromFile[ClientConfig](configFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return conf, err
		}
		log.Println("Skipping config file:", err)
	}

	conf.Servers = config.EnvToStrSlice("GOS_SERVERS", conf.Servers)
	conf.APIKey = config.EnvToStr("GOS_API_KEY", conf.APIKey)
	conf.Editor = config.EnvToStr("GOS_EDITOR", "EDITOR", conf.Editor, "vi")

	defaultDataDir := fmt.Sprintf("%s/.gos/data", os.Getenv("HOME"))
	conf.DataDir = config.EnvToStr("GOS_DATA_DIR", conf.DataDir, defaultDataDir)
	conf.ComposeFile = config.EnvToStr("GOS_COMPOSE_FILE", conf.ComposeFile, "compose.txt")

	defaultLogFile := fmt.Sprintf("%s/.gos/gos.log", os.Getenv("HOME"))
	conf.LogFile = config.EnvToStr("GOS_LOG_FILE", conf.LogFile, defaultLogFile)

	return conf, nil
}
