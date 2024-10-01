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
const oldestValidTime = "20240912-102800"

type Entry struct {
	Path  string
	Time  time.Time
	State State
}

func (e Entry) String() string {
	return fmt.Sprintf("Path:%s;Stamp:%s,State:%s",
		e.Path, e.Time.Format(format.Time), e.State)
}

var Zero = Entry{}

// filePath format: /foo/foobarbaz.something.here.txt.STAMP.{posted,queued}
func New(filePath string) (Entry, error) {
	e := Entry{Path: filePath}

	// We want to get the STAMP!
	parts := strings.Split(filePath, ".")
	if len(parts) < 4 {
		return e, fmt.Errorf("not a valid entry path: %s", filePath)
	}

	switch parts[len(parts)-1] {
	case "queued":
		e.State = Queued
	case "posted":
		e.State = Posted
	default:
		return e, fmt.Errorf("can't parse state from path: %s", filePath)
	}

	var err error
	if e.Time, err = time.Parse(format.Time, parts[len(parts)-2]); err != nil {
		return e, err
	}

	oldestValid, err := time.Parse(format.Time, oldestValidTime)
	if err != nil {
		panic(err)
	}

	if e.Time.Before(oldestValid) {
		return e, fmt.Errorf("entry time does not seem legit, it is too old: %v", e.Time)
	}

	return e, nil
}
