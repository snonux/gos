package types

import "testing"

func TestEntryChecksum(t *testing.T) {
	t.Parallel()

	ent, err := NewEntry([]byte(`{"Body": "Body text here"}`))
	if err != nil {
		t.Error(err)
		return
	}

	expected := "8618a63380fe6d365422cae6ef143a88bb6bd78df567fea3822074cc748f52f8"
	got := ent.Checksum()

	if expected != got {
		t.Errorf("expected checksum '%s' but got '%s'", expected, got)
		return
	}
	t.Log(ent.Checksum())
}
