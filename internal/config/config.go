package config

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
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
func EnvToStr(keys ...any) string {
	for _, key := range keys {
		switch key := key.(type) {
		case string:
			if key == "" {
				continue
			}
			if !isAllUpperCase(key) {
				return key
			}
			if value := os.Getenv(key); value != "" {
				return value
			}
		case func() string:
			return key()
		}
	}

	return ""
}

func EnvToInt(keys ...any) int {
	for _, key := range keys {
		switch key := key.(type) {
		case string:
			if key == "" || !isAllUpperCase(key) {
				continue
			}
			strValue := os.Getenv(key)
			if strValue == "" {
				continue
			}
			if intValue, err := strconv.Atoi(strValue); err == nil {
				return intValue
			}
		case int:
			return key
		case func() int:
			return key()
		}
	}

	return 0
}

func isAllUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
		if unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
