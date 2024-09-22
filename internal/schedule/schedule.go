package schedule

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
)

var NothingToSchedule = errors.New("nothing to schedule")

func Run(args config.Args, platform string) (string, error) {
	dir := fmt.Sprintf("%s/db/platforms/%s", args.GosDir, strings.ToLower(platform))
	stats, err := newStats(dir)
	if err != nil {
		return "", err
	}

	log.Println("For", platform, "stats:", stats)

	return "", NothingToSchedule
}
