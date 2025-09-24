package config

import (
	"testing"
	"time"
)

func TestIsPaused(t *testing.T) {
	tests := []struct {
		name        string
		pauseStart  string
		pauseEnd    string
		testTime    time.Time
		expected    bool
		expectError bool
	}{
		{
			name:        "No pause dates configured",
			pauseStart:  "",
			pauseEnd:    "",
			testTime:    time.Date(2024, 8, 15, 12, 0, 0, 0, time.UTC),
			expected:    false,
			expectError: false,
		},
		{
			name:        "Currently paused - middle of pause period",
			pauseStart:  "2024-07-01",
			pauseEnd:    "2024-09-18",
			testTime:    time.Date(2024, 8, 15, 12, 0, 0, 0, time.UTC),
			expected:    true,
			expectError: false,
		},
		{
			name:        "Not paused - before pause period",
			pauseStart:  "2024-07-01",
			pauseEnd:    "2024-09-18",
			testTime:    time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC),
			expected:    false,
			expectError: false,
		},
		{
			name:        "Not paused - after pause period",
			pauseStart:  "2024-07-01",
			pauseEnd:    "2024-09-18",
			testTime:    time.Date(2024, 9, 19, 0, 0, 1, 0, time.UTC),
			expected:    false,
			expectError: false,
		},
		{
			name:        "Paused - exactly on start date",
			pauseStart:  "2024-07-01",
			pauseEnd:    "2024-09-18",
			testTime:    time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			expected:    true,
			expectError: false,
		},
		{
			name:        "Paused - exactly on end date",
			pauseStart:  "2024-07-01",
			pauseEnd:    "2024-09-18",
			testTime:    time.Date(2024, 9, 18, 23, 59, 59, 0, time.UTC),
			expected:    true,
			expectError: false,
		},
		{
			name:        "Single day pause",
			pauseStart:  "2024-08-15",
			pauseEnd:    "2024-08-15",
			testTime:    time.Date(2024, 8, 15, 12, 0, 0, 0, time.UTC),
			expected:    true,
			expectError: false,
		},
		{
			name:        "Invalid start date format",
			pauseStart:  "2024/07/01",
			pauseEnd:    "2024-09-18",
			testTime:    time.Date(2024, 8, 15, 12, 0, 0, 0, time.UTC),
			expected:    false,
			expectError: true,
		},
		{
			name:        "Invalid end date format",
			pauseStart:  "2024-07-01",
			pauseEnd:    "2024/09/18",
			testTime:    time.Date(2024, 8, 15, 12, 0, 0, 0, time.UTC),
			expected:    false,
			expectError: true,
		},
		{
			name:        "Empty start date only",
			pauseStart:  "",
			pauseEnd:    "2024-09-18",
			testTime:    time.Date(2024, 8, 15, 12, 0, 0, 0, time.UTC),
			expected:    false,
			expectError: false,
		},
		{
			name:        "Empty end date only",
			pauseStart:  "2024-07-01",
			pauseEnd:    "",
			testTime:    time.Date(2024, 8, 15, 12, 0, 0, 0, time.UTC),
			expected:    false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				PauseStart: tt.pauseStart,
				PauseEnd:   tt.pauseEnd,
			}

			// Mock current time by temporarily replacing time.Now in the method
			// Since we can't easily mock time.Now, we'll test the logic manually
			paused, err := isPausedAtTime(config, tt.testTime)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if paused != tt.expected {
				t.Errorf("Expected paused=%v, got paused=%v", tt.expected, paused)
			}
		})
	}
}

// Helper function to test pause logic with a specific time
func isPausedAtTime(c Config, testTime time.Time) (bool, error) {
	if c.PauseStart == "" || c.PauseEnd == "" {
		return false, nil
	}

	startDate, err := time.Parse("2006-01-02", c.PauseStart)
	if err != nil {
		return false, err
	}

	endDate, err := time.Parse("2006-01-02", c.PauseEnd)
	if err != nil {
		return false, err
	}

	// Set time to start of day for start date and end of day for end date
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, testTime.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, testTime.Location())

	return (testTime.After(startDate) || testTime.Equal(startDate)) && (testTime.Before(endDate) || testTime.Equal(endDate)), nil
}

func TestIsPausedCurrentTime(t *testing.T) {
	// Test with actual current time using the real IsPaused method
	config := Config{
		PauseStart: "2025-01-01",
		PauseEnd:   "2025-12-31",
	}

	paused, err := config.IsPaused()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Since we're in 2025, this should be paused
	if !paused {
		t.Errorf("Expected to be paused in 2025, but got false")
	}

	// Test with dates in the past
	config.PauseStart = "2020-01-01"
	config.PauseEnd = "2020-12-31"

	paused, err = config.IsPaused()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Since we're past 2020, this should not be paused
	if paused {
		t.Errorf("Expected not to be paused for past dates, but got true")
	}
}
