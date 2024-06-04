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

	entry1, _ := makeAnEntry(fs)
	entry2, _ := makeAnotherEntry(fs)
	entries := []types.Entry{entry1, entry2}

	for _, entry := range entries {
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

func makeAnEntry(fs fs) (types.Entry, error) {
	entry := `
		{
			"Body": "Body text here",
			"Shared": [
				{ "Name": "Foo", "Is": true },
				{ "Name": "Bar", "Is": false }
			]
		}
	`
	return types.NewEntry([]byte(entry), fs)
}

func makeAnotherEntry(fs fs) (types.Entry, error) {
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
	return types.NewEntry([]byte(entry), fs)
}

// TODO: Write unit tests for the remainder of the repo methods
