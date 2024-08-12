package server

import (
	"fmt"
	"log"
	"os"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
)

type ServerConfig struct {
	ListenAddr         string `json:"ListenAddr,omitempty"`
	Partner            string `json:"Partner,omitempty"`
	APIKey             string `json:"APIKey,omitempty"`
	DataDir            string `json:"StateDir,omitempty"`
	EmailTo            string `json:"EmailTo,omitempty"`
	EmailFrom          string `json:"EmailFrom,omitempty"`
	SMTPServer         string `json:"SMTPServer,omitempty"`
	CRONMergeIntervalS int    `json:"CRONMergeInterval,omitempty"`
}

func New(configFile string) (ServerConfig, error) {
	conf, _ := config.FromFile[ServerConfig](configFile)

	conf.ListenAddr = config.EnvToStr("GOS_LISTEN_ADDR", conf.ListenAddr, "localhost:8080")
	conf.Partner = config.EnvToStr("GOS_PARTNER", conf.Partner)
	conf.APIKey = config.EnvToStr("GOS_API_KEY", conf.APIKey)
	conf.DataDir = config.EnvToStr("GOS_DATA_DIR", conf.DataDir, "data")
	conf.EmailTo = config.EnvToStr("GOS_EMAIL_TO", conf.EmailTo)
	conf.EmailFrom = config.EnvToStr("GOS_EMAIL_FROM", conf.EmailFrom)

	conf.SMTPServer = config.EnvToStr("GOS_SMTP_SERVER", conf.SMTPServer, func() string {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
		return fmt.Sprintf("%s:25", hostname)
	})

	conf.CRONMergeIntervalS = config.EnvToInt("GOS_CRON_MERGE_INTERVAL", 3600)
	return conf, nil
}

func (conf ServerConfig) Partners() []string {
	if partners := strings.Split(conf.Partner, ","); partners[0] != "" {
		return partners
	}
	return []string{}
}
