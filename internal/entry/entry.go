package entry

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/prompt"
	"codeberg.org/snonux/gos/internal/timestamp"
)

type State int

const (
	Unknown State = iota
	Inboxed       // TODO: Implement
	Queued
	Posted
)

var (
	validTags            = []string{"ask", "prio", "now"}
	ErrSizeLimitExceeded = errors.New("message size limit exceeded")
)

func (s State) String() string {
	switch s {
	case Unknown:
		return "unknown"
	case Inboxed:
		return "inboxed"
	case Queued:
		return "queued"
	case Posted:
		return "posted"
	default:
		panic(fmt.Sprintf("unknown state: %d", int(s)))
	}
}

var Zero = Entry{}

type Entry struct {
	Path  string
	Time  time.Time
	State State
	tags  []string
}

func (e Entry) String() string {
	if e.State == Inboxed {
		return fmt.Sprintf("Path:%s;State:%s", e.Path, e.State)
	}
	return fmt.Sprintf("Path:%s;Stamp:%s,State:%s", e.Path, e.Time.Format(timestamp.Format), e.State)
}

// filePath format: /foo/foobarbaz.something.here.txt.STAMP.{posted,queued}
// or for inboxed: /foo.txt
// or inboxed with tags: /foo.prio.ask.txt
func New(filePath string) (Entry, error) {
	e := Entry{Path: filePath}

	// We want to get the STAMP!
	parts := strings.Split(filePath, ".")
	if len(parts) < 2 {
		// Could be 2 if inboxed
		return e, fmt.Errorf("not a valid entry path: %s", filePath)
	}

	for _, part := range parts {
		if slices.Contains(validTags, part) {
			e.tags = append(e.tags, part)
		}
	}

	switch parts[len(parts)-1] {
	case "queued":
		e.State = Queued
	case "posted":
		e.State = Posted
	default:
		e.State = Inboxed
		return e, nil
	}

	if len(parts) < 4 {
		// If not inboxed, must be longer.
		return e, fmt.Errorf("not a valid entry path: %s", filePath)
	}
	var err error
	if e.Time, err = timestamp.Parse(parts[len(parts)-2]); err != nil {
		return e, err
	}

	if e.Time.Before(timestamp.OldestValidTime()) {
		return e, fmt.Errorf("entry time does not seem legit, it is too old: %v", e.Time)
	}

	return e, nil
}

func (e *Entry) Content() (string, []string, error) {
	bytes, err := os.ReadFile(e.Path)
	if err != nil {
		return "", []string{}, err
	}
	content := strings.TrimSpace(string(bytes))
	return content, extractURLs(content), nil
}

func (e Entry) ContentWithLimit(sizeLimit int) (string, []string, error) {
	content, urls, err := e.Content()
	if err != nil {
		return "", urls, err
	}
	if len(content) > sizeLimit {
		err := fmt.Errorf("%w (%d > %d): %v", ErrSizeLimitExceeded, len(content), sizeLimit, e)
		if err2 := prompt.Acknowledge("You need to shorten the content as "+err.Error(), content); err2 != nil {
			return "", urls, errors.Join(err, err2)
		}
		if err2 := e.Edit(); err2 != nil {
			return "", urls, errors.Join(err, err2)
		}
		return e.ContentWithLimit(sizeLimit)
	}
	return content, urls, nil
}

func (e *Entry) MarkPosted() error {
	if e.State == Inboxed {
		return errors.New("entry still inboxed, can not mark as posted")
	}
	if e.State != Queued {
		return errors.New("entry is not queued")
	}
	if e.State == Posted {
		return errors.New("entry is already posted")
	}
	newPath, err := timestamp.UpdateInFilename(strings.TrimSuffix(e.Path, ".queued")+".posted", -2)
	if err != nil {
		return err
	}
	if err := os.Rename(e.Path, newPath); err != nil {
		return err
	}
	e.State = Posted
	return nil
}

func (e Entry) HasTag(tag string) bool {
	return slices.Contains(e.tags, tag)
}

func (e Entry) Edit() error {
	if err := prompt.EditFile(e.Path); err != nil {
		return err
	}
	return nil
}

func extractURLs(input string) []string {
	urlPattern := `(http://|https://|ftp://)[^\s]+`
	re := regexp.MustCompile(urlPattern)
	return re.FindAllString(input, -1)
}
