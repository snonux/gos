package server

import (
	"fmt"
	"log"
	"os"

	"codeberg.org/snonux/gos/internal/config"
)

type ServerConfig struct {
	ListenAddr        string   `json:"ListenAddr,omitempty"`
	Partners          []string `json:"Partners,omitempty"`
	APIKey            string   `json:"APIKey,omitempty"`
	DataDir           string   `json:"StateDir,omitempty"`
	EmailTo           string   `json:"EmailTo,omitempty"`
	EmailFrom         string   `json:"EmailFrom,omitempty"`
	SMTPServer        string   `json:"SMTPServer,omitempty"`
	MergeIntervalS    int      `json:"MergeInterval,omitempty"`
	ScheduleIntervalS int      `json:"ScheduleInterval,omitempty"`
	// SocialPlatformsEnable []string      `json:"SocialPlatformsEnable,omitempty"`
	Secrets SecretsConfig `json:"Secrets,omitempty"`
}

func New(configFile, secretsFile string) (ServerConfig, error) {
	conf, err := config.FromFile[ServerConfig](configFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return conf, err
		}
		log.Println("Skipping config file:", err)
	}

	if conf.Secrets, err = newSecretsConfig(secretsFile); err != nil {
		return conf, err
	}

	conf.ListenAddr = config.Env[config.Str]("GOS_LISTEN_ADDR", conf.ListenAddr, "localhost:8080")
	conf.Partners = config.Env[config.StrSlice]("GOS_PARTNERS", conf.Partners)
	conf.APIKey = config.Env[config.Str]("GOS_API_KEY", conf.APIKey)
	conf.DataDir = config.Env[config.Str]("GOS_DATA_DIR", conf.DataDir, "data")
	conf.EmailTo = config.Env[config.Str]("GOS_EMAIL_TO", conf.EmailTo)
	conf.EmailFrom = config.Env[config.Str]("GOS_EMAIL_FROM", conf.EmailFrom)

	conf.SMTPServer = config.Env[config.Str]("GOS_SMTP_SERVER", conf.SMTPServer, func() string {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
		return fmt.Sprintf("%s:25", hostname)
	})

	const oneHour = 3600
	conf.MergeIntervalS = config.Env[config.Int]("GOS_MERGE_INTERVAL", oneHour)
	conf.ScheduleIntervalS = config.Env[config.Int]("GOS_SCHEDULER_INTERVAL", oneHour*6)

	return conf, nil
}
