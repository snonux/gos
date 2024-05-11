package server

import (
	"fmt"
	"log"
	"os"

	"codeberg.org/snonux/gos/internal/config"
)

type ServerConfig struct {
	ListenAddr string `json:"ListenAddr,omitempty"`
	Partner    string `json:"Partner,omitempty"`
	ApiKey     string `json:"ApiKey,omitempty"`
	DataDir    string `json:"StateDir,omitempty"`
	EmailTo    string `json:"EmailTo,omitempty"`
	EmailFrom  string `json:"EmailFrom,omitempty"`
	SMTPServer string `json:"SMTPServer,omitempty"`
}

func New(configFile string) (ServerConfig, error) {
	conf, _ := config.FromFile[ServerConfig](configFile)
	conf.ListenAddr = config.FromENV("ListenAddr", conf.ListenAddr, "localhost:8080")
	conf.Partner = config.FromENV("Partner", conf.Partner)
	conf.ApiKey = config.FromENV("ApiKey", conf.ApiKey)
	conf.DataDir = config.FromENV("DataDir", conf.DataDir, "data")
	conf.EmailTo = config.FromENV("EmailTo", conf.EmailTo)
	conf.EmailFrom = config.FromENV("EmailFrom", conf.EmailFrom)
	conf.SMTPServer = config.FromENV("SMTPServer", conf.SMTPServer)

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
