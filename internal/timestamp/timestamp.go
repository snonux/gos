package timestamp

import (
	"strings"
	"time"
)

const (
	Format = "20060102-150405"
)

func Now() string {
	return time.Now().Format(Format)
}

func NowTime() time.Time {
	t, err := Parse(Now())
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(timeStr string) (time.Time, error) {
	return time.Parse(Format, timeStr)
}
func OldestValidTime() time.Time {
	// The time this code was written a:round, actually.
	oldestValidTime, err := Parse("20240912-102800")
	if err != nil {
		// Never expected
		panic(err)
	}
	return oldestValidTime
}

func UpdateInFilename(filename string, rIndex int) string {
	parts := strings.Split(filename, ".")
	parts[len(parts)+rIndex] = Now()
	return strings.Join(parts, ".")
}
