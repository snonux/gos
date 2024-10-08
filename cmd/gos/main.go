package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal"
	"codeberg.org/snonux/gos/internal/config"
)

const versionStr = "v0.0.0"

func main() {
	dry := flag.Bool("dry", false, "Dry run")
	version := flag.Bool("version", false, "Display version")
	gosDir := flag.String("gosDir", "./gosdir", "Gos' directory")
	secretsConfigPath := filepath.Join(os.Getenv("HOME"), ".config/gos/gosec.json")
	secretsConfigPath = *flag.String("secretsConfig", secretsConfigPath, "Gos' secret config")
	platforms := flag.String("platforms", "Mastodon,LinkedIn", "Platforms enabled")
	target := flag.Int("target", 2, "How many posts per week are the target?")
	lookback := flag.Int("lookback", 30, "How many days look back in time for posting history")
	flag.Parse()

	secrets, err := config.NewSecrets(secretsConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	args := config.Args{
		DryRun:            *dry,
		GosDir:            *gosDir,
		Platforms:         strings.Split(*platforms, ","),
		Target:            *target,
		Lookback:          time.Duration(*lookback) * time.Hour * 24,
		SecretsConfigPath: secretsConfigPath,
		Secrets:           secrets,
	}

	if err := args.Validate(); err != nil {
		log.Fatal(err)
	}

	if *version {
		fmt.Printf("This is Gos version %s; (C) by Paul Buetow\n", versionStr)
		fmt.Println("https://codeberg.org/snonux/gos")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1*time.Minute))
	defer cancel()

	if err := internal.Run(ctx, args); err != nil {
		log.Fatal(err)
	}
}
