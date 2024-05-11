package config

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"unicode"
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
func FromENV(configKey string, defaultValue ...string) string {
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
