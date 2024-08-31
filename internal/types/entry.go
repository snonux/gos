package types

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"sort"
	"strings"
)

type Entry struct {
	// The unique ID of this entry.
	ID     string            `json:"id,omitempty"`
	Body   string            `json:"body"`
	Shared map[string]Shared `json:"shared,omitempty"`
	Epoch  int               `json:"epoch,omitempty"`

	// The checksum of the whole entry, can change depending on the state.
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
		e.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(e.Body)))
	}

	return e, nil
}

func NewEntryFromCopy(other Entry) (Entry, error) {
	e := other
	e.initialize()
	e.Shared = maps.Clone(other.Shared)

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
	if e.Shared == nil {
		e.Shared = make(map[string]Shared)
	}
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
		// case len(e.Shared) != len(other.Shared):
		// 	return false
	}

	return maps.Equal(e.Shared, other.Shared)
}

func (e Entry) IsShared(name string) bool {
	shared, ok := e.Shared[name]
	if !ok {
		return false
	}
	return shared.Is
}

/**
 * This updates the entry with the other entry. The Shared slice will also be
 * updated. If entry is missing, it will be added. If entry is there, the shared
 * Is status will eventually flip to true but never to false.
 */
func (e Entry) Update(other Entry) (Entry, bool, error) {
	if e.ID != other.ID {
		return e, false, fmt.Errorf("can update entry only with other entry with same ID: this(%s) other(%s)", e, other)
	}

	var changed bool

	if e.Body != other.Body {
		e.Body = other.Body
		changed = true
	}

	if e.Epoch != other.Epoch {
		e.Epoch = other.Epoch
		changed = true
	}

	for otherName, otherShared := range other.Shared {
		shared, ok := e.Shared[otherName]
		switch {
		case !ok:
			e.Shared[otherName] = shared
			changed = true
		case otherShared.Is && !shared.Is:
			shared.Is = true
			e.Shared[otherName] = shared
			changed = true
		}
	}

	if changed {
		e.checksumDirty = true
	}

	return e, changed, nil
}

func (e Entry) JSONMarshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e Entry) String() string {
	return e.checksumBase()
}

/**
 * Used to calculate the checksum, better don't change the output, otherwise
 * repository database will get confused with entry checksum mismatches.
 */
func (e Entry) checksumBase() string {
	var sb strings.Builder

	sb.WriteString("ID:")
	sb.WriteString(e.ID)
	sb.WriteString(";")
	sb.WriteString(fmt.Sprintf("Epoch:%d;", e.Epoch))
	sb.WriteString("Shared:{")

	keys := make([]string, 0, len(e.Shared))
	for key := range e.Shared {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for i, sharedName := range keys {
		if i > 0 {
			sb.WriteString(",")
		}
		shared := e.Shared[sharedName]
		sb.WriteString(sharedName)
		sb.WriteString(":{")
		sb.WriteString(shared.String())
		sb.WriteString("}")
	}

	sb.WriteString("};")
	sb.WriteString("Body:")
	sb.WriteString(e.Body)

	return sb.String()
}

func (e *Entry) Checksum() string {
	if !e.checksumDirty {
		return e.checksum
	}

	e.checksum = fmt.Sprintf("%x", sha256.Sum256([]byte(e.checksumBase())))
	e.checksumDirty = false
	return e.checksum
}
