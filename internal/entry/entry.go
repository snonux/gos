package entry

import (
	"fmt"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/format"
)

type State int

const (
	Unknown State = iota
	Queued
	Posted
)

func (s State) String() string {
	switch s {
	case Unknown:
		return "unknown"
	case Queued:
		return "queued"
	case Posted:
		return "posted"
	default:
		panic(fmt.Sprintf("unknown state: %d", int(s)))
	}
}

// The time this code was written a:round, actually.
const oldestValidTime = "20240922-102800"

type Entry struct {
	Path  string
	Time  time.Time
	State State
}

// filePath format: /foo/foobarbaz.something.here.txt.STAMP.{posted,queued}
func New(filePath string) (Entry, error) {
	ent := Entry{Path: filePath}

	// We want to get the STAMP!
	parts := strings.Split(filePath, ".")
	if len(parts) < 4 {
		return ent, fmt.Errorf("not a valid entry path: %s", filePath)
	}

	switch parts[len(parts)-1] {
	case "queued":
		ent.State = Queued
	case "posted":
		ent.State = Posted
	default:
		return ent, fmt.Errorf("can't parse state from path: %s", filePath)
	}

	var err error
	if ent.Time, err = time.Parse(format.Time, parts[len(parts)-2]); err != nil {
		return ent, err
	}

	oldestValid, err := time.Parse(format.Time, oldestValidTime)
	if err != nil {
		panic(err)
	}

	if ent.Time.Before(oldestValid) {
		return ent, fmt.Errorf("entry time does not seem legic, it is too old: %v", ent.Time)
	}

	return ent, nil
}
