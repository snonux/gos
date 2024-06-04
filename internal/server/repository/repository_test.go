package repository

import (
	"testing"

	"codeberg.org/snonux/gos/internal/types"
	"codeberg.org/snonux/gos/internal/vfs"
)

func TestRepositoryPutGet(t *testing.T) {
	t.Parallel()

	fs := make(vfs.MemoryFS)
	repo := newRepository("./data", fs)

	for _, entry := range makeEntries(t) {
		t.Run(entry.ID, func(t *testing.T) {
			_ = repo.put(entry)
			entryGot, ok := repo.Get(entry.ID)
			if !ok {
				t.Errorf("could not find entry with id %s in repo", entry.ID)
			}
			if !entryGot.Equals(entry) {
				t.Error("expected to get", entry, "but got", entryGot)
			}
		})
	}
}

func TestRepositoryLoad(t *testing.T) {
	t.Parallel()

	fs := make(vfs.MemoryFS)
	repo := newRepository("./data", fs)
	entries := makeEntries(t)

	// Write entries into the VFS
	for _, entry := range entries {
		bytes, _ := entry.Serialize()
		_ = repo.fs.WriteFile(repo.entryPath(entry), bytes)
	}

	// Load entries from VFS into the repo
	if err := repo.load(); err != nil {
		t.Error(err)
	}

	for _, entry := range entries {
		t.Run(entry.ID, func(t *testing.T) {
			entryGot, ok := repo.Get(entry.ID)
			if !ok {
				t.Errorf("could not find entry with id %s in repo", entry.ID)
			}
			if !entryGot.Equals(entry) {
				t.Error("expected to get", entry, "but got", entryGot)
			}
		})
	}
}

func TestRepositoryList(t *testing.T) {
	t.Parallel()

	fs := make(vfs.MemoryFS)
	repo := newRepository("./data", fs)
	entries := makeEntries(t)

	for _, entry := range entries {
		_ = repo.put(entry)
	}

	list, _ := repo.List()
	t.Log(string(list))
}

func makeEntries(t *testing.T) []types.Entry {
	entry1, err := makeAnEntry()
	if err != nil {
		t.Error(err)
	}
	entry2, err := makeAnotherEntry()
	if err != nil {
		t.Error(err)
	}
	return []types.Entry{entry1, entry2}
}

func makeAnEntry() (types.Entry, error) {
	entry := `
		{
			"Body": "Body text here",
			"Shared": [
				{ "Name": "Foo", "Is": true },
				{ "Name": "Bar", "Is": false }
			]
		}
	`
	return types.NewEntry([]byte(entry))
}

func makeAnotherEntry() (types.Entry, error) {
	entry := `
		{
			"Body": "Another text here",
			"Shared": [
				{ "Name": "Foo", "Is": true },
				{ "Name": "Bar", "Is": true },
				{ "Name": "Baz", "Is": false }
			]
		}
	`
	return types.NewEntry([]byte(entry))
}

// TODO: Write unit tests for the remainder of the repo methods
