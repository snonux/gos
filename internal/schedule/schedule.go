package schedule

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/oi"
)

var (
	ErrNothingToSchedule = errors.New("nothing to schedule")
	ErrNothingQueued     = errors.New("nothing queued")
)

func Run(args config.Args, platform string) (entry.Entry, error) {
	dir := fmt.Sprintf("%s/db/platforms/%s", args.GosDir, strings.ToLower(platform))
	stats, err := newStats(dir, args.Lookback, args.Target)
	if err != nil {
		return entry.Zero, err
	}

	log.Println("For", platform, "stats:", stats)
	if stats.targetHit() {
		log.Println("Target hit, not posting at", platform)
		return entry.Zero, ErrNothingToSchedule
	}

	// Schedule random qeued entry for platform
	ent, err := oi.ReadDirRandom(dir, func(file os.DirEntry) (entry.Entry, bool) {
		ent, err := entry.New(filepath.Join(dir, file.Name()))
		if err != nil {
			log.Println(err)
			return entry.Zero, false
		}
		return ent, ent.State == entry.Queued
	})

	if err != nil {
		return entry.Zero, fmt.Errorf("%w: %w", ErrNothingQueued, err)
	}
	return ent, nil
}
