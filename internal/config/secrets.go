package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Secrets struct {
	MastodonURL         string
	MastodonAccessToken string
	LinkedInClientID    string
	LinkedInSecret      string
	LinkedInRedirectURL string
}

func NewSecrets(configPath string) (Secrets, error) {
	var sec Secrets
	file, err := os.Open(configPath)
	if err != nil {
		return sec, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return sec, fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(bytes, &sec); err != nil {
		return sec, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return sec, nil
}
