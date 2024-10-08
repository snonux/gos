package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// The config file containing all the secrets and credentials.
type Secrets struct {
	MastodonURL         string
	MastodonAccessToken string
	LinkedInClientID    string
	LinkedInSecret      string
	LinkedInRedirectURL string
	// Will be updated by gos automatically, after successful oauth2
	LinkedInAccessToken string `json:"LinedInAccessToken,omitempty"`
	// Will be updated by gos automatically, after successful oauth2
	LinkedInPersonID string `json:"LinedInPersonID,omitempty"`
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

func (s Secrets) WriteToDisk(configPath string) error {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	tmpConfigPath := fmt.Sprintf("%s.tmp", configPath)
	file, err := os.Create(tmpConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return os.Rename(tmpConfigPath, configPath)
}
