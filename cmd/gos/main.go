package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
	platforms := flag.String("platforms", "Mastodon,LinkedIn", "Platforms enabled")
	lookback := flag.Int("lookback", 30, "How many days look back in time for posting history")
	flag.Parse()

	args := config.Args{
		DryRun:    *dry,
		GosDir:    *gosDir,
		Platforms: strings.Split(*platforms, ","),
		Lookback:  time.Duration(*lookback) * time.Hour * 24,
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
