package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"codeberg.org/snonux/gos/internal/colour"
)

// The config file containing all the secrets and credentials.
type Secrets struct {
	MastodonURL         string
	MastodonAccessToken string
	LinkedInClientID    string
	LinkedInSecret      string
	LinkedInRedirectURL string
	// Will be updated by gos automatically, after successful oauth2
	LinkedInAccessToken string `json:"LinkedInAccessToken,omitempty"`
	// Will be updated by gos automatically, after successful oauth2
	LinkedInPersonID string `json:"LinkedInPersonID,omitempty"`
}

func NewSecrets(configPath string, composeEntry bool) (Secrets, error) {
	var sec Secrets
	if composeEntry {
		// In compose mode, no need to read the secrets.
		return sec, nil
	}

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
	colour.Infoln("Writing", configPath)

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
