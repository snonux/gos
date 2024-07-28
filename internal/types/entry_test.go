package types

import "testing"

func TestEntryChecksum(t *testing.T) {
	t.Parallel()

	ent, err := NewEntry([]byte(`{"Body": "Body text here"}`))
	if err != nil {
		t.Error(err)
		return
	}

	expected := "e139c0788fbc0d9cce370e4918c1cbc8862184d9461bd1238c02b7f80cb042fe"
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

func TestUpdate(t *testing.T) {
	t.Parallel()

	ent1, ent2, err := twoDifferentEntries()
	if err != nil {
		t.Error(err)
	}

	if ent1.Changed {
		t.Error("didn't expect the entry to be changed before the update", ent1)
	}

	ent1, _ = ent1.Update(ent2)
	if len(ent1.Shared) != 3 {
		t.Error("expected 3 entries after update", ent1)
	}

	if !ent1.Changed {
		t.Error("expected the entry to be changed after update")
	}

	var isShared int
	for _, shared := range ent1.Shared {
		if shared.Is {
			isShared++
		}
	}

	if isShared != 2 {
		t.Error("expected 2 shared entries after update but got", isShared, ent1)
	}
}

func twoDifferentEntries() (ent1, ent2 Entry, err error) {
	ent1Str := `
		{
			"Body": "Body text here",
			"Shared": [
				{ "Name": "Foo", "Is": true },
				{ "Name": "Bar", "Is": false }
			]
		}
	`
	ent1, err = NewEntry([]byte(ent1Str))
	if err != nil {
		return
	}

	ent2Str := `
		{
			"Body": "Body text here",
			"Shared": [
				{ "Name": "Foo", "Is": true },
				{ "Name": "Bar", "Is": true },
				{ "Name": "Baz", "Is": false }
			]
		}
	`
	ent2, err = NewEntry([]byte(ent2Str))
	return
}
