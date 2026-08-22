package tags

import (
	"slices"
	"strings"

	"github.com/snonux/gos/internal/config"
	"github.com/snonux/gos/internal/platforms"
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
				raw := strings.ToLower(t[1:])
				if p, err := platforms.New(raw); err == nil {
					raw = p.String()
				}
				s.Excludes = append(s.Excludes, raw)
			} else {
				raw := strings.ToLower(t)
				if p, err := platforms.New(raw); err == nil {
					raw = p.String()
				}
				s.Includes = append(s.Includes, raw)
			}
		}
	}

	// If there is no share tag, by default include all platforms but noop-like
	// platforms (Noop, Snonux) that are only queued when explicitly named in share:.
	if len(s.Includes) == 0 {
		for platformStr := range args.Platforms {
			p, err := platforms.New(platformStr)
			if err != nil {
				continue
			}
			name := p.String()
			if name == "noop" || name == "snonux" {
				continue
			}
			if slices.Contains(s.Excludes, name) {
				continue
			}
			s.Includes = append(s.Includes, name)
		}
	}

	return s, nil
}

func (s Share) Excluded(platformStr string) bool {
	return slices.Contains(s.Excludes, platformStr) || !slices.Contains(s.Includes, platformStr)
}
