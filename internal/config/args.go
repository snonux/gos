package config

import (
	"strconv"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/colour"
)

type Args struct {
	GosDir            string
	CacheDir          string
	DryRun            bool
	Platforms         map[string]int // Platform and post size limits
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
	a.Platforms = make(map[string]int)

	for _, platformInfo := range strings.Split(platformStrs, ",") {
		// E.g. Mastodon:500
		parts := strings.Split(platformInfo, ":")
		platformStr := parts[0]

		// E.g. args.Platform["mastodon"] = 500
		if len(parts) > 1 {
			var err error
			a.Platforms[platformStr], err = strconv.Atoi(parts[1])
			if err != nil {
				return err
			}
		} else {
			colour.Infoln("No message length specified for", platformStr, "so assuming 500")
			a.Platforms[platformStr] = 500
		}
	}
	return nil
}
