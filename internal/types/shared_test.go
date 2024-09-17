package types

import (
	"testing"
	"time"
)

func TestShared(t *testing.T) {
	t.Parallel()

	before := time.Now().Unix()
	shared := newShared(false)

	if before > shared.Timestamp {
		t.Errorf("expected %d to be after or equal %d", shared.Timestamp, before)
	}

	after := time.Now().Unix()
	if after < shared.Timestamp {
		t.Errorf("expected %d to be before or equal %d", shared.Timestamp, before)
	}
}
