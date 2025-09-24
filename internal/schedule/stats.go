package schedule

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/config"
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

func newStats(gosDir string, platformName string, lookback time.Duration, target, pauseDays, maxQueuedDays int, cfg config.Config) (stats, error) {
	dir := filepath.Join(gosDir, "db", "platforms", strings.ToLower(platformName))
	s := stats{postsPerDayTarget: float64(target) / 7, pauseDays: pauseDays}

	if err := s.gatherPostedStats(dir, pastTime(lookback), cfg); err != nil {
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

		colour.Infoln(platformName, "- Increasing posts per day target", s.postsPerDayTarget, "by", add, "to", newTarget)
		s.postsPerDayTarget = newTarget

		colour.Infoln(platformName, "- Decreasing pause days from", s.pauseDays, "to", s.pauseDays-1)
		s.pauseDays--
	}

	return s, nil
}

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

func (s *stats) gatherPostedStats(dir string, lookbackTime time.Time, cfg config.Config) error {
	var (
		now         time.Time = timestamp.NowTime()
		newest      time.Time = timestamp.OldestValidTime()
		oldest      time.Time = now // Oldest since lookbackTime
		totalOldest time.Time = now // All time oldest
	)

	err := oi.ForeachDirEntry(dir, func(file os.DirEntry) error {
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

	// Subtract paused days from the calculation period
	pausedDays := calculatePausedDays(oldest, now, cfg)
	activeDays := s.sinceDays - pausedDays
	if activeDays > 0 {
		s.postsPerDay = float64(s.posted) / activeDays
	} else {
		s.postsPerDay = 0
	}
	s.lastPostDaysAgo = now.Sub(newest).Hours() / 24.0

	since = now.Sub(totalOldest)
	s.totalSinceDays = since.Abs().Hours() / 24.0
	totalPausedDays := calculatePausedDays(totalOldest, now, cfg)
	totalActiveDays := s.totalSinceDays - totalPausedDays
	if totalActiveDays > 0 {
		s.totalPostsPerDay = float64(s.totalPosted) / totalActiveDays
	} else {
		s.totalPostsPerDay = 0
	}

	return nil
}

func (s *stats) gatherQueuedStats(dir string) error {
	var firstQueuedPath string

	err := oi.ForeachDirEntry(dir, func(file os.DirEntry) error {
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
	table.New().
		WithColor(colour.AttentionCol).
		Header(platform.String(), "value", "Lifetime stats", "value").
		Row("Since (days)", s.sinceDays, "Total since (days)", s.totalSinceDays).
		Row("#Posted entries", s.posted, "#Total posted entries", s.totalPosted).
		Row("#Queued entries", s.queued, "", "").
		Row("Enough for (days)", s.queuedForDays, "", "").
		Row("Last post (days ago)", s.lastPostDaysAgo, "Pause days", s.pauseDays).
		Row("Posts per day", s.postsPerDay, "Total posts per day", s.totalPostsPerDay).
		Row("Posts per day target", s.postsPerDayTarget, "", "").
		MustRender()
}

func pastTime(duration time.Duration) time.Time {
	return timestamp.NowTime().Add(-duration)
}

// calculatePausedDays calculates the number of days that fall within pause periods
// between startTime and endTime
func calculatePausedDays(startTime, endTime time.Time, cfg config.Config) float64 {
	if cfg.PauseStart == "" || cfg.PauseEnd == "" {
		return 0
	}

	pauseStart, err := time.Parse("2006-01-02", cfg.PauseStart)
	if err != nil {
		return 0 // If parse fails, assume no pause
	}

	pauseEnd, err := time.Parse("2006-01-02", cfg.PauseEnd)
	if err != nil {
		return 0 // If parse fails, assume no pause
	}

	// Set to start and end of day
	pauseStart = time.Date(pauseStart.Year(), pauseStart.Month(), pauseStart.Day(), 0, 0, 0, 0, startTime.Location())
	pauseEnd = time.Date(pauseEnd.Year(), pauseEnd.Month(), pauseEnd.Day(), 23, 59, 59, 999999999, endTime.Location())

	// Find intersection of [startTime, endTime] with [pauseStart, pauseEnd]
	intersectionStart := maxTime(startTime, pauseStart)
	intersectionEnd := minTime(endTime, pauseEnd)

	if intersectionStart.Before(intersectionEnd) || intersectionStart.Equal(intersectionEnd) {
		return intersectionEnd.Sub(intersectionStart).Hours() / 24.0
	}

	return 0
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func PrintAllStats(args config.Args) {
	for platformName := range args.Platforms {
		platform, err := platforms.New(platformName)
		if err != nil {
			colour.Warnln("Error creating platform for", platformName, ":", err)
			continue
		}
		s, err := newStats(args.GosDir, platformName, args.Lookback, args.Target, args.PauseDays, args.MaxDaysQueued, args.Config)
		if err != nil {
			colour.Warnln("Error gathering stats for", platformName, ":", err)
			continue
		}
		s.RenderTable(platform)
	}
}
