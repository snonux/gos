package schedule

import (
	"fmt"
	"os"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/format"
	"codeberg.org/snonux/gos/internal/oi"
)

// Posting stats
type stats struct {
	posted      int
	queued      int
	sinceDays   float64
	postsPerDay float64
}

func (s stats) String() string {
	return fmt.Sprintf("posted:%d,queued:%d,sinceDays:%v,postsPerDay:%v",
		s.posted, s.queued, s.sinceDays, s.postsPerDay,
	)
}

func newStats(dir string) (stats, error) {
	var stats stats

	if err := stats.gatherPostedStats(dir); err != nil {
		return stats, err
	}
	if err := stats.gatherQueuedStats(dir); err != nil {
		return stats, err
	}

	return stats, nil
}

func (s *stats) gatherPostedStats(dir string) error {
	ch, err := oi.ReadDirFilter(dir, func(file os.DirEntry) bool {
		return strings.HasSuffix(file.Name(), ".posted")
	})
	if err != nil {
		return err
	}

	var (
		now    time.Time = nowTime()
		oldest time.Time = now
	)

	for filePath := range ch {
		newOldest, err := parseEntryPath(filePath)
		if err != nil {
			return err
		}
		if newOldest.Before(oldest) {
			oldest = newOldest
		}
		s.posted++
	}

	since := now.Sub(oldest)
	s.sinceDays = since.Abs().Hours() / 24
	s.postsPerDay = float64(s.posted) / s.sinceDays
	return nil
}

func (s *stats) gatherQueuedStats(dir string) error {
	ch, err := oi.ReadDirFilter(dir, func(file os.DirEntry) bool {
		return strings.HasSuffix(file.Name(), ".queued")
	})
	if err != nil {
		return err
	}

	var firstQueuedPath string
	for filePath := range ch {
		if _, err := parseEntryPath(filePath); err != nil {
			return err
		}
		if firstQueuedPath == "" {
			firstQueuedPath = filePath
		}
		s.queued++
	}

	return nil
}

// Make a simpler "now" time which gets rid of any extra information like offsets etc.
func nowTime() time.Time {
	simplerNow, err := time.Parse(format.Time, time.Now().Format(format.Time))
	if err != nil {
		panic(err)
	}
	return simplerNow
}

func parseEntryPath(filePath string) (time.Time, error) {
	// Format: foobarbaz.something.here.txt.STAMP.{posted,queued}
	// We want to get the STAMP!
	parts := strings.Split(filePath, ".")
	if len(parts) < 4 {
		return time.Time{}, fmt.Errorf("not a valid entry path: %s", filePath)
	}
	return time.Parse(format.Time, parts[len(parts)-2])
}
