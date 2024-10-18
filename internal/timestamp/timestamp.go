package timestamp

import (
	"fmt"
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
		panic(err) // Never expected
	}
	return oldestValidTime
}

func UpdateInFilename(filename string, rIndex int) (string, error) {
	parts := strings.Split(filename, ".")
	ind := len(parts) + rIndex
	if ind < 0 || ind >= len(parts) {
		return "", fmt.Errorf("unable to update timestamp in %s, invalid index %d: %v",
			filename, ind, parts)
	}
	parts[ind] = Now()
	return strings.Join(parts, "."), nil
}
