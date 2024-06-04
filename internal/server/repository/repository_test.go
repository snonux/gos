package repository

import (
	"testing"

	"codeberg.org/snonux/gos/internal/types"
	"codeberg.org/snonux/gos/internal/vfs"
)

func TestRepositoryGet(t *testing.T) {
	t.Parallel()

	repo, entry, _, err := setupRepository()
	if err != nil {
		t.Error(err)
		return
	}

	if err := repo.put(entry); err != nil {
		t.Error(err)
		return
	}
	list, err := repo.List()
	if err != err {
		t.Error(err)
		return
	}
	t.Log(string(list))

	entryGot, ok := repo.Get(entry.ID)
	if !ok {
		t.Errorf("could not find entry with id %s in repo", entry.ID)
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
