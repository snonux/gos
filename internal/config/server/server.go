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

	conf.ListenAddr = config.StrFromENV("GOS_LISTEN_ADDR", conf.ListenAddr, "localhost:8080")
	conf.Partner = config.StrFromENV("GOS_PARTNER", conf.Partner)
	conf.APIKey = config.StrFromENV("GOS_API_KEY", conf.APIKey)
	conf.DataDir = config.StrFromENV("GOS_DATA_DIR", conf.DataDir, "data")
	conf.EmailTo = config.StrFromENV("GOS_EMAIL_TO", conf.EmailTo)
	conf.EmailFrom = config.StrFromENV("GOS_EMAIL_FROM", conf.EmailFrom)

	conf.SMTPServer = config.StrFromENV("GOS_SMTP_SERVER", conf.SMTPServer)
	if conf.SMTPServer == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
		conf.SMTPServer = fmt.Sprintf("%s:25", hostname)
		log.Println("Set SMTPServer to " + conf.SMTPServer)
	}

	conf.CRONMergeIntervalS = config.IntFromENV("GOS_CRON_MERGE_INTERVAL", 3600)
	return conf, nil
}

func (conf ServerConfig) Partners() []string {
	return strings.Split(conf.Partner, ",")
}
