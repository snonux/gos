package server

import (
	"fmt"
	"log"
	"os"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
)

type ServerConfig struct {
	ListenAddr        string `json:"ListenAddr,omitempty"`
	Partner           string `json:"Partner,omitempty"`
	APIKey            string `json:"APIKey,omitempty"`
	DataDir           string `json:"StateDir,omitempty"`
	EmailTo           string `json:"EmailTo,omitempty"`
	EmailFrom         string `json:"EmailFrom,omitempty"`
	SMTPServer        string `json:"SMTPServer,omitempty"`
	MergeIntervalS    int    `json:"MergeInterval,omitempty"`
	ScheduleIntervalS int    `json:"ScheduleInterval,omitempty"`
}

func New(configFile string) (ServerConfig, error) {
	conf, _ := config.FromFile[ServerConfig](configFile)

	conf.ListenAddr = config.EnvToStr("GOS_LISTEN_ADDR", conf.ListenAddr, "localhost:8080")
	conf.Partner = config.EnvToStr("GOS_PARTNER", "GOS_PARTNERS", conf.Partner)
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

	const oneHour = 3600
	conf.MergeIntervalS = config.EnvToInt("GOS_MERGE_INTERVAL", oneHour)
	conf.ScheduleIntervalS = config.EnvToInt("GOS_SCHEDULER_INTERVAL", oneHour*6)

	return conf, nil
}

func (conf ServerConfig) Partners() ([]string, error) {
	if partners := strings.Split(conf.Partner, ","); partners[0] != "" {
		return partners, nil
	}

	return []string{}, fmt.Errorf("no partners configured")
}
