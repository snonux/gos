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

	pairs, _ := repo.List()
	if len(entries) != len(pairs) {
		t.Error("expected as many entries as pairs")
	}

	for _, entry := range entries {
		var found bool
		for _, pair := range pairs {
			if entry.ID == pair.ID && entry.Checksum() == pair.Checksum {
				found = true
				t.Log("entry matches pair", entry, pair)
				break
			}
		}
		if !found {
			t.Error("could not find entry", entry, "in", pairs)
		}
	}
}

func TestRepositoryHasSameEntry(t *testing.T) {
	t.Parallel()

	fs := make(vfs.MemoryFS)
	repo := newRepository("./data", fs)
	entry, _ := makeAnEntry()
	_ = repo.put(entry)

	pair := EntryPair{entry.ID, entry.Checksum()}
	if !repo.HasSameEntry(pair) {
		t.Error("repo does not contain entry corresponding to pair", pair)
	}

	pair = EntryPair{"nonexistent", "nonexistent"}
	if repo.HasSameEntry(pair) {
		t.Error("repo does contain entry corresponding to pair", pair, "but that should not be")
	}
}

func TestRepositoryMerge(t *testing.T) {
	t.Parallel()

	fs := make(vfs.MemoryFS)
	repo := newRepository("./data", fs)
	entry1, _ := makeAnEntry()
	_ = repo.put(entry1)

	entry2, _ := makeAnotherEntry()
	// Need to have the same IDs so that the entries will actually be merged
	entry2.ID = entry1.ID
	// Merge a modified entry2 into the repository.
	entry2.Body = "merged"
	entry2.Epoch = 12345
	_ = repo.Merge(entry2)

	pairs, _ := repo.List()
	// Ensuring the merge didn't add a new entry
	if len(pairs) != 1 {
		t.Error("expected exactly one element in the repo but got", pairs)
	}

	entryGot, _ := repo.Get(entry1.ID)
	if entryGot.Body != "merged" {
		t.Error("unexpected body", entryGot.Body)
	}
	if entryGot.Epoch != 12345 {
		t.Error("unexpected epoch", entryGot.Epoch)
	}
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
