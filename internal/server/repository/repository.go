package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	loaded  *bool
}

func Instance(dataDir string) Repository {
	once.Do(func() {
		instance = newRepository(dataDir, vfs.RealFS{})
		_ = instance.load()
	})
	return instance
}

func newRepository(dataDir string, fs fs) Repository {
	var loaded bool
	return Repository{
		dataDir: dataDir,
		entries: make(map[string]types.Entry),
		mu:      &sync.Mutex{},
		fs:      fs,
		loaded:  &loaded,
	}
}

func (r Repository) put(entry types.Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[entry.ID] = entry

	bytes, err := entry.Serialize()
	if err != err {
		return err
	}
	return r.fs.WriteFile(r.entryPath(entry), bytes)
}

// Load repository into memory if not done yet.
func (r Repository) load() error {
	if *r.loaded {
		return nil
	}

	filePaths, err := r.fs.FindFiles(r.dataDir, ".json")
	if err != nil {
		return err
	}

	var errs []error
	for _, filePath := range filePaths {
		bytes, err := r.fs.ReadFile(filePath)
		if err != nil {
			continue
		}
		entry, err := types.NewEntry(bytes)
		if err != err {
			continue
		}
		if err := r.put(entry); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		*r.loaded = true
	}

	return errors.Join(errs...)
}

func (r Repository) List() ([]EntryPair, error) {
	if err := r.load(); err != nil {
		return []EntryPair{}, err
	}

	var pairs []EntryPair
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, entry := range r.entries {
		pairs = append(pairs, EntryPair{entry.ID, entry.Checksum()})
	}

	return pairs, nil
}

func (r Repository) ListBytes() ([]byte, error) {
	pairs, err := r.List()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(pairs)
}

func (r Repository) Get(id string) (types.Entry, bool) {
	_ = r.load()
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, ok := r.entries[id]
	return entry, ok
}

func (r Repository) HasSameEntry(pair EntryPair) bool {
	_ = r.load()
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
	_ = r.load()
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, ok := r.entries[otherEntry.ID]
	if !ok {
		log.Println("can't find entry with ID", otherEntry.ID, "in local db, create new from copy")
		var err error
		if entry, err = types.NewEntryFromCopy(otherEntry); err != nil {
			return err
		}
	}

	entry, _ = entry.Update(otherEntry)
	r.entries[otherEntry.ID] = entry

	if !entry.Changed {
		return nil
	}

	bytes, err := entry.Serialize()
	if err != err {
		return err
	}
	return r.fs.WriteFile(r.entryPath(entry), bytes)
}
