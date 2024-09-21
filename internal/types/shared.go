package types

import (
	"fmt"
	"time"
)

type UnixEpoch int64

// Tells me whether the entry was Shared to the sm platform named Name
type Shared struct {
	Is        bool      `json:"is,omitempty"`
	Timestamp UnixEpoch `json:"timestamp,omitempty"`
}

func newShared(is bool) Shared {
	return Shared{
		Is:        is,
		Timestamp: UnixEpoch(time.Now().Unix()),
	}
}

func (s Shared) String() string {
	return fmt.Sprintf("Is:%v;Timestamp:%v", s.Is, s.Timestamp)
}

func (s Shared) Equals(other Shared) bool {
	return s.Timestamp == other.Timestamp && s.Is == other.Is
}
