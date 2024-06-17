package server

import (
	"fmt"
	"log"
	"os"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
)

type ServerConfig struct {
	ListenAddr string `json:"ListenAddr,omitempty"`
	Partner    string `json:"Partner,omitempty"`
	APIKey     string `json:"APIKey,omitempty"`
	DataDir    string `json:"StateDir,omitempty"`
	EmailTo    string `json:"EmailTo,omitempty"`
	EmailFrom  string `json:"EmailFrom,omitempty"`
	SMTPServer string `json:"SMTPServer,omitempty"`
}

func New(configFile string) (ServerConfig, error) {
	conf, _ := config.FromFile[ServerConfig](configFile)
	conf.ListenAddr = config.FromENV("GOS_LISTEN_ADDR", conf.ListenAddr, "localhost:8080")
	conf.Partner = config.FromENV("GOS_PARTNER", conf.Partner)
	conf.APIKey = config.FromENV("GOS_API_KEY", conf.APIKey)
	conf.DataDir = config.FromENV("GOS_DATA_DIR", conf.DataDir, "data")
	conf.EmailTo = config.FromENV("GOS_EMAIL_TO", conf.EmailTo)
	conf.EmailFrom = config.FromENV("GOS_EMAIL_FROM", conf.EmailFrom)
	conf.SMTPServer = config.FromENV("GOS_SMTP_SERVER", conf.SMTPServer)

	if conf.SMTPServer == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
		conf.SMTPServer = fmt.Sprintf("%s:25", hostname)
		log.Println("Set SMTPServer to " + conf.SMTPServer)
	}

	return conf, nil
}

func (conf ServerConfig) Partners() []string {
	return strings.Split(conf.Partner, ",")
}
