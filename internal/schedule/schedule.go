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

var NothingToSchedule = errors.New("nothing to schedule")

func Run(args config.Args, platform string) (string, error) {
	dir := fmt.Sprintf("%s/db/platforms/%s", args.GosDir, strings.ToLower(platform))
	stats, err := newStats(dir, args.Lookback, args.Target)
	if err != nil {
		return "", err
	}

	log.Println("For", platform, "stats:", stats)
	if stats.targetHit() {
		log.Println("Target hit, not posting at", platform)
		return "", NothingToSchedule
	}

	// Schedule random qeued entry for platform
	return oi.ReadDirRandomEntry(dir, func(file os.DirEntry) bool {
		return strings.HasSuffix(file.Name(), ".queued")
	})
}
