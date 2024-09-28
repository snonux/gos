package config

import "time"

type Args struct {
	GosDir    string
	DryRun    bool
	Platforms []string
	Lookback  time.Duration
}
