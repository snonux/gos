package repository

import (
	"testing"

	"codeberg.org/snonux/gos/internal/types"
	"codeberg.org/snonux/gos/internal/vfs"
)

func TestRepositoryGet(t *testing.T) {
	t.Parallel()

	entry, _, err := twoDifferentEntries()
	if err != nil {
		t.Error(err)
		return
	}

	repo.put(entry)
	t.Log(fs)

	entryGot, err := repo.Get(entry.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if !entryGot.Equals(entry) {
		t.Error("expected to get", entry, "but got", entryGot)
	}
}

// TODO: Write unit tests for the remainder of the repo methods

func setupRepository() (repo Repository, entry1, entry2 types.Entry, err error) {
	fs := make(vfs.MemoryFS)
	repo = newRepository("./data", fs)

	entry1Str := `
		{
			"Body": "Body text here",
			"Shared": [
				{ "Name": "Foo", "Is": true },
				{ "Name": "Bar", "Is": false }
			]
		}
	`
	entry1, err = types.NewEntry([]byte(entry1Str), fs)
	if err != nil {
		return
	}

	entry2Str := `
		{
			"Body": "Body text here",
			"Shared": [
				{ "Name": "Foo", "Is": true },
				{ "Name": "Bar", "Is": true },
				{ "Name": "Baz", "Is": false }
			]
		}
	`
	entry2, err = types.NewEntry([]byte(entry2Str), fs)
	return
}
