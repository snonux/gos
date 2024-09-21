package repository

import "codeberg.org/snonux/gos/internal/types"

// Keeps track of how many messages were posted to social media over the last week and month.
type stats struct {
	// Sliding window of entries shared last 7 days
	last7Days map[types.PlatformName][]types.UnixEpoch
	// Sliding window of entries shared last 30 days
	last30Days map[types.PlatformName][]types.UnixEpoch
}

func newStats() stats {
	return stats{
		last7Days:  make(map[types.PlatformName][]types.UnixEpoch),
		last30Days: make(map[types.PlatformName][]types.UnixEpoch),
	}
}

// func (s stats) add(platform types.PlatformName, entry types.Entry) {

// }
