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

func (s stats) Render(platform string) {
	var sb strings.Builder

	sep := colour.SInfo2f("+%s+%s+", strings.Repeat("-", 22), strings.Repeat("-", 13))
	sb.WriteString(sep)
	sb.WriteString("\n")
	sb.WriteString(colour.SInfo2f("| %-20s | %-11s |", platform, "Stat. value"))
	sb.WriteString("\n")
	sb.WriteString(sep)
	sb.WriteString("\n")
	sb.WriteString(colour.SInfo2f("| %-20s | %-11s |", "Stats since (days)", fmt.Sprintf("%.02f", s.sinceDays)))
	sb.WriteString("\n")
	sb.WriteString(colour.SInfo2f("| %-20s | %-11s |", "#Posted entries", strconv.Itoa(s.posted)))
	sb.WriteString("\n")
	sb.WriteString(colour.SInfo2f("| %-20s | %-11s |", "#Queued entries", strconv.Itoa(s.queued)))
	sb.WriteString("\n")
	sb.WriteString(colour.SInfo2f("| %-20s | %-11s |", "Enough for (days)", fmt.Sprintf("%.02f", s.queuedForDays)))
	sb.WriteString("\n")
	sb.WriteString(colour.SInfo2f("| %-20s | %-11s |", "Last post (days ago)", fmt.Sprintf("%.02f", s.lastPostDaysAgo)))
	sb.WriteString("\n")
	sb.WriteString(colour.SInfo2f("| %-20s | %-11s |", "Posts per day", fmt.Sprintf("%.02f", s.postsPerDay)))
	sb.WriteString("\n")
	sb.WriteString(colour.SInfo2f("| %-20s | %-11s |", "Posts per day target", fmt.Sprintf("%.02f", s.postsPerDayTarget)))
	sb.WriteString("\n")
	sb.WriteString(sep)
	sb.WriteString("\n")

	fmt.Print(sb.String())
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

	s.queuedForDays = float64(s.queued) / s.postsPerDayTarget

	return err
}

func pastTime(duration time.Duration) time.Time {
	return timestamp.NowTime().Add(-duration)
}
