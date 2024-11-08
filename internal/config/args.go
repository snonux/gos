package config

import (
	"fmt"
	"log"
	"slices"
	"strconv"
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
	MaxDaysQueued     int
	PauseDays         int
	Lookback          time.Duration
	CacheDir          string
	SecretsConfigPath string
	Secrets           Secrets
	OAuth2Browser     string
}

func (a *Args) ParsePlatforms(platforms string) error {
	for _, platform := range strings.Split(platforms, ",") {
		// E.g. Mastodon:500
		parts := strings.Split(platform, ":")
		var err error
		// E.g. args.Platform["mastodon"] = 500
		if len(parts) > 1 {
			a.Platforms[parts[0]], err = strconv.Atoi(parts[1])
			if err != nil {
				return err
			}
		} else {
			log.Println("No message length specified for", platform, "so assuming 500")
			a.Platforms[parts[0]] = 500
		}
	}
	return nil
}

func (a *Args) Validate() error {
	for platform := range a.Platforms {
		if !slices.Contains(validPlatforms, strings.ToLower(platform)) {
			return fmt.Errorf("Platform %s not supported", platform)
		}
	}
	return nil
}
