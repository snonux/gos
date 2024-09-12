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

func TestPendingDelete(t *testing.T) {
	pending := newPending()
	pending.add(types.LinkedIn, "fooid")

	entries, ok := pending.get(types.LinkedIn)
	if !ok {
		t.Error("expected ok return status")
	}
	if len(entries) != 1 {
		t.Error("expected one entry")
	}

	pending.delete(types.LinkedIn, "fooid")
	if entries, ok = pending.get(types.LinkedIn); ok {
		t.Error("expected not an ok", entries)
	}
}

func TestPendingNext(t *testing.T) {
	pending := newPending()

	id, ok := pending.next(types.LinkedIn)
	if ok {
		t.Error("not expected ok return status", id)
	}

	pending.add(types.LinkedIn, "fooid")
	id, ok = pending.next(types.LinkedIn)
	if !ok {
		t.Error("expected ok return status")
	}
	if id != "fooid" {
		t.Error("expected entry ID fooid")
	}
}
