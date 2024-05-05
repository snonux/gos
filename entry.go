package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type shared struct {
	Name string `json:"id"`
	Is   bool   `json:"is,omitempty"`
}

type entry struct {
	Body   string   `json:"body"`
	Shared []shared `json:"shared,omitempty"`
	Epoch  int      `json:"epoch,omitempty"`
	id     string
}

func newEntry(bytes []byte) (entry, error) {
	var entry entry
	if err := json.Unmarshal(bytes, &entry); err != nil {
		return entry, fmt.Errorf("unable to deserialise payload: %w", err)
	}
	entry.id = fmt.Sprintf("%x", sha256.Sum256(bytes))
	return entry, nil
}

func (e entry) serialize() ([]byte, error) {
	return json.Marshal(e)
}
