package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"codeberg.org/snonux/gos/internal/types"
	"codeberg.org/snonux/gos/internal/vfs"
)

var (
	instance Repository
	once     sync.Once
)

type fs interface {
	ReadFile(name string) ([]byte, error)
	WriteFile(filePath string, bytes []byte) error
	FindFiles(dataPath, suffix string) ([]string, error)
}

// Contains an Entry ID and its checksumm, for the list and merge operations.
type EntryPair struct {
	ID, Checksum string
}

type Repository struct {
	dataDir string
	entries map[string]types.Entry
	mu      *sync.Mutex
	fs      fs
}

func Instance(dataDir string) Repository {
	once.Do(func() {
		instance = newRepository(dataDir, vfs.RealFS{})
	})
	return instance
}

func newRepository(dataDir string, fs fs) Repository {
	return Repository{
		dataDir: dataDir,
		entries: make(map[string]types.Entry),
		mu:      &sync.Mutex{},
		fs:      fs,
	}
}

func (r Repository) put(entry types.Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[entry.ID] = entry
	return entry.SaveFile(r.entryPath(entry))
}

// Load repository into memory
func (r Repository) load() error {
	filePaths, err := r.fs.FindFiles(r.dataDir, ".json")
	if err != nil {
		return err
	}

	var errs []error
	for _, filePath := range filePaths {
		entry, err := types.NewEntryFromFile(filePath, r.fs)
		if err != err {
			errs = append(errs, err)
			continue
		}
		if err := r.put(entry); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (r Repository) List() ([]byte, error) {
	if err := r.load(); err != nil {
		return []byte{}, err
	}

	var pairs []EntryPair
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, entry := range r.entries {
		pairs = append(pairs, EntryPair{entry.ID, entry.Checksum()})
	}

	return json.Marshal(pairs)
}

func (r Repository) Get(id string) (types.Entry, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, ok := r.entries[id]
	return entry, ok
}

func (r Repository) HasSameEntry(pair EntryPair) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	ent, ok := r.entries[pair.ID]
	if !ok || ent.Checksum() != pair.Checksum {
		return false
	}
	return true
}

func (r Repository) entryPath(entry types.Entry) string {
	return fmt.Sprintf("%s/%s/%s.json", r.dataDir, time.Now().Format("2006"), entry.ID)
}

func (r Repository) Merge(otherEntry types.Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, ok := r.entries[otherEntry.ID]
	if !ok {
		var err error
		if entry, err = types.NewEntryFromCopy(otherEntry, r.fs); err != nil {
			return err
		}
	}

	entry, _ = entry.Update(otherEntry)
	r.entries[otherEntry.ID] = entry

	// TODO: Only save to file when actually changed
	return entry.SaveFile(r.entryPath(entry))
}
