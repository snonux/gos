package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"codeberg.org/snonux/gos/internal/config/server"
	"codeberg.org/snonux/gos/internal/easyhttp"
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
type entryPair struct {
	ID, Checksum string
}

// Holds all entries in the database / stores them to the disks..
type Repository struct {
	pending
	conf    server.ServerConfig
	entries map[types.EntryID]types.Entry
	mu      *sync.Mutex
	fs      fs
	loaded  *bool
	getIdRe *regexp.Regexp
}

func Instance(conf server.ServerConfig) Repository {
	once.Do(func() {
		instance = newRepository(conf, vfs.RealFS{})
	})
	return instance
}

// Need to register all social platforms for in-memory representation of shared posts and so on.
func newRepository(conf server.ServerConfig, fs fs) Repository {
	var loaded bool
	return Repository{
		pending: newPending(),
		conf:    conf,
		entries: make(map[types.EntryID]types.Entry),
		mu:      &sync.Mutex{},
		fs:      fs,
		loaded:  &loaded,
		getIdRe: regexp.MustCompile(`^[a-z0-9]{64}$`),
	}
}

// Gets next entry to be shared for the given social platform.
func (r Repository) Next(platform types.PlatformName) (types.Entry, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, entry := range r.entries {
		if !entry.IsShared(platform) {
			return entry, true
		}
	}

	return types.Entry{}, false // No entry found
}

// Load repository into memory if not done yet.
func (r Repository) load() error {
	if *r.loaded {
		return nil
	}

	filePaths, err := r.fs.FindFiles(r.conf.DataDir, ".json")
	if err != nil {
		return err
	}

	var errs []error
	for _, filePath := range filePaths {
		log.Println("loading entry", filePath)

		bytes, err := r.fs.ReadFile(filePath)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		entry, err := types.NewEntry(bytes)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		r.mu.Lock()
		r.entries[entry.ID] = entry
		r.mu.Unlock()
	}

	if len(errs) == 0 {
		*r.loaded = true
	}

	return errors.Join(errs...)
}

func (r Repository) List() ([]entryPair, error) {
	if err := r.load(); err != nil {
		return []entryPair{}, err
	}

	var pairs []entryPair
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, entry := range r.entries {
		pairs = append(pairs, entryPair{entry.ID, entry.Checksum()})
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

// put writes exact the same entry to the repository. Whereas merge
// Is a bit more refined, tries to merge the same entry wich are slightly
// different into the same entry.
func (r Repository) put(entry types.Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[entry.ID] = entry

	bytes, err := entry.JSONMarshal()
	if err != err {
		return err
	}
	return r.fs.WriteFile(r.entryPath(entry), bytes)
}

func (r Repository) Get(id types.EntryID) (types.Entry, error) {
	if !r.getIdRe.MatchString(id) {
		return types.Entry{}, fmt.Errorf("invalid id %s", id)
	}
	if err := r.load(); err != nil {
		return types.Entry{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	entry, ok := r.entries[id]
	if !ok {
		return entry, fmt.Errorf("no entry with id %s found", id)
	}
	return entry, nil
}

func (r Repository) GetJSON(id types.EntryID) (string, error) {
	entry, err := r.Get(id)
	if err != nil {
		return "", err
	}

	bytes, err := entry.JSONMarshal()
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func (r Repository) hasSameEntry(pair entryPair) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, ok := r.entries[pair.ID]
	if !ok || entry.Checksum() != pair.Checksum {
		return false
	}
	return true
}

func (r Repository) entryPath(ent types.Entry) string {
	return fmt.Sprintf("%s/%s/%s.json", r.conf.DataDir, time.Now().Format("2006"), ent.ID)
}

func (r Repository) Update(ent types.Entry) error {
	// Update is just an alias for the merge, makes the intention clearer.
	return r.Merge(ent)
}

func (r Repository) Merge(otherEnt types.Entry) error {
	if err := r.load(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	entry, ok := r.entries[otherEnt.ID]
	if !ok {
		log.Println("can't find entry with ID", otherEnt.ID, "in local db, create new from copy")
		var err error
		if entry, err = types.NewEntryFromCopy(otherEnt); err != nil {
			return err
		}
	}

	var changed bool
	entry, changed, _ = entry.Update(otherEnt)
	r.entries[otherEnt.ID] = entry

	if !changed {
		// Hasn't changed, so no need to write anything to file.
		return nil
	}

	bytes, err := entry.JSONMarshal()
	if err != err {
		return err
	}

	return r.fs.WriteFile(r.entryPath(entry), bytes)
}

func (r Repository) MergeRemotely(ctx context.Context) error {
	var errs []error

	if len(r.conf.Partners) == 0 {
		log.Println("No partners configured - skipping remote merge operation")
		return nil
	}

	for _, partner := range r.conf.Partners {
		if err := r.mergeRemotelyFromPartner(ctx, partner); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// Makes it mockable/testable
type getPairDataFunc func(context.Context, string, *[]entryPair) error
type getEntryDataFunc func(context.Context, string, string, *types.Entry) error

func (r Repository) mergeRemotelyFromPartner(ctx context.Context, partner string) error {
	getPair := func(ctx context.Context, partner string, pairs *[]entryPair) error {
		uri := fmt.Sprintf("%s/list", partner)
		return easyhttp.GetData(ctx, uri, r.conf.APIKey, pairs)
	}

	getEntry := func(ctx context.Context, partner, id types.EntryID, entry *types.Entry) error {
		uri := fmt.Sprintf("%s/get?id=%s", partner, id)
		return easyhttp.GetData(ctx, uri, r.conf.APIKey, entry)
	}

	return r.mergeFromPartner(ctx, partner, getPair, getEntry)
}

func (r Repository) mergeFromPartner(ctx context.Context, partner string,
	getPair getPairDataFunc, getEntry getEntryDataFunc) error {

	if err := r.load(); err != nil {
		return err
	}

	var (
		errs  []error
		pairs []entryPair
	)

	if err := getPair(ctx, partner, &pairs); err != nil {
		return err
	}

	for _, pair := range pairs {
		if r.hasSameEntry(pair) {
			continue
		}

		log.Println("pair", pair, "missing in local reposotory, going to merge it")

		var entry types.Entry
		if err := getEntry(ctx, partner, pair.ID, &entry); err != nil {
			errs = append(errs, err)
			continue
		}

		// In theory, this should never happen
		if pair.ID != entry.ID {
			errs = append(errs, fmt.Errorf("pair ID %s does not match entry id %s", pair.ID, entry.ID))
			continue
		}

		errs = append(errs, r.Merge(entry))
	}

	return errors.Join(errs...)
}
