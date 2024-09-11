package repository

import "testing"

func TestPendingAdd(t *testing.T) {
	pending := newPending()

	entries, ok := pending.get("LinkedIn")
	if ok {
		t.Error("expected no ok return status")
	}
	if len(entries) != 0 {
		t.Error("expected no entries")
	}

	// TODO REFACTOR: Use constants for types.PlatformName's
	// TODO REFACTOR: Don't use a type alias for types.PlatformName anymore, but an own type.
	pending.add("LinkedIn", "foo")
	pending.add("LinkedIn", "bar")

	entries, ok = pending.get("LinkedIn")
	if !ok {
		t.Error("expected ok return status")
	}
	if len(entries) != 2 {
		t.Error("expected two entries")
	}
}
