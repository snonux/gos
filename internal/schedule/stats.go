package schedule

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/oi"
	"codeberg.org/snonux/gos/internal/platforms"
	"codeberg.org/snonux/gos/internal/table"
	"codeberg.org/snonux/gos/internal/timestamp"
)

// Posting stats
type stats struct {
	posted            int
	queued            int
	queuedForDays     float64
	sinceDays         float64
	postsPerDay       float64
	postsPerDayTarget float64
	lastPostDaysAgo   float64

	totalPosted      int
	totalSinceDays   float64
	totalPostsPerDay float64

	pauseDays int
}

func newStats(dir string, lookback time.Duration, target, pauseDays, maxQueuedDays int) (stats, error) {
	s := stats{postsPerDayTarget: float64(target) / 7, pauseDays: pauseDays}

	if err := s.gatherPostedStats(dir, pastTime(lookback)); err != nil {
		return s, err
	}
	if err := s.gatherQueuedStats(dir); err != nil {
		return s, err
	}

	// Dynamically increase the target when there are many entries queued.
	if s.queuedForDays > float64(maxQueuedDays) {
		add := (s.queuedForDays - float64(maxQueuedDays)) * 0.01
		if add > 0.5 {
			add = 0.5
		}
		newTarget := s.postsPerDayTarget + add

		colour.Infoln("Increasing posts per day target", s.postsPerDayTarget, "by", add, "to", newTarget)
		s.postsPerDayTarget = newTarget

		colour.Infoln("Decreasing pause days from", s.pauseDays, "to", s.pauseDays-1)
		s.pauseDays--
	}

	return s, nil
}

// func (s stats) String() string {
// 	return fmt.Sprintf("posted:%d,queued:%d,sinceDays:%v,postsPerDayTarget:%v>?%v,lastPostDaysAgo:%v",
// 		s.posted, s.queued, s.sinceDays, s.postsPerDay, s.postsPerDayTarget, s.lastPostDaysAgo,
// 	)
// }

func (s stats) targetHit() bool {
	if s.postsPerDay >= s.postsPerDayTarget {
		colour.Infoln("Posts per day target hit", s.postsPerDay, "is greater or equal than", s.postsPerDayTarget)
		return true
	}
	if s.lastPostDaysAgo <= float64(s.pauseDays) {
		colour.Infoln("Need to wait a bit longer as last post isn't", s.pauseDays, "days ago yet")
		return true
	}
	return false
}

func (s *stats) gatherPostedStats(dir string, lookbackTime time.Time) error {
	var (
		now         time.Time = timestamp.NowTime()
		newest      time.Time = timestamp.OldestValidTime()
		oldest      time.Time = now // Oldest since lookbackTime
		totalOldest time.Time = now // All time oldest
	)

	err := oi.TraverseDir(dir, func(file os.DirEntry) error {
		filePath := filepath.Join(dir, file.Name())
		ent, err := entry.New(filePath)
		if err != nil {
			return err
		}
		if ent.State != entry.Posted {
			return nil
		}
		if ent.Time.Before(totalOldest) {
			totalOldest = ent.Time
		}
		s.totalPosted++
		if ent.Time.Before(lookbackTime) {
			return nil
		}
		// Ignore .now.
		if strings.Contains(file.Name(), ".now.") {
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
	s.sinceDays = since.Abs().Hours() / 24.0
	s.postsPerDay = float64(s.posted) / float64(s.sinceDays)
	s.lastPostDaysAgo = now.Sub(newest).Hours() / 24.0

	since = now.Sub(totalOldest)
	s.totalSinceDays = since.Abs().Hours() / 24.0
	s.totalPostsPerDay = float64(s.totalPosted) / float64(s.totalSinceDays)

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

	s.queuedForDays = float64(s.queued) / s.postsPerDayTarget

	return err
}

func (s stats) RenderTable(platform platforms.Platform) {
	err := table.New(platform.String(), "value", "Lifetime stats", "value").
		WithColor(colour.Info2Col).
		Add("Since (days)", s.sinceDays, "Total since (days)", s.totalSinceDays).
		Add("#Posted entries", s.posted, "#Total posted entries", s.totalPosted).
		Add("#Queued entries", s.queued, "", "").
		Add("Enough for (days)", s.queuedForDays, "", "").
		Add("Last post (days ago)", s.lastPostDaysAgo, "Pause days", s.pauseDays).
		Add("Posts per day", s.postsPerDay, "Total posts per day", s.totalPostsPerDay).
		Add("Posts per day target", s.postsPerDayTarget, "", "").
		Render()

	if err != nil {
		panic(err)
	}
}

func pastTime(duration time.Duration) time.Time {
	return timestamp.NowTime().Add(-duration)
}
