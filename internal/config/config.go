package config

import (
	"encoding/json"
	"io"
	"os"
)

func FromFile[T any](configFile string) (T, error) {
	var conf T

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
func FromENV(envKey string, defaultValue ...string) string {
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
