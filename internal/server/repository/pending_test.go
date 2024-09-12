package repository

import (
	"testing"

	"codeberg.org/snonux/gos/internal/types"
)

func TestPendingAdd(t *testing.T) {
	pending := newPending()

	entries, ok := pending.get(types.LinkedIn)
	if ok {
		t.Error("expected no ok return status")
	}
	if len(entries) != 0 {
		t.Error("expected no entries")
	}

	// TODO REFACTOR: Use constants for types.PlatformName's
	// TODO REFACTOR: Don't use a type alias for types.PlatformName anymore, but an own type.
	pending.add(types.LinkedIn, "foo")
	pending.add(types.LinkedIn, "bar")

	entries, ok = pending.get(types.LinkedIn)
	if !ok {
		t.Error("expected ok return status")
	}
	if len(entries) != 2 {
		t.Error("expected two entries")
	}
}
