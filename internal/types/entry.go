package types

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"codeberg.org/snonux/gos/internal/vfs"
)

type fs interface {
	ReadFile(name string) ([]byte, error)
	WriteFile(filePath string, bytes []byte) error
}

// Tells me whether the entry was shared to the sm platform named Name
type Shared struct {
	Name string `json:"name"`
	Is   bool   `json:"is,omitempty"`
}

func (s Shared) String() string {
	return fmt.Sprintf("Name:%s;Is:%v", s.Name, s.Is)
}

func (s Shared) Equals(other Shared) bool {
	switch {
	case s.Name != other.Name:
		return false
	case s.Is != other.Is:
		return false
	default:
		return true
	}
}

type Entry struct {
	// The unique ID of this entry.
	ID     string   `json:"id,omitempty"`
	Body   string   `json:"body"`
	Shared []Shared `json:"shared,omitempty"`
	Epoch  int      `json:"epoch,omitempty"`
	fs     fs

	// The checksum of the whole entry, can change depending on the state.
	checksum      string
	checksumDirty bool
	mu            *sync.Mutex
}

func NewEntry(bytes []byte) (Entry, error) {
	var e Entry
	if err := json.Unmarshal(bytes, &e); err != nil {
		return e, fmt.Errorf("unable to deserialise payload: %w", err)
	}
	e.initialize()
	if e.ID == "" {
		e.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(e.Body)))
	}
	return e, nil
}

func NewEntryFromFile(filePath string, fsToUse ...fs) (Entry, error) {
	var (
		bytes []byte
		err   error
		fs    fs = vfs.RealFS{}
	)

	if len(fsToUse) > 0 {
		fs = fsToUse[0]
	}

	bytes, err = fs.ReadFile(filePath)
	if err != err {
		return Entry{}, err
	}
	e, err := NewEntry(bytes)
	e.fs = fs
	return e, err
}

func NewEntryFromCopy(other Entry) (Entry, error) {
	var e Entry
	e.initialize()
	return e.Update(other)
}

func (e *Entry) initialize() {
	e.mu = &sync.Mutex{}
	e.checksumDirty = true
	e.fs = vfs.RealFS{}
}

func (e Entry) Equals(other Entry) bool {
	switch {
	case e.Body != other.Body:
		return false
	case e.Epoch != other.Epoch:
		return false
	case e.ID != other.ID:
		return false
	case len(e.Shared) != len(other.Shared):
		return false
	}

	otherShared := make(map[string]Shared)
	for _, shared := range other.Shared {
		otherShared[shared.Name] = shared
	}

	for _, shared := range e.Shared {
		otherShared, ok := otherShared[shared.Name]
		if !ok || !shared.Equals(otherShared) {
			return false
		}
	}

	return true
}

/**
 * This updates the entry with the other entry. The Shared slice will also be
 * updated. If entry is missing, it will be added. If entry is there, the shared
 * Is status will eventually flip to true but never to false.
 */
func (e Entry) Update(other Entry) (Entry, error) {
	if e.ID != other.ID {
		return e, fmt.Errorf("can update entry only with other entry with same ID: %s %s", e, other)
	}
	e.checksumDirty = true

	if e.Body != other.Body {
		e.Body = other.Body
	}

	if e.Epoch != other.Epoch {
		e.Epoch = other.Epoch
	}

	sharedMap := make(map[string]Shared)
	for _, shared := range e.Shared {
		sharedMap[shared.Name] = shared
	}

	for _, otherShared := range other.Shared {
		shared, ok := sharedMap[otherShared.Name]
		switch {
		case !ok:
			sharedMap[otherShared.Name] = shared
			continue
		case otherShared.Is:
			shared.Is = true
			sharedMap[otherShared.Name] = shared
		}
	}

	e.Shared = e.Shared[:0]
	for _, shared := range sharedMap {
		e.Shared = append(e.Shared, shared)
	}

	return e, nil
}

func (e Entry) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e Entry) SaveFile(filePath string) error {
	jsonStr, err := e.Serialize()
	if err != nil {
		return err
	}

	return e.fs.WriteFile(filePath, jsonStr)
}

func (e Entry) String() string {
	var sb strings.Builder

	sb.WriteString("ID:")
	sb.WriteString(e.ID)
	sb.WriteString(";")
	sb.WriteString(fmt.Sprintf("Epoch:%d;", e.Epoch))
	sb.WriteString("Shared:[")
	for i, shared := range e.Shared {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(shared.String())
	}
	sb.WriteString("];")
	sb.WriteString("Body:")
	sb.WriteString(e.Body)

	return sb.String()
}

func (e *Entry) Checksum() string {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.checksumDirty {
		return e.checksum
	}

	e.checksum = fmt.Sprintf("%x", sha256.Sum256([]byte(e.String())))
	e.checksumDirty = false
	return e.checksum
}
