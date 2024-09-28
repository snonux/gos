package schedule

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/oi"
)

var (
	ErrNothingToSchedule = errors.New("nothing to schedule")
	ErrNothingQueued     = errors.New("nothing queued")
)

func Run(args config.Args, platform string) (string, error) {
	dir := fmt.Sprintf("%s/db/platforms/%s", args.GosDir, strings.ToLower(platform))
	stats, err := newStats(dir, args.Lookback, args.Target)
	if err != nil {
		return "", err
	}

	log.Println("For", platform, "stats:", stats)
	if stats.targetHit() {
		log.Println("Target hit, not posting at", platform)
		return "", ErrNothingToSchedule
	}

	// Schedule random qeued entry for platform
	randomEntry, err := oi.ReadDirRandomEntry(dir, func(file os.DirEntry) bool {
		return strings.HasSuffix(file.Name(), ".queued")
	})

	if err != nil {
		// TODO: FIX THIS
		return randomEntry, fmt.Errorf("%w: %w", ErrNothingQueued, err)
	}
	return randomEntry, nil
}
