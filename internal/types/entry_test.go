package types

import (
	"testing"
)

func oneEntry() (Entry, error) {
	ent := `
		{
			"body": "Body text here",
			"shared": {
				"Foo": { "Is": true },
				"Bar": { "Is": false }
			}
		}
	`
	return NewEntry([]byte(ent))
}

func anotherEntry() (Entry, error) {
	ent := `
		{
			"body": "Body text here",
			"shared": {
				"Foo": { "Is": true },
				"Bar": { "Is": true },
				"Baz": { "Is": false }
			}
		}
	`
	return NewEntry([]byte(ent))
}

func twoDifferentEntries() (ent1, ent2 Entry, err error) {
	if ent1, err = oneEntry(); err != nil {
		return
	}
	ent2, err = anotherEntry()
	return
}

func TestNewEntryFromJSON(t *testing.T) {
	ent1, err := oneEntry()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("ent1", ent1)
	if len(ent1.Shared) != 2 {
		t.Error("expected to have two shared entries in ent1")
	}
	if !ent1.IsShared("Foo") {
		t.Error("Foo should be shared")
	}
	if ent1.IsShared("Bar") {
		t.Error("Bar should not be shared")
	}
}

func TestEntryChecksum(t *testing.T) {
	t.Parallel()

	ent, err := NewEntry([]byte(`{"Body": "Body text here"}`))
	if err != nil {
		t.Error(err)
		return
	}

	expected := "4dbd4f04d7917b1f1bd0807cf39a260efe51085d49b40469fca27b7f89cc73bd"
	got := ent.Checksum()

	if expected != got {
		t.Errorf("expected checksum '%s' but got '%s'", expected, got)
		return
	}
	t.Log(ent.Checksum())
}

func TestEquals(t *testing.T) {
	t.Parallel()

	ent1, ent2, err := twoDifferentEntries()
	if err != nil {
		t.Error(err)
		return
	}

	if ent1.Equals(ent2) {
		t.Error("entries should not be equal", ent1, ent2)
	}

	t.Log("both entries differ", ent1, ent2)
}

func TestNewEntryFromCopy(t *testing.T) {
	ent1, _, err := twoDifferentEntries()
	if err != nil {
		t.Error(err)
	}

	ent2, err := NewEntryFromCopy(ent1)
	if err != nil {
		t.Error(err)
	}

	if !ent1.Equals(ent2) {
		t.Error("copy of entry ent1 does not equal")
		t.Error("original:", ent1)
		t.Error("copy:    ", ent2)
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	ent1, ent2, err := twoDifferentEntries()
	if err != nil {
		t.Error(err)
	}

	var changed bool
	if ent1, changed, err = ent1.Update(ent2); err != nil {
		t.Error(err)
	}

	if len(ent1.Shared) != 3 {
		t.Error("expected 3 entries after update", ent1)
	}

	if !changed {
		t.Error("expected the entry to be changed after update")
	}

	var sharedCount int
	for _, shared := range ent1.Shared {
		if shared.Is {
			sharedCount++
		}
	}

	if sharedCount != 2 {
		t.Error("expected 2 shared entries after update but got", sharedCount, ent1)
	}
}
