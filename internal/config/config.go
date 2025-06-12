package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"codeberg.org/snonux/gos/internal/colour"
)

// The config file containing all the secrets and credentials plus maybe more.
type Config struct {
	LastRunEpoch        int64 `json:"LastRunEpoch,omitempty"`
	MastodonURL         string
	MastodonAccessToken string
	LinkedInClientID    string
	LinkedInSecret      string
	LinkedInRedirectURL string
	// Will be updated by gos automatically, after successful oauth2
	LinkedInAccessToken string `json:"LinkedInAccessToken,omitempty"`
	// Will be updated by gos automatically, after successful oauth2
	LinkedInPersonID string `json:"LinkedInPersonID,omitempty"`
	// Pause posting between these dates (format: "2006-01-02")
	PauseStart string `json:"PauseStart,omitempty"`
	PauseEnd   string `json:"PauseEnd,omitempty"`
}

func New(configPath string, composeEntry bool) (Config, error) {
	var conf Config

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		if !composeEntry {
			return conf, fmt.Errorf("No config file %s", configPath)
		}
		// Create empty new config for compose mode.
		return conf, conf.WriteToDisk(configPath)
	}

	file, err := os.Open(configPath)
	if err != nil {
		return conf, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return conf, fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(bytes, &conf); err != nil {
		return conf, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return conf, nil
}

func (s Config) WriteToDisk(configPath string) error {
	colour.Infoln("Writing", configPath)
	if err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm); err != nil {
		return err
	}

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

// IsPaused checks if the current time falls within the configured pause period
func (c Config) IsPaused() (bool, error) {
	if c.PauseStart == "" || c.PauseEnd == "" {
		return false, nil
	}

	now := time.Now()
	startDate, err := time.Parse("2006-01-02", c.PauseStart)
	if err != nil {
		return false, fmt.Errorf("invalid PauseStart date format '%s', expected YYYY-MM-DD: %w", c.PauseStart, err)
	}

	endDate, err := time.Parse("2006-01-02", c.PauseEnd)
	if err != nil {
		return false, fmt.Errorf("invalid PauseEnd date format '%s', expected YYYY-MM-DD: %w", c.PauseEnd, err)
	}

	// Set time to start of day for start date and end of day for end date
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, now.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, now.Location())

	return (now.After(startDate) || now.Equal(startDate)) && (now.Before(endDate) || now.Equal(endDate)), nil
}
