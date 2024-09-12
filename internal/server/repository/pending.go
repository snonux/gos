package repository

import "codeberg.org/snonux/gos/internal/types"

type pendingEntries map[types.EntryID]struct{}

// Keep track of pending entries per social platform
type pending struct {
	platforms map[types.PlatformName]pendingEntries
}

func newPending() pending {
	return pending{platforms: make(map[types.PlatformName]pendingEntries)}
}

func (p pending) add(platform types.PlatformName, id types.EntryID) {
	pe, ok := p.get(platform)
	if !ok {
		pe = make(pendingEntries)
	}
	pe[id] = struct{}{}
	p.platforms[platform] = pe
}

func (p pending) get(platform types.PlatformName) (pendingEntries, bool) {
	pe, ok := p.platforms[platform]
	return pe, ok
}
