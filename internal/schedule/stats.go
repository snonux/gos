package schedule

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/oi"
	"codeberg.org/snonux/gos/internal/platforms"
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
}

func newStats(dir string, lookback time.Duration, target int) (stats, error) {
	s := stats{postsPerDayTarget: float64(target) / 7}

	if err := s.gatherPostedStats(dir, pastTime(lookback)); err != nil {
		return s, err
	}
	if err := s.gatherQueuedStats(dir); err != nil {
		return s, err
	}

	return s, nil
}

func (s stats) String() string {
	return fmt.Sprintf("posted:%d,queued:%d,sinceDays:%v,postsPerDayTarget:%v>?%v,lastPostDaysAgo:%v",
		s.posted, s.queued, s.sinceDays, s.postsPerDay, s.postsPerDayTarget, s.lastPostDaysAgo,
	)
}

func (s stats) targetHit(pauseDays, maxQueuedDays int) bool {
	if s.queuedForDays > float64(maxQueuedDays) {
		s.postsPerDayTarget++
		pauseDays--
	}
	if s.postsPerDay >= s.postsPerDayTarget {
		colour.Infoln("Posts per day target hit")
		return true
	}
	if s.lastPostDaysAgo <= float64(pauseDays) {
		colour.Infoln("Need to wait a bit longer as last post isn't", pauseDays, "days ago yet")
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
	var sb strings.Builder

	dataRow := func(descr1, val1, descr2, val2 string) {
		const format = "| %-21s | %-11s | %-21s | %-11s |"
		sb.WriteString(colour.SInfo2f(format, descr1, val1, descr2, val2))
		sb.WriteString("\n")
	}

	sep := colour.SInfo2f("+%s+%s+%s+%s+", strings.Repeat("-", 23),
		strings.Repeat("-", 13), strings.Repeat("-", 23), strings.Repeat("-", 13))

	separator := func() {
		sb.WriteString(sep)
		sb.WriteString("\n")
	}

	val := func(val any) string {
		switch v := val.(type) {
		case int:
			return strconv.Itoa(v)
		case float64:
			return fmt.Sprintf("%0.2f", v)
		default:
			panic("unexpeced type")
		}
	}

	separator()
	dataRow(platform.String(), "value", "Lifetime stats", "value")
	separator()
	dataRow("Since (days)", val(s.sinceDays), "Total since (days)", val(s.totalSinceDays))
	dataRow("#Posted entries", val(s.posted), "#Total posted entries", val(s.totalPosted))
	dataRow("#Queued entries", val(s.queued), "", "")
	dataRow("Enough for (days)", val(s.queuedForDays), "", "")
	dataRow("Last post (days ago)", val(s.lastPostDaysAgo), "", "")
	dataRow("Posts per day", val(s.postsPerDay), "Total posts per day", val(s.totalPostsPerDay))
	dataRow("Posts per day target", val(s.postsPerDayTarget), "", "")
	separator()

	fmt.Print(sb.String())
}

func pastTime(duration time.Duration) time.Time {
	return timestamp.NowTime().Add(-duration)
}
