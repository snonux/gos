package schedule

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/oi"
	"codeberg.org/snonux/gos/internal/timestamp"
)

// Posting stats
type stats struct {
	posted            int
	queued            int
	sinceDays         float64
	postsPerDay       float64
	postsPerDayTarget float64
	lastPostDaysAgo   float64
}

func (s stats) String() string {
	return fmt.Sprintf("posted:%d,queued:%d,sinceDays:%v,postsPerDayTarget:%v>?%v,lastPostDaysAgo:%v",
		s.posted, s.queued, s.sinceDays, s.postsPerDay, s.postsPerDayTarget, s.lastPostDaysAgo,
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

func (s stats) targetHit(pauseDays int) bool {
	if s.postsPerDay >= s.postsPerDayTarget {
		log.Println("Posts per day target hit")
		return true
	}
	if s.lastPostDaysAgo <= float64(pauseDays) {
		log.Println("Need to wait a bit longer as last post isn't", pauseDays, "ago yet")
		return true

	}
	return false
}

func (s *stats) gatherPostedStats(dir string, lookbackTime time.Time) error {
	var (
		now    time.Time = timestamp.NowTime()
		oldest time.Time = now
		newest time.Time = timestamp.OldestValidTime()
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
		if ent.Time.After(newest) {
			newest = ent.Time
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
	s.lastPostDaysAgo = now.Sub(newest).Hours() / 24.0

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

func pastTime(duration time.Duration) time.Time {
	return timestamp.NowTime().Add(-duration)
}
