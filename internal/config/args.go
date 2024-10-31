package config

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

var validPlatforms = []string{"mastodon", "linkedin"}

type Args struct {
	GosDir            string
	DryRun            bool
	Platforms         map[string]int // Platform name and post size limits
	Target            int
	MinQueued         int
	PauseDays         int
	Lookback          time.Duration
	SecretsConfigPath string
	Secrets           Secrets
	OAuth2Browser     string
}

func (a Args) Validate() error {
	for platform := range a.Platforms {
		if !slices.Contains(validPlatforms, strings.ToLower(platform)) {
			return fmt.Errorf("Platform %s not supported", platform)
		}
	}
	return nil
}
