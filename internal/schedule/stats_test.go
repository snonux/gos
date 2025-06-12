package schedule

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/timestamp"
)

func TestGatherPostedStats(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "gos_stats_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test posted files with different timestamps
	now := time.Now()
	testFiles := []struct {
		filename string
		content  string
		timestamp time.Time
	}{
		{
			// Posted entry from 5 days ago
			filename: "post1.txt." + now.AddDate(0, 0, -5).Format(timestamp.Format) + ".posted",
			content:  "Test post 1",
			timestamp: now.AddDate(0, 0, -5),
		},
		{
			// Posted entry from 3 days ago
			filename: "post2.txt." + now.AddDate(0, 0, -3).Format(timestamp.Format) + ".posted",
			content:  "Test post 2",
			timestamp: now.AddDate(0, 0, -3),
		},
		{
			// Posted entry from 1 day ago
			filename: "post3.txt." + now.AddDate(0, 0, -1).Format(timestamp.Format) + ".posted",
			content:  "Test post 3",
			timestamp: now.AddDate(0, 0, -1),
		},
		{
			// Posted entry from 10 days ago (outside lookback period)
			filename: "old_post.txt." + now.AddDate(0, 0, -10).Format(timestamp.Format) + ".posted",
			content:  "Old test post",
			timestamp: now.AddDate(0, 0, -10),
		},
		{
			// Queued entry (should be ignored)
			filename: "queued_post.txt." + now.AddDate(0, 0, -2).Format(timestamp.Format) + ".queued",
			content:  "Queued post",
			timestamp: now.AddDate(0, 0, -2),
		},
		{
			// Posted entry with .now. tag (should be ignored)
			filename: "now_post.now.txt." + now.AddDate(0, 0, -2).Format(timestamp.Format) + ".posted",
			content:  "Now post",
			timestamp: now.AddDate(0, 0, -2),
		},
	}

	// Create test files
	for _, tf := range testFiles {
		filePath := filepath.Join(tmpDir, tf.filename)
		if err := os.WriteFile(filePath, []byte(tf.content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", tf.filename, err)
		}
	}

	// Initialize stats and run gatherPostedStats
	s := &stats{}
	lookbackTime := now.AddDate(0, 0, -7) // 7 days lookback
	cfg := config.Config{} // Empty config (no pause)

	err = s.gatherPostedStats(tmpDir, lookbackTime, cfg)
	if err != nil {
		t.Fatalf("gatherPostedStats failed: %v", err)
	}

	// Verify results
	expectedPosted := 3 // post1, post2, post3 (excludes old_post, queued_post, now_post)
	if s.posted != expectedPosted {
		t.Errorf("Expected posted=%d, got posted=%d", expectedPosted, s.posted)
	}

	expectedTotalPosted := 5 // post1, post2, post3, old_post, now_post (excludes only queued_post)
	if s.totalPosted != expectedTotalPosted {
		t.Errorf("Expected totalPosted=%d, got totalPosted=%d", expectedTotalPosted, s.totalPosted)
	}

	// Check sinceDays calculation (should be approximately 5 days from oldest to now)
	expectedSinceDays := 5.0
	if s.sinceDays < expectedSinceDays-0.1 || s.sinceDays > expectedSinceDays+0.1 {
		t.Errorf("Expected sinceDays≈%.1f, got sinceDays=%.1f", expectedSinceDays, s.sinceDays)
	}

	// Check postsPerDay calculation
	expectedPostsPerDay := float64(expectedPosted) / s.sinceDays
	if s.postsPerDay < expectedPostsPerDay-0.01 || s.postsPerDay > expectedPostsPerDay+0.01 {
		t.Errorf("Expected postsPerDay≈%.2f, got postsPerDay=%.2f", expectedPostsPerDay, s.postsPerDay)
	}

	// Check lastPostDaysAgo (should be approximately 1 day)
	expectedLastPostDaysAgo := 1.0
	if s.lastPostDaysAgo < expectedLastPostDaysAgo-0.1 || s.lastPostDaysAgo > expectedLastPostDaysAgo+0.1 {
		t.Errorf("Expected lastPostDaysAgo≈%.1f, got lastPostDaysAgo=%.1f", expectedLastPostDaysAgo, s.lastPostDaysAgo)
	}

	// Check totalSinceDays (should be approximately 10 days from oldest to now)
	expectedTotalSinceDays := 10.0
	if s.totalSinceDays < expectedTotalSinceDays-0.1 || s.totalSinceDays > expectedTotalSinceDays+0.1 {
		t.Errorf("Expected totalSinceDays≈%.1f, got totalSinceDays=%.1f", expectedTotalSinceDays, s.totalSinceDays)
	}

	// Check totalPostsPerDay calculation
	expectedTotalPostsPerDay := float64(expectedTotalPosted) / s.totalSinceDays
	if s.totalPostsPerDay < expectedTotalPostsPerDay-0.01 || s.totalPostsPerDay > expectedTotalPostsPerDay+0.01 {
		t.Errorf("Expected totalPostsPerDay≈%.2f, got totalPostsPerDay=%.2f", expectedTotalPostsPerDay, s.totalPostsPerDay)
	}
}

func TestGatherPostedStatsEmptyDir(t *testing.T) {
	// Create a temporary directory with no files
	tmpDir, err := os.MkdirTemp("", "gos_stats_empty_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	s := &stats{}
	lookbackTime := time.Now().AddDate(0, 0, -7)
	cfg := config.Config{} // Empty config (no pause)

	err = s.gatherPostedStats(tmpDir, lookbackTime, cfg)
	if err != nil {
		t.Fatalf("gatherPostedStats failed: %v", err)
	}

	// All stats should be zero or NaN for division by zero
	if s.posted != 0 {
		t.Errorf("Expected posted=0, got posted=%d", s.posted)
	}
	if s.totalPosted != 0 {
		t.Errorf("Expected totalPosted=0, got totalPosted=%d", s.totalPosted)
	}
	// postsPerDay and totalPostsPerDay will be NaN due to division by zero when no posts exist
	// This is the current behavior of the code
}

func TestGatherPostedStatsWithPause(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "gos_stats_pause_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test scenario:
	// - Analysis period: 10 days ago to now
	// - Pause period: 6 days ago to 2 days ago (4 days pause)
	// - Posts: 8 days ago, 7 days ago, 1 day ago (3 posts)
	// - Expected: 3 posts over 6 active days (10 - 4 pause days) = 0.5 posts/day

	now := time.Now()
	testFiles := []struct {
		filename string
		content  string
	}{
		{
			// Posted entry from 8 days ago (before pause)
			filename: "post1.txt." + now.AddDate(0, 0, -8).Format(timestamp.Format) + ".posted",
			content:  "Test post 1",
		},
		{
			// Posted entry from 7 days ago (before pause)
			filename: "post2.txt." + now.AddDate(0, 0, -7).Format(timestamp.Format) + ".posted",
			content:  "Test post 2",
		},
		{
			// Posted entry from 1 day ago (after pause)
			filename: "post3.txt." + now.AddDate(0, 0, -1).Format(timestamp.Format) + ".posted",
			content:  "Test post 3",
		},
	}

	// Create test files
	for _, tf := range testFiles {
		filePath := filepath.Join(tmpDir, tf.filename)
		if err := os.WriteFile(filePath, []byte(tf.content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", tf.filename, err)
		}
	}

	// Configure pause period (6 days ago to 2 days ago)
	pauseStart := now.AddDate(0, 0, -6).Format("2006-01-02")
	pauseEnd := now.AddDate(0, 0, -2).Format("2006-01-02")
	cfg := config.Config{
		PauseStart: pauseStart,
		PauseEnd:   pauseEnd,
	}

	// Initialize stats and run gatherPostedStats
	s := &stats{}
	lookbackTime := now.AddDate(0, 0, -10) // 10 days lookback

	err = s.gatherPostedStats(tmpDir, lookbackTime, cfg)
	if err != nil {
		t.Fatalf("gatherPostedStats failed: %v", err)
	}

	// Verify results
	expectedPosted := 3 // All 3 posts are within lookback period
	if s.posted != expectedPosted {
		t.Errorf("Expected posted=%d, got posted=%d", expectedPosted, s.posted)
	}

	// Check sinceDays should be around 8 days (oldest post to now)
	expectedSinceDays := 8.0
	if s.sinceDays < expectedSinceDays-0.1 || s.sinceDays > expectedSinceDays+0.1 {
		t.Errorf("Expected sinceDays≈%.1f, got sinceDays=%.1f", expectedSinceDays, s.sinceDays)
	}

	// Key test: postsPerDay should exclude the pause period
	// Let's calculate this step by step:
	
	// 1. Find the actual time range of our posts
	oldestPost := now.AddDate(0, 0, -8)
	
	// 2. Calculate paused days using our helper function
	actualPausedDays := calculatePausedDays(oldestPost, now, cfg)
	
	// 3. Calculate expected values
	expectedActiveDays := s.sinceDays - actualPausedDays
	expectedPostsPerDay := float64(expectedPosted) / expectedActiveDays
	
	// Debug info for understanding the calculation
	// t.Logf("Debug: sinceDays=%.1f, actualPausedDays=%.1f, expectedActiveDays=%.1f", 
	//	s.sinceDays, actualPausedDays, expectedActiveDays)
	
	if s.postsPerDay < expectedPostsPerDay-0.01 || s.postsPerDay > expectedPostsPerDay+0.01 {
		t.Errorf("Expected postsPerDay≈%.2f (%.0f posts / %.1f active days), got postsPerDay=%.2f", 
			expectedPostsPerDay, float64(expectedPosted), expectedActiveDays, s.postsPerDay)
	}
}

func TestCalculatePausedDays(t *testing.T) {
	cfg := config.Config{
		PauseStart: "2024-07-01",
		PauseEnd:   "2024-07-10",
	}

	tests := []struct {
		name        string
		startTime   string
		endTime     string
		expected    float64
		description string
	}{
		{
			name:        "No overlap - before pause",
			startTime:   "2024-06-20",
			endTime:     "2024-06-30",
			expected:    0,
			description: "Period entirely before pause",
		},
		{
			name:        "No overlap - after pause",
			startTime:   "2024-07-15",
			endTime:     "2024-07-20",
			expected:    0,
			description: "Period entirely after pause",
		},
		{
			name:        "Full pause period overlap",
			startTime:   "2024-06-25",
			endTime:     "2024-07-15",
			expected:    10,
			description: "Period encompasses entire pause",
		},
		{
			name:        "Partial overlap - start during pause",
			startTime:   "2024-07-05",
			endTime:     "2024-07-15",
			expected:    6, // 2024-07-05 to 2024-07-10 inclusive = 6 days
			description: "Period starts during pause",
		},
		{
			name:        "Partial overlap - end during pause",
			startTime:   "2024-06-25",
			endTime:     "2024-07-05",
			expected:    5,
			description: "Period ends during pause",
		},
		{
			name:        "Exact pause period",
			startTime:   "2024-07-01",
			endTime:     "2024-07-10",
			expected:    10,
			description: "Period exactly matches pause",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime, _ := time.Parse("2006-01-02", tt.startTime)
			endTime, _ := time.Parse("2006-01-02", tt.endTime)
			endTime = endTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second) // End of day

			result := calculatePausedDays(startTime, endTime, cfg)
			if result < tt.expected-0.1 || result > tt.expected+0.1 {
				t.Errorf("Expected %.1f paused days, got %.1f", tt.expected, result)
			}
		})
	}
}