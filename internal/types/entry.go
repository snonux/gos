package types

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Tells me whether the entry was shared to the sm platform named Name
type Shared struct {
	Name string `json:"name"`
	Is   bool   `json:"is,omitempty"`
}

func (s Shared) String() string {
	return fmt.Sprintf("Name:%s;Is:%v", s.Name, s.Is)
}

func (s Shared) Equals(other Shared) bool {
	return s.Name == other.Name && s.Is == other.Is
}

type Entry struct {
	// The unique ID of this entry.
	ID     string   `json:"id,omitempty"`
	Body   string   `json:"body"`
	Shared []Shared `json:"shared,omitempty"`
	Epoch  int      `json:"epoch,omitempty"`

	// The checksum of the whole entry, can change depending on the state.
	checksum      string
	checksumDirty bool
	mu            *sync.Mutex

	// To identify whether this entry was changed.
	Changed bool `json:"-"`
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

func NewEntryFromCopy(other Entry) (Entry, error) {
	e := other
	e.initialize()

	e.Shared = make([]Shared, len(other.Shared))
	copy(e.Shared, other.Shared)

	return e, nil
}

func NewEntryFromTextFile(filePath string) (Entry, error) {
	var e Entry

	data, err := os.ReadFile(filePath)
	if err != nil {
		return e, err
	}

	e.Body = string(data)
	if e.ID == "" {
		e.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(e.Body)))
	}

	e.initialize()
	e.Checksum()

	return e, nil
}

func (e *Entry) initialize() {
	e.mu = &sync.Mutex{}
	e.checksumDirty = true
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
		return e, fmt.Errorf("can update entry only with other entry with same ID: this(%s) other(%s)", e, other)
	}
	e.checksumDirty = true
	e.Changed = true

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

func (e Entry) JSONMarshal() ([]byte, error) {
	return json.Marshal(e)
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
