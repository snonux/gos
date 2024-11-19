package config

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/platforms"
)

var validPlatforms = []string{"mastodon", "linkedin"}

type Args struct {
	GosDir            string
	CacheDir          string
	DryRun            bool
	Platforms         map[platforms.Platform]int // Platform and post size limits
	Target            int
	MinQueued         int
	MaxDaysQueued     int
	PauseDays         int
	Lookback          time.Duration
	SecretsConfigPath string
	Secrets           Secrets
	OAuth2Browser     string
}

func (a *Args) ParsePlatforms(platformStrs string) error {
	a.Platforms = make(map[platforms.Platform]int)

	for _, platformStr := range strings.Split(platformStrs, ",") {
		// E.g. Mastodon:500
		parts := strings.Split(platformStr, ":")
		platform, err := platforms.New(parts[0])
		if err != nil {
			return err
		}

		// E.g. args.Platform["mastodon"] = 500
		if len(parts) > 1 {
			a.Platforms[platform], err = strconv.Atoi(parts[1])
			if err != nil {
				return err
			}
		} else {
			colour.Infoln("No message length specified for", platform, "so assuming 500")
			a.Platforms[platform] = 500
		}
	}
	return nil
}

func (a *Args) Validate() error {
	for platform := range a.Platforms {
		if !slices.Contains(validPlatforms, platform.String()) {
			return fmt.Errorf("Platform %s not supported", platform)
		}
	}
	return nil
}
