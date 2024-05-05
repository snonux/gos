package main

import (
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
}

func newEntry(data []byte) (entry, error) {
	var entry entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return entry, fmt.Errorf("unable to deserialise payload: %w", err)
	}
	return entry, nil
}

func (e entry) serialize() ([]byte, error) {
	return json.Marshal(e)
}
