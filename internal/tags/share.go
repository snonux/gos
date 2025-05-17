package tags

import (
	"slices"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
)

// Share tags.
type Share struct {
	Includes []string // The platforms to include
	Excludes []string // The platforms to exclude
}

func NewShare(args config.Args, tags map[string]struct{}) (Share, error) {
	var s Share

	for tag := range tags {
		if !strings.HasPrefix(tag, "share:") {
			continue
		}
		for _, t := range strings.Split(tag[6:], ":") {
			if strings.HasPrefix(t, "-") {
				s.Excludes = append(s.Excludes, strings.ToLower(t[1:]))
			} else {
				s.Includes = append(s.Includes, strings.ToLower(t))
			}
		}
	}

	// If there is no share tag, by default include all platforms but "Noop"
	if len(s.Includes) == 0 {
		for platformStr := range args.Platforms {
			if slices.Contains(s.Excludes, strings.ToLower(platformStr)) {
				continue
			}
			if platformStr == "Noop" {
				continue
			}
			s.Includes = append(s.Includes, strings.ToLower(platformStr))
		}
	}

	return s, nil
}

func (s Share) Excluded(platformStr string) bool {
	return slices.Contains(s.Excludes, platformStr) || !slices.Contains(s.Includes, platformStr)
}
