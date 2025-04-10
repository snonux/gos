package internal

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/table"
)

func Main(composeModeDefault bool) {
	dry := flag.Bool("dry", false, "Dry run")
	version := flag.Bool("version", false, "Display version")
	composeMode := flag.Bool("compose", composeModeDefault, "Compose a new entry")
	gosDir := flag.String("gosDir", filepath.Join(os.Getenv("HOME"), ".gosdir"), "Gos' queue and DB directory")
	cacheDir := flag.String("cacheDir", filepath.Join(*gosDir, "cache"), "Go's cache dir")
	browser := flag.String("browser", "firefox", "OAuth2 browser")
	configPath := filepath.Join(os.Getenv("HOME"), ".config/gos/gos.json")
	configPath = *flag.String("configPath", configPath, "Gos' config file path")
	platforms := flag.String("platforms", "Mastodon:500,LinkedIn:1000", "Platforms enabled plus their post size limits")
	target := flag.Int("target", 2, "How many posts per week are the target?")
	minQueued := flag.Int("minQueued", 4, "Minimum of queued items until printing a warn message!")
	maxDaysQueued := flag.Int("maxDaysQueued", 1000, "Maximum days worth of queued posts until target++ and pauseDays--")
	pauseDays := flag.Int("pauseDays", 3, "How many days until next post can be posted?")
	runInterval := flag.Int("runInterval", 12, "How many hours to wait for the next run.")
	lookback := flag.Int("lookback", 30, "How many days look back in time for posting history")
	geminiSummaryFor := flag.String("geminiSummaryFor", "", "Generate a summary in Gemini Gemtext format, format is coma separated string of months, e.g. 202410,202411")
	geminiCapsules := flag.String("geminiCapsules", "foo.zone", "Comma sepaeated list Gemini capsules. Used by geminiEnable to detect Gemtext links")
	gemtexterEnable := flag.Bool("gemtexterEnable", false, "Add special Gemtexter (the static site generator) tags to the Gemini Gemtext summary")
	dev := flag.Bool("dev", false, "For internal development purposes only")
	flag.Parse()

	if *dev {
		table.New().
			Header("foo", "bar", "baz", 3).
			Row("hans", "wurst", "klo", 3.3).
			Row("Klopapier", "Hans", "Wurst", 7).
			MustRender()
		os.Exit(0)
	}

	conf, err := config.New(configPath, *composeMode)
	if err != nil {
		log.Fatal(err)
	}

	args := config.Args{
		DryRun:          *dry,
		GosDir:          *gosDir,
		Target:          *target,
		MinQueued:       *minQueued,
		MaxDaysQueued:   *maxDaysQueued,
		PauseDays:       *pauseDays,
		RunInterval:     time.Duration(*runInterval) * time.Hour, // TODO: Document
		Lookback:        time.Duration(*lookback) * time.Hour * 24,
		ConfigPath:      configPath,
		Config:          conf,
		CacheDir:        *cacheDir,
		OAuth2Browser:   *browser,
		GemtexterEnable: *gemtexterEnable,
		GeminiCapsules:  strings.Split(*geminiCapsules, ","),
		ComposeMode:     *composeMode,
	}
	if *geminiSummaryFor != "" {
		args.GeminiSummaryFor = strings.Split(*geminiSummaryFor, ",")
	}

	if err := args.ParsePlatforms(*platforms); err != nil {
		log.Fatal(err)
	}

	if *version {
		printVersion()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx, args); err != nil {
		log.Fatal(err)
	}
}
