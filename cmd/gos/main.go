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

const versionStr = "v0.0.2"

func main() {
	dry := flag.Bool("dry", false, "Dry run")
	version := flag.Bool("version", false, "Display version")
	gosDir := flag.String("gosDir", filepath.Join(os.Getenv("HOME"), ".gosdir"), "Gos' queue and DB directory")
	cacheDir := flag.String("cacheDir", filepath.Join(*gosDir, "cache"), "Go's cache dir")
	browser := flag.String("browser", "firefox", "OAuth2 browser")
	secretsConfigPath := filepath.Join(os.Getenv("HOME"), ".config/gos/gosec.json")
	secretsConfigPath = *flag.String("secretsConfig", secretsConfigPath, "Gos' secret config")
	platforms := flag.String("platforms", "Mastodon:500,LinkedIn:1000", "Platforms enabled plus their post size limits")
	target := flag.Int("target", 2, "How many posts per week are the target?")
	minQueued := flag.Int("minQueued", 4, "Minimum of queued items until printing a warn message!")
	maxDaysQueued := flag.Int("maxDaysQueued", 365, "Maximum days worth of queued posts until target++ and pauseDays--")
	pauseDays := flag.Int("pauseDays", 3, "How many days until next post can be posted?")
	lookback := flag.Int("lookback", 30, "How many days look back in time for posting history")
	summaryFor := flag.String("summaryFor", "", "Generate a summary in Gemtext format, format is coma separated string of months, e.g. 202410,202411")
	gemtexterEnable := flag.Bool("gemtexterEnable", true, "Add special Gemtexter tags to the Gemtext summary")
	flag.Parse()

	secrets, err := config.NewSecrets(secretsConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	args := config.Args{
		DryRun:            *dry,
		GosDir:            *gosDir,
		Target:            *target,
		MinQueued:         *minQueued,
		MaxDaysQueued:     *maxDaysQueued,
		PauseDays:         *pauseDays,
		Lookback:          time.Duration(*lookback) * time.Hour * 24,
		SecretsConfigPath: secretsConfigPath,
		CacheDir:          *cacheDir,
		Secrets:           secrets,
		OAuth2Browser:     *browser,
		GemtexterEnable:   *gemtexterEnable,
	}
	if *summaryFor != "" {
		args.SummaryFor = strings.Split(*summaryFor, ",")
	}

	if err := args.ParsePlatforms(*platforms); err != nil {
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
