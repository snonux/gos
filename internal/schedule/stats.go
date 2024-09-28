package schedule

import (
	"errors"
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

func newStats(dir string, lookback time.Duration) (stats, error) {
	var stats stats

	if err := stats.gatherPostedStats(dir, pastTime(lookback)); err != nil {
		return stats, err
	}
	if err := stats.gatherQueuedStats(dir); err != nil {
		return stats, err
	}

	return stats, nil
}

func (s *stats) gatherPostedStats(dir string, lookbackTime time.Time) error {
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

	var errs []error
	// TODO: Maybe refactor to include in ReadDirFilter filter
	for filePath := range ch {
		entryTime, err := parseEntryPath(filePath)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if entryTime.Before(lookbackTime) {
			continue
		}
		if entryTime.Before(oldest) {
			oldest = entryTime
		}
		s.posted++
	}

	since := now.Sub(oldest)
	s.sinceDays = since.Abs().Hours() / 24
	s.postsPerDay = float64(s.posted) / s.sinceDays
	return errors.Join(errs...)
}

func (s *stats) gatherQueuedStats(dir string) error {
	ch, err := oi.ReadDirFilter(dir, func(file os.DirEntry) bool {
		return strings.HasSuffix(file.Name(), ".queued")
	})
	if err != nil {
		return err
	}

	var (
		firstQueuedPath string
		errs            []error
	)
	for filePath := range ch {
		if _, err := parseEntryPath(filePath); err != nil {
			errs = append(errs, err)
			continue
		}
		if firstQueuedPath == "" {
			firstQueuedPath = filePath
		}
		s.queued++
	}

	return errors.Join(errs...)
}

// Make a simpler "now" time which gets rid of any extra information like offsets etc.
func nowTime() time.Time {
	simplerNow, err := time.Parse(format.Time, time.Now().Format(format.Time))
	if err != nil {
		panic(err)
	}
	return simplerNow
}

func pastTime(duration time.Duration) time.Time {
	return nowTime().Add(-duration)
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
