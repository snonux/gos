package types

import "testing"

func TestEntryChecksum(t *testing.T) {
	t.Parallel()

	entry, err := NewEntry([]byte(`{"Body": "Body text here"}`))
	if err != nil {
		t.Error(err)
		return
	}

	expected := "8618a63380fe6d365422cae6ef143a88bb6bd78df567fea3822074cc748f52f8"
	got := entry.Checksum()

	if expected != got {
		t.Errorf("expected checksum '%s' but got '%s'", expected, got)
		return
	}
	t.Log(entry.Checksum())
}

func twoDifferentEntries(t *testing.T) (entry1, entry2 Entry, err error) {
	entry1, err = NewEntry([]byte(`{"Body": "Body text here"}`))
	if err != nil {
		return
	}

	entry2, err = NewEntry([]byte(`{"Body": "Body text here 2"}`))
	return
}

func TestEquals(t *testing.T) {
	t.Parallel()

	entry1, entry2, err := twoDifferentEntries(t)
	if err != nil {
		t.Error(err)
		return
	}

	if entry1.Equals(entry2) {
		t.Error("entries should not be equal", entry1, entry2)
	}

	t.Log("both entries differ", entry1, entry2)
}
