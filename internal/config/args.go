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
	Platforms         []string
	Target            int
	Lookback          time.Duration
	SecretsConfigPath string
	Secrets           Secrets
}

func (a Args) Validate() error {
	for _, platform := range a.Platforms {
		if !slices.Contains(validPlatforms, strings.ToLower(platform)) {
			return fmt.Errorf("Platform %s not supported", platform)
		}
	}
	return nil
}
