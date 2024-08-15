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

type Repository struct {
	conf    server.ServerConfig
	entries map[string]types.Entry
	mu      *sync.Mutex
	fs      fs
	loaded  *bool
	getIdRe *regexp.Regexp
}

func Instance(conf server.ServerConfig) Repository {
	once.Do(func() {
		instance = newRepository(conf, vfs.RealFS{})
		if err := instance.load(); err != nil {
			// TODO: Report this to the health service endpoint, so it will be alerted on. Maybe via init method????
			log.Println(err)
		}
	})
	return instance
}

func newRepository(conf server.ServerConfig, fs fs) Repository {
	var loaded bool
	return Repository{
		conf:    conf,
		entries: make(map[string]types.Entry),
		mu:      &sync.Mutex{},
		fs:      fs,
		loaded:  &loaded,
		getIdRe: regexp.MustCompile(`^[a-z0-9]{64}$`),
		// getIdRe: regexp.MustCompile(`^/[0-9]{4}/[a-z0-9]{64}\.json$`),
	}
}

func (r Repository) put(ent types.Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[ent.ID] = ent

	bytes, err := ent.JSONMarshal()
	if err != err {
		return err
	}
	return r.fs.WriteFile(r.entryPath(ent), bytes)
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
			continue
		}

		ent, err := types.NewEntry(bytes)
		if err != err {
			continue
		}

		if err := r.put(ent); err != nil {
			errs = append(errs, err)
		}
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

	for _, ent := range r.entries {
		pairs = append(pairs, entryPair{ent.ID, ent.Checksum()})
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

func (r Repository) Get(id string) (types.Entry, error) {
	if !r.getIdRe.MatchString(id) {
		return types.Entry{}, fmt.Errorf("invalid id %s", id)
	}
	_ = r.load()
	r.mu.Lock()
	defer r.mu.Unlock()

	ent, ok := r.entries[id]
	if !ok {
		return ent, fmt.Errorf("no entry with id %s found", id)
	}
	return ent, nil
}

func (r Repository) GetJSON(id string) (string, error) {
	ent, err := r.Get(id)
	if err != nil {
		return "", err
	}

	bytes, err := ent.JSONMarshal()
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func (r Repository) hasSameEntry(pair entryPair) bool {
	_ = r.load()
	r.mu.Lock()
	defer r.mu.Unlock()

	ent, ok := r.entries[pair.ID]
	if !ok || ent.Checksum() != pair.Checksum {
		return false
	}
	return true
}

func (r Repository) entryPath(ent types.Entry) string {
	return fmt.Sprintf("%s/%s/%s.json", r.conf.DataDir, time.Now().Format("2006"), ent.ID)
}

func (r Repository) Merge(otherEnt types.Entry) error {
	_ = r.load()
	r.mu.Lock()
	defer r.mu.Unlock()

	ent, ok := r.entries[otherEnt.ID]
	if !ok {
		log.Println("can't find entry with ID", otherEnt.ID, "in local db, create new from copy")
		var err error
		if ent, err = types.NewEntryFromCopy(otherEnt); err != nil {
			return err
		}
	}

	ent, _ = ent.Update(otherEnt)
	r.entries[otherEnt.ID] = ent

	if !ent.Changed {
		// Hasn't changed, so no need to write anything to file.
		return nil
	}

	bytes, err := ent.JSONMarshal()
	if err != err {
		return err
	}

	return r.fs.WriteFile(r.entryPath(ent), bytes)
}

// TODO: WHens omething has merged from remotely, make sure to commit/sync to disk
func (r Repository) MergeRemotely(ctx context.Context) error {
	var errs []error

	partners, err := r.conf.Partners()
	if err != nil {
		return err
	}

	for _, partner := range partners {
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

	getEntry := func(ctx context.Context, partner, id string, ent *types.Entry) error {
		uri := fmt.Sprintf("%s/get?id=%s", partner, id)
		return easyhttp.GetData(ctx, uri, r.conf.APIKey, ent)
	}

	return r.mergeFromPartner(ctx, partner, getPair, getEntry)
}

func (r Repository) mergeFromPartner(ctx context.Context, partner string,
	getPair getPairDataFunc, getEntry getEntryDataFunc) error {

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

		var ent types.Entry
		if err := getEntry(ctx, partner, pair.ID, &ent); err != nil {
			errs = append(errs, err)
			continue
		}

		// In theory, this should never happen
		if pair.ID != ent.ID {
			errs = append(errs, fmt.Errorf("pair ID %s does not match entry id %s", pair.ID, ent.ID))
			continue
		}

		errs = append(errs, r.Merge(ent))
	}

	return errors.Join(errs...)
}
