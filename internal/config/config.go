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

type Config struct {
	ListenAddr string `json:"ListenAddr,omitempty"`
	Partner    string `json:"Partner,omitempty"`
	ApiKey     string `json:"ApiKey,omitempty"`
	DataDir    string `json:"StateDir,omitempty"`
	EmailTo    string `json:"EmailTo,omitempty"`
	EmailFrom  string `json:"EmailFrom,omitempty"`
	SMTPServer string `json:"SMTPServer,omitempty"`
}

func New(configFile string) (Config, error) {
	conf, _ := newFromConfigFile(configFile)
	conf.ListenAddr = fromEnv("ListenAddr", conf.ListenAddr, "localhost:8080")
	conf.Partner = fromEnv("Partner", conf.Partner)
	conf.ApiKey = fromEnv("ApiKey", conf.ApiKey)
	conf.DataDir = fromEnv("DataDir", conf.DataDir, "data")
	conf.EmailTo = fromEnv("EmailTo", conf.EmailTo)
	conf.EmailFrom = fromEnv("EmailFrom", conf.EmailFrom)
	conf.SMTPServer = fromEnv("SMTPServer", conf.SMTPServer)

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

func newFromConfigFile(configFile string) (Config, error) {
	var conf Config

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
	return conf, err
}

// Set config from envoronment variable if present, e.g. hansWurst from GOS_HANS_WURST
func fromEnv(configKey string, defaultValue ...string) string {
	envKey := camelToSnakeWithPrefix("GOS", configKey)
	if value := os.Getenv(envKey); value != "" {
		return value
	}

	// Use first non-empty default value.
	for _, value := range defaultValue {
		if value != "" {
			return value
		}
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
