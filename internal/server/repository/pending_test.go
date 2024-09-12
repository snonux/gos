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

	pending.add(types.LinkedIn, "fooid")
	pending.add(types.LinkedIn, "barid")

	entries, ok = pending.get(types.LinkedIn)
	if !ok {
		t.Error("expected ok return status")
	}
	if len(entries) != 2 {
		t.Error("expected two entries")
	}
}
