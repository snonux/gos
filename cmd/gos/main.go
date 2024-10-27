package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal"
	"codeberg.org/snonux/gos/internal/config"
)

const versionStr = "v0.0.1"

// TODO: edit tag, to edit post before it is queued.
// TODO: now tag, to post a post immediately, ignoring the stats.
func main() {
	dry := flag.Bool("dry", false, "Dry run")
	version := flag.Bool("version", false, "Display version")
	gosDir := flag.String("gosDir", filepath.Join(os.Getenv("HOME"), ".gosdir"), "Gos' queue and DB directory")
	browser := flag.String("browser", "firefox", "OAuth2 browser")
	secretsConfigPath := filepath.Join(os.Getenv("HOME"), ".config/gos/gosec.json")
	secretsConfigPath = *flag.String("secretsConfig", secretsConfigPath, "Gos' secret config")
	platforms := flag.String("platforms", "Mastodon:500,LinkedIn:1000", "Platforms enabled plus their post size limits")
	target := flag.Int("target", 2, "How many posts per week are the target?")
	pauseDays := flag.Int("pauseDays", 3, "How many days until next post can be posted?")
	lookback := flag.Int("lookback", 30, "How many days look back in time for posting history")
	flag.Parse()

	secrets, err := config.NewSecrets(secretsConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	args := config.Args{
		DryRun:            *dry,
		GosDir:            *gosDir,
		Platforms:         make(map[string]int),
		Target:            *target,
		PauseDays:         *pauseDays,
		Lookback:          time.Duration(*lookback) * time.Hour * 24,
		SecretsConfigPath: secretsConfigPath,
		Secrets:           secrets,
		OAuth2Browser:     *browser,
	}
	for _, platform := range strings.Split(*platforms, ",") {
		// E.g. Mastodon:500
		parts := strings.Split(platform, ":")
		var err error
		// E.g. args.Platform["mastodon"] = 500
		if len(parts) > 1 {
			args.Platforms[parts[0]], err = strconv.Atoi(parts[1])
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Println("No message length specified for", platform, "so assuming 500")
			args.Platforms[parts[0]] = 500
		}
	}

	if err := args.Validate(); err != nil {
		log.Fatal(err)
	}

	if *version {
		fmt.Printf("This is Gos version %s; (C) by Paul Buetow\n", versionStr)
		fmt.Println("https://codeberg.org/snonux/gos")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := internal.Run(ctx, args); err != nil {
		log.Fatal(err)
	}
}
