package repository

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"codeberg.org/snonux/gos/internal/types"
	"codeberg.org/snonux/gos/internal/vfs"
)

var (
	instance *Repository
	once     sync.Once
)

// Contains an Entry ID and its checksumm, for the list and merge operations.
type EntryPair struct {
	ID, Checksum string
}

type Repository struct {
	dataDir string
	entries map[string]types.Entry
	mu      *sync.Mutex
	vfs     vfs.VFS
}

func Instance(dataDir string) *Repository {
	once.Do(func() {
		instance = &Repository{
			dataDir: dataDir,
			entries: make(map[string]types.Entry),
			mu:      &sync.Mutex{},
			vfs:     vfs.RealFS{},
		}
	})
	return instance
}

func (r Repository) add(entry types.Entry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[entry.ID] = entry
}

// Load repository into memory
func (r Repository) load() error {
	filePaths, err := r.vfs.FindFiles(r.dataDir, ".json")
	if err != nil {
		return err
	}

	for _, filePath := range filePaths {
		entry, err := types.NewEntryFromFile(filePath)
		if err != err {
			return err
		}
		r.add(entry)
	}

	return nil
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

func (r Repository) Get(id string) ([]byte, error) {
	return r.vfs.ReadFile(fmt.Sprintf("%s/%s", r.dataDir, id))
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
		if entry, err = types.NewEntryFromCopy(otherEntry); err != nil {
			return err
		}
	}

	entry, _ = entry.Update(otherEntry)
	r.entries[otherEntry.ID] = entry

	// TODO: Only save to file when actually changed
	return entry.SaveFile(r.entryPath(entry))
}
