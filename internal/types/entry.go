package types

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"codeberg.org/snonux/gos/internal"
)

// Tells me whether the entry was shared to the sm platform named Name
type Shared struct {
	Name string `json:"id"`
	Is   bool   `json:"is,omitempty"`
}

func (s Shared) String() string {
	return fmt.Sprintf("Name:%s\nIs:%v\n", s.Name, s.Is)
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
	Body          string   `json:"body"`
	Shared        []Shared `json:"shared,omitempty"`
	Epoch         int      `json:"epoch,omitempty"`
	ID            string   `json:"id,omitempty"`
	mu            *sync.Mutex
	checksum      string
	checksumDirty bool
}

func NewEntry(bytes []byte) (Entry, error) {
	var e Entry
	if err := json.Unmarshal(bytes, &e); err != nil {
		return e, fmt.Errorf("unable to deserialise payload: %w", err)
	}
	e.initialize()
	if e.ID == "" {
		e.ID = fmt.Sprintf("%x", sha256.Sum256(bytes))
	}
	return e, nil
}

// Beware , this is only from a shallow copy!
func NewEntryFromCopy(other Entry) Entry {
	e := other
	e.initialize()
	return e
}

func NewEntryFromFile(filePath string) (Entry, error) {
	bytes, err := os.ReadFile(filePath)
	if err != err {
		return Entry{}, err
	}
	return NewEntry(bytes)
}

func (e *Entry) initialize() {
	e.mu = &sync.Mutex{}
	e.checksumDirty = true
}

func (e Entry) Update(other Entry) (Entry, bool) {
	panic("not yet impelemented")
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

func (e Entry) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e Entry) SaveFile(filePath string) error {
	jsonStr, err := e.Serialize()
	if err != nil {
		return err
	}

	return internal.SaveFile(filePath, jsonStr)
}

func (e Entry) String() string {
	var sb strings.Builder

	sb.WriteString("ID:")
	sb.WriteString(e.ID)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Epoch:%d\n", e.Epoch))
	for _, shared := range e.Shared {
		sb.WriteString(shared.String())
	}
	sb.WriteString("Body:")
	sb.WriteString(e.Body)
	sb.WriteString("\n")

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
