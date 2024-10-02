package schedule

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/format"
	"codeberg.org/snonux/gos/internal/oi"
)

// Posting stats
type stats struct {
	posted            int
	queued            int
	sinceDays         float64
	postsPerDay       float64
	postsPerDayTarget float64
}

func (s stats) String() string {
	return fmt.Sprintf("posted:%d,queued:%d,sinceDays:%v,postsPerDay:%v >? postsPerDayTarget:%v",
		s.posted, s.queued, s.sinceDays, s.postsPerDay, s.postsPerDayTarget,
	)
}

func newStats(dir string, lookback time.Duration, target int) (stats, error) {
	stats := stats{postsPerDayTarget: float64(target) / 7}

	if err := stats.gatherPostedStats(dir, pastTime(lookback)); err != nil {
		return stats, err
	}
	if err := stats.gatherQueuedStats(dir); err != nil {
		return stats, err
	}

	return stats, nil
}

func (s stats) targetHit() bool {
	return s.postsPerDay >= s.postsPerDayTarget
}

func (s *stats) gatherPostedStats(dir string, lookbackTime time.Time) error {
	var (
		now    time.Time = nowTime()
		oldest time.Time = now
	)

	err := oi.TraverseDir(dir, func(file os.DirEntry) error {
		filePath := filepath.Join(dir, file.Name())
		ent, err := entry.New(filePath)
		if err != nil {
			return err
		}
		if ent.State != entry.Posted || ent.Time.Before(lookbackTime) {
			return nil
		}
		if ent.Time.Before(oldest) {
			oldest = ent.Time
		}
		s.posted++
		return nil
	})
	if err != nil {
		return err
	}

	since := now.Sub(oldest)
	s.sinceDays = since.Abs().Hours() / 24
	s.postsPerDay = float64(s.posted) / float64(s.sinceDays)
	return nil
}

func (s *stats) gatherQueuedStats(dir string) error {
	var firstQueuedPath string

	err := oi.TraverseDir(dir, func(file os.DirEntry) error {
		filePath := filepath.Join(dir, file.Name())
		ent, err := entry.New(filePath)
		if err != nil {
			return err
		}
		if ent.State == entry.Queued {
			if firstQueuedPath == "" {
				firstQueuedPath = filePath
			}
			s.queued++
		}
		return nil
	})

	return err
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
