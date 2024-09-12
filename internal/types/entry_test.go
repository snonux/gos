package types

import (
	"testing"
)

func TestNewEntryFromJSON(t *testing.T) {
	entry1, err := oneEntry()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("entry1", entry1)
	if len(entry1.Shared) != 2 {
		t.Error("expected to have two shared entries in entry1")
	}
	if !entry1.IsShared("Mastodon") {
		t.Error("Mastodon should be shared")
	}
	if entry1.IsShared("LinkedIn") {
		t.Error("LinkedIn should not be shared")
	}
}

func TestEntryChecksum(t *testing.T) {
	t.Parallel()

	entry, err := NewEntry([]byte(`{"Body": "Body text here"}`))
	if err != nil {
		t.Error(err)
		return
	}

	expected := "4dbd4f04d7917b1f1bd0807cf39a260efe51085d49b40469fca27b7f89cc73bd"
	got := entry.Checksum()

	if expected != got {
		t.Errorf("expected checksum '%s' but got '%s'", expected, got)
		return
	}
	t.Log(entry.Checksum())
}

func TestEquals(t *testing.T) {
	t.Parallel()

	entry1, entry2, err := twoDifferentEntries()
	if err != nil {
		t.Error(err)
		return
	}

	if entry1.Equals(entry2) {
		t.Error("entries should not be equal", entry1, entry2)
	}

	t.Log("both entries differ", entry1, entry2)
}

func TestNewEntryFromCopy(t *testing.T) {
	entry1, _, err := twoDifferentEntries()
	if err != nil {
		t.Error(err)
	}

	entry2, err := NewEntryFromCopy(entry1)
	if err != nil {
		t.Error(err)
	}

	if !entry1.Equals(entry2) {
		t.Error("copy of entry entry1 does not equal")
		t.Error("original:", entry1)
		t.Error("copy:    ", entry2)
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	entry1, entry2, err := twoDifferentEntries()
	if err != nil {
		t.Error(err)
	}

	var changed bool
	if entry1, changed, err = entry1.Update(entry2); err != nil {
		t.Error(err)
	}

	if len(entry1.Shared) != 3 {
		t.Error("expected 3 entries after update", entry1)
	}

	if !changed {
		t.Error("expected the entry to be changed after update")
	}

	var sharedCount int
	for _, shared := range entry1.Shared {
		if shared.Is {
			sharedCount++
		}
	}

	if sharedCount != 2 {
		t.Error("expected 2 shared entries after update but got", sharedCount, entry1)
	}
}

func oneEntry() (Entry, error) {
	entry := `
		{
			"body": "Body text here",
			"shared": {
				"Mastodon": { "is": true },
				"LinkedIn": { "is": false }
			}
		}
	`
	return NewEntry([]byte(entry))
}

func anotherEntry() (Entry, error) {
	entry := `
		{
			"body": "Body text here",
			"shared": {
				"Mastodon": { "is": true },
				"LinkedIn": { "is": true },
				"Textfile": { "is": false }
			}
		}
	`
	return NewEntry([]byte(entry))
}

func twoDifferentEntries() (entry1, entry2 Entry, err error) {
	if entry1, err = oneEntry(); err != nil {
		return
	}
	entry2, err = anotherEntry()
	return
}
