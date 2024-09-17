package repository

import "codeberg.org/snonux/gos/internal/types"

type pendingEntries map[types.EntryID]struct{}

// Keep track of pending entries per social platform
type pending struct {
	platforms map[types.PlatformName]pendingEntries
}

func newPending() pending {
	return pending{make(map[types.PlatformName]pendingEntries)}
}

// Returns number of pending entries for the platform
// func (p pending) num(platform types.PlatformName) int {
// 	pe, ok := p.platforms[platform]
// 	if !ok {
// 		return 0
// 	}
// 	return len(pe)
// }

func (p pending) add(platform types.PlatformName, id types.EntryID) {
	pe, ok := p.platforms[platform]
	if !ok {
		pe = make(pendingEntries)
	}
	pe[id] = struct{}{}
	p.platforms[platform] = pe
}

func (p pending) delete(platform types.PlatformName, id types.EntryID) {
	pe, ok := p.platforms[platform]
	if !ok {
		return
	}
	delete(pe, id)
	p.platforms[platform] = pe
}

func (p pending) get(platform types.PlatformName) (pendingEntries, bool) {
	pe, ok := p.platforms[platform]
	return pe, ok && len(pe) > 0
}

func (p pending) next(platform types.PlatformName) (types.EntryID, bool) {
	pe, ok := p.get(platform)
	if !ok {
		return "", false
	}
	for id := range pe {
		return id, true
	}
	return "", false
}
