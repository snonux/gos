package entry

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/oi"
	"codeberg.org/snonux/gos/internal/prompt"
	"codeberg.org/snonux/gos/internal/timestamp"
)

type State int

const (
	Unknown State = iota
	Inboxed
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
	Tags  map[string]struct{}
}

func (en Entry) String() string {
	if en.State == Inboxed {
		return fmt.Sprintf("Path:%s;State:%s", en.Path, en.State)
	}
	return fmt.Sprintf("Path:%s;Stamp:%s,State:%s", en.Path, en.Time.Format(timestamp.Format), en.State)
}

// filePath format: /foo/foobarbaz.something.here.txt.STAMP.{posted,queued}
// or for inboxed: /foo.txt
// or inboxed with tags: /foo.prio.ask.txt
func New(filePath string) (Entry, error) {
	en := Entry{Path: filePath, Tags: make(map[string]struct{})}

	// We want to get the STAMP!
	parts := strings.Split(filePath, ".")
	if len(parts) < 2 {
		// Could be 2 if inboxed
		return en, fmt.Errorf("not a valid entry path: %s", filePath)
	}
	en.extractTags(parts)

	switch parts[len(parts)-1] {
	case "queued":
		en.State = Queued
	case "posted":
		en.State = Posted
	default:
		en.State = Inboxed
		return en, nil
	}

	if len(parts) < 4 {
		// If not inboxed, must be longer.
		return en, fmt.Errorf("not a valid entry path: %s", filePath)
	}

	var err error
	if en.Time, err = timestamp.Parse(parts[len(parts)-2]); err != nil {
		return en, err
	}

	if en.Time.Before(timestamp.OldestValidTime()) {
		return en, fmt.Errorf("entry time does not seem legit, it is too old: %v", en.Time)
	}

	return en, nil
}

func (en *Entry) Content() (string, []string, error) {
	content, err := oi.SlurpAndTrim(en.Path)
	return content, extractURLs(content), err
}

// Returns the Name, e.g. foo.bar.baz from /path/foo.bar.baz.TIMESTAMP.posted
func (en *Entry) Name() string {
	base := filepath.Base(en.Path)
	parts := strings.Split(base, ".")

	offset := len(parts) - 1

	switch en.State {
	case Queued:
		fallthrough
	case Posted:
		offset -= 2
	}

	// TODO: Unit test this
	return strings.Join(parts[:offset], ".")
}

// Returns the content and also checks for the size limit
func (en Entry) ContentWithLimit(sizeLimit int) (string, []string, error) {
	content, urls, err := en.Content()
	if err != nil {
		return "", urls, err
	}
	if len(content) > sizeLimit {
		err := fmt.Errorf("%w (%d > %d): %v", ErrSizeLimitExceeded, len(content), sizeLimit, en)
		if err2 := prompt.Acknowledge("You need to shorten the content as "+err.Error(), content); err2 != nil {
			return "", urls, errors.Join(err, err2)
		}
		if err2 := en.Edit(); err2 != nil {
			return "", urls, errors.Join(err, err2)
		}
		return en.ContentWithLimit(sizeLimit)
	}
	return content, urls, nil
}

func (en *Entry) MarkPosted() error {
	if en.State == Inboxed {
		return errors.New("entry still inboxed, can not mark as posted")
	}
	if en.State != Queued {
		return errors.New("entry is not queued")
	}
	if en.State == Posted {
		return errors.New("entry is already posted")
	}
	newPath, err := timestamp.UpdateInFilename(strings.TrimSuffix(en.Path, ".queued")+".posted", -2)
	if err != nil {
		return err
	}
	if err := os.Rename(en.Path, newPath); err != nil {
		return err
	}
	en.State = Posted
	return nil
}

func (en Entry) HasTag(tag string) bool {
	_, ok := en.Tags[tag]
	return ok
}

func (en Entry) Edit() error {
	if err := prompt.EditFile(en.Path); err != nil {
		return err
	}
	return nil
}

func (en Entry) Remove() error {
	return os.Remove(en.Path)
}

func (en Entry) FileAction(question string) error {
	content, _, err := en.Content()
	if err != nil {
		return err
	}
	_, err = prompt.FileAction(question, content, en.Path)
	return err
}

func (en Entry) extractTags(parts []string) {
	for _, part := range parts {
		if slices.Contains(validTags, part) || strings.HasPrefix(part, "share:") {
			en.Tags[part] = struct{}{}
		}
	}
}

func extractURLs(input string) []string {
	urlPattern := `(http://|https://|ftp://)[^\s]+`
	re := regexp.MustCompile(urlPattern)
	return re.FindAllString(input, -1)
}
