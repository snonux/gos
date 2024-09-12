package server

import (
	"fmt"
	"log"
	"os"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/types"
)

type ServerConfig struct {
	ListenAddr             string        `json:"ListenAddr,omitempty"`
	Partners               []string      `json:"Partners,omitempty"`
	APIKey                 string        `json:"APIKey,omitempty"`
	DataDir                string        `json:"StateDir,omitempty"`
	EmailTo                string        `json:"EmailTo,omitempty"`
	EmailFrom              string        `json:"EmailFrom,omitempty"`
	SMTPServer             string        `json:"SMTPServer,omitempty"`
	MergeIntervalS         int           `json:"MergeInterval,omitempty"`
	ScheduleIntervalS      int           `json:"ScheduleInterval,omitempty"`
	SocialPlatformsEnabled []string      `json:"SocialPlatformsEnabled,omitempty"`
	Secrets                SecretsConfig `json:"Secrets,omitempty"`
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

	conf.ListenAddr = config.Str("GOS_LISTEN_ADDR", conf.ListenAddr, "localhost:8080")
	conf.Partners = config.StrSlice("GOS_PARTNERS", conf.Partners)
	conf.APIKey = config.Str("GOS_API_KEY", conf.APIKey)
	conf.DataDir = config.Str("GOS_DATA_DIR", conf.DataDir, "data")
	conf.EmailTo = config.Str("GOS_EMAIL_TO", conf.EmailTo)
	conf.EmailFrom = config.Str("GOS_EMAIL_FROM", conf.EmailFrom)
	conf.SocialPlatformsEnabled = config.StrSlice("GOS_SOCIAL_PLATFORMS_ENABLED",
		[]string{types.Mastodon, types.LinkedIn})

	conf.SMTPServer = config.Str("GOS_SMTP_SERVER", conf.SMTPServer, func() string {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
		return fmt.Sprintf("%s:25", hostname)
	})

	const oneHour = 3600
	conf.MergeIntervalS = config.Int("GOS_MERGE_INTERVAL", oneHour)
	conf.ScheduleIntervalS = config.Int("GOS_SCHEDULER_INTERVAL", oneHour*6)

	return conf, nil
}
