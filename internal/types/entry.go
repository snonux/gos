package types

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Shared struct {
	Name string `json:"id"`
	Is   bool   `json:"is,omitempty"`
}

type Entry struct {
	Body   string   `json:"body"`
	Shared []Shared `json:"shared,omitempty"`
	Epoch  int      `json:"epoch,omitempty"`
	ID     string   `json:"id,omitempty"`
}

func NewEntry(bytes []byte) (Entry, error) {
	var entry Entry
	if err := json.Unmarshal(bytes, &entry); err != nil {
		return entry, fmt.Errorf("unable to deserialise payload: %w", err)
	}
	if entry.ID == "" {
		entry.ID = fmt.Sprintf("%x", sha256.Sum256(bytes))
	}
	return entry, nil
}

func (e Entry) Serialize() ([]byte, error) {
	return json.Marshal(e)
}
