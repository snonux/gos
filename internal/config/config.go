package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

type config struct {
	EmailTo    string `json:"EmailTo,omitempty"`
	EmailFrom  string `json:"EmailFrom,omitempty"`
	SMTPServer string `json:"SMTPServer,omitempty"`
	DataDir    string `json:"StateDir,omitempty"`
	Partner    string `json:"Partner,omitempty"`
}

func newConfig(configFile string) (config, error) {
	conf := config{
		EmailTo:    fromEnv("EmailTo"),
		EmailFrom:  fromEnv("EmailFrom"),
		SMTPServer: fromEnv("SMTPServer"),
		DataDir:    fromEnv("DataDir", "data"),
		Partner:    fromEnv("Partner"),
	}

	file, err := os.Open(configFile)
	if err != nil {
		return conf, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return conf, err
	}

	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		return conf, err
	}

	if conf.SMTPServer == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
		conf.SMTPServer = fmt.Sprintf("%s:25", hostname)
		log.Println("Set SMTPServer to " + conf.SMTPServer)
	}

	if conf.DataDir == "" {
		conf.DataDir = "data"
		log.Println("Set data dir to " + conf.DataDir)
	}

	return conf, nil
}

// Set config from envoronment variable if present, e.g. hansWurst from GOS_HANS_WURST
func fromEnv(configKey string, defaultValue ...string) string {
	envKey := camelToSnakeWithPrefix("GOS", configKey)
	if value := os.Getenv(envKey); value != "" {
		return value
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// camelToSnaeWithPrefix converts camelCase strings to UPPER_SNAKE_CASE with a prefix.
func camelToSnakeWithPrefix(prefix, s string) string {
	var builder strings.Builder
	builder.WriteString(strings.ToUpper(prefix))
	builder.WriteRune('_')

	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			builder.WriteRune('_')
		}
		builder.WriteRune(unicode.ToUpper(r))
	}

	return builder.String()
}
