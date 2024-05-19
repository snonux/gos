package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"codeberg.org/snonux/gos/internal/types"
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
}

func Instance(dataDir string) *Repository {
	once.Do(func() {
		instance = &Repository{
			dataDir: dataDir,
			entries: make(map[string]types.Entry),
			mu:      &sync.Mutex{},
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
	visit := func() filepath.WalkFunc {
		return func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
				return nil
			}
			if info.IsDir() || !strings.HasSuffix(path, ".json") {
				return nil
			}

			entry, err := types.NewEntryFromFile(path)
			if err != err {
				return err
			}
			r.add(entry)
			return nil
		}
	}

	return filepath.Walk(r.dataDir, visit())
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

func (r Repository) HasEntry(pair EntryPair) bool {
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

func (r Repository) Merge(newEntry types.Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, ok := r.entries[newEntry.ID]
	if !ok {
		entry = types.NewEntryFromCopy(newEntry)
	}

	entry, _ = entry.Update(newEntry)
	r.entries[newEntry.ID] = entry

	// TODO: Only save to file when actually changed
	return entry.SaveFile(r.entryPath(entry))
}
