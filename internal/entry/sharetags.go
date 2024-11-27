package entry

import (
	"slices"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
)

// TODO: Own package only dealing with tags, and put all tag code in there.
type shareTags struct {
	includes []string // The platforms to include
	excludes []string // The platforms to exclude
}

func newShareTags(args config.Args, tags map[string]struct{}) (shareTags, error) {
	var s shareTags

	for tag := range tags {
		if !strings.HasPrefix(tag, "share:") {
			continue
		}
		for _, t := range strings.Split(tag[6:], ":") {
			if strings.HasPrefix(t, "-") {
				s.excludes = append(s.excludes, strings.ToLower(t[1:]))
			} else {
				s.includes = append(s.includes, strings.ToLower(t))
			}
		}
	}

	if len(s.includes) == 0 {
		for platformStr := range args.Platforms {
			if slices.Contains(s.excludes, strings.ToLower(platformStr)) {
				continue
			}
			s.includes = append(s.includes, strings.ToLower(platformStr))
		}
	}

	return s, nil
}
