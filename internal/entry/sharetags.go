package entry

import (
	"fmt"
	"slices"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
)

type shareTags struct {
	includes []string // The platforms to include
	excludes []string // The platforms to exclude
}

func newShareTags(args config.Args, filePath string) (shareTags, error) {
	var s shareTags

	parts := strings.Split(filePath, ".")
	if len(parts) < 4 {
		return s, fmt.Errorf("invalid file path: %s", filePath)
	}
	tagStr := parts[len(parts)-4]

	if len(parts) > 2 && strings.HasPrefix(tagStr, "share:") {
		for _, tag := range strings.Split(tagStr[6:], ":") {
			if strings.HasPrefix(tag, "-") {
				s.excludes = append(s.excludes, strings.ToLower(tag[1:]))
			} else {
				s.includes = append(s.includes, strings.ToLower(tag))
			}
		}
	}

	if len(s.includes) == 0 {
		for platform := range args.Platforms {
			if slices.Contains(s.excludes, strings.ToLower(platform)) {
				continue
			}
			s.includes = append(s.includes, strings.ToLower(platform))
		}
	}

	return s, nil
}
