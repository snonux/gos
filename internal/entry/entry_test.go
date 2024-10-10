package entry

import (
	"fmt"
	"testing"

	"codeberg.org/snonux/gos/internal/timestamp"
)

func TestEntry(t *testing.T) {
	states := []State{Queued, Posted}
	stamps := []string{"20240928-111835", "20241028-120135"}

	for _, state := range states {
		for _, stamp := range stamps {
			queuedPath := fmt.Sprintf("gosdir/db/platforms/linkedin/helloworld.txt.%s.%s", stamp, state)

			ent, err := New(queuedPath)
			if err != nil {
				t.Error(err)
			}
			if ent.Path != queuedPath {
				t.Errorf("expected path %s but got %s", queuedPath, ent.Path)
			}
			if ent.State != state {
				t.Errorf("expected state %s but got %s", state, ent.State)
			}

			expectedTime, err := timestamp.Parse(stamp)
			if err != nil {
				t.Error(err)
			}
			if ent.Time != expectedTime {
				t.Errorf("expected time to be %v but got %v", expectedTime, ent.Time)
			}
		}
	}
}
