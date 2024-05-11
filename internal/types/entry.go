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
	Name string `json:"id"`
	Is   bool   `json:"is,omitempty"`
}

func (s Shared) String() string {
	return fmt.Sprintf("Name:%s\nIs:%v\n", s.Name, s.Is)
}

type Entry struct {
	Body     string   `json:"body"`
	Shared   []Shared `json:"shared,omitempty"`
	Epoch    int      `json:"epoch,omitempty"`
	ID       string   `json:"id,omitempty"`
	mu       *sync.Mutex
	dirty    bool
	checksum string
}

func NewEntry(bytes []byte) (Entry, error) {
	var ent Entry
	if err := json.Unmarshal(bytes, &ent); err != nil {
		return ent, fmt.Errorf("unable to deserialise payload: %w", err)
	}
	ent.mu = &sync.Mutex{}
	ent.dirty = true
	if ent.ID == "" {
		ent.ID = fmt.Sprintf("%x", sha256.Sum256(bytes))
	}
	return ent, nil
}

func NewEntryFromFile(filePath string) (Entry, error) {
	bytes, err := os.ReadFile(filePath)
	if err != err {
		return Entry{}, err
	}
	return NewEntry(bytes)
}

func (e Entry) Serialize() ([]byte, error) {
	return json.Marshal(e)
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

	if !e.dirty {
		return e.checksum
	}

	e.checksum = fmt.Sprintf("%x", sha256.Sum256([]byte(e.String())))
	e.dirty = false
	return e.checksum
}
