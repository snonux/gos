package repository

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"codeberg.org/snonux/gos/internal/types"
)

type Repository struct {
	dataDir string
	entries map[string]types.Entry
	mu      *sync.Mutex
}

func New(dataDir string) Repository {
	return Repository{
		dataDir: dataDir,
		entries: make(map[string]types.Entry),
		mu:      &sync.Mutex{},
	}
}

func (r Repository) store(entry types.Entry) {
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
			r.store(entry)
			return nil
		}
	}

	return filepath.Walk(r.dataDir, visit())
}

func (r Repository) List() ([]byte, error) {
	if err := r.load(); err != nil {
		return []byte{}, err
	}

	type pair struct {
		ID, Checksum string
	}

	var pairs []pair
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, entry := range r.entries {
		pairs = append(pairs, pair{entry.ID, entry.Checksum()})
	}

	return json.Marshal(pairs)
}
