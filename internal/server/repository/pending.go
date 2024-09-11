package repository

import "codeberg.org/snonux/gos/internal/types"

type pendingEntries []types.EntryID

// Keep track of pending entries per social platform
type pending struct {
	platforms map[types.PlatformName]pendingEntries
}

func newPending() pending {
	return pending{platforms: make(map[types.PlatformName]pendingEntries)}
}

func (p pending) add(platform types.PlatformName, id types.EntryID) {
	pe, _ := p.get(platform)
	p.platforms[platform] = append(pe, id)
}

func (p pending) get(platform types.PlatformName) (pendingEntries, bool) {
	pe, ok := p.platforms[platform]
	return pe, ok
}
