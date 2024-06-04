package types

import "testing"

func TestEntryChecksum(t *testing.T) {
	t.Parallel()

	entry, err := NewEntry([]byte(`{"Body": "Body text here"}`))
	if err != nil {
		t.Error(err)
		return
	}

	expected := "e139c0788fbc0d9cce370e4918c1cbc8862184d9461bd1238c02b7f80cb042fe"
	got := entry.Checksum()

	if expected != got {
		t.Errorf("expected checksum '%s' but got '%s'", expected, got)
		return
	}
	t.Log(entry.Checksum())
}

func twoDifferentEntries() (entry1, entry2 Entry, err error) {
	entry1Str := `
		{
			"Body": "Body text here",
			"Shared": [
				{ "Name": "Foo", "Is": true },
				{ "Name": "Bar", "Is": false }
			]
		}
	`
	entry1, err = NewEntry([]byte(entry1Str))
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
	entry2, err = NewEntry([]byte(entry2Str))
	return
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

func TestUpdate(t *testing.T) {
	t.Parallel()

	entry1, entry2, err := twoDifferentEntries()
	if err != nil {
		t.Error(err)
		return
	}

	entry1, _ = entry1.Update(entry2)
	if len(entry1.Shared) != 3 {
		t.Error("expected 3 entries after update", entry1)
		return
	}

	var isShared int
	for _, shared := range entry1.Shared {
		if shared.Is {
			isShared++
		}
	}

	if isShared != 2 {
		t.Error("expected 2 shared entries after update but got", isShared, entry1)
		return
	}

	t.Log("entry as expected after update", entry1)
}
