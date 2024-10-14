package queue

import (
	"strings"

	"codeberg.org/snonux/gos/internal/config"
)

type shareTags struct {
	// The platforms to include
	includes []string
	// The platforms to exclude
	excludes []string
}

// Valid tags are: share:foo[,...]
// whereas foo can be a supported plutform such as linkedin, mastodon, etc.
// foo can also be prefixed with - to exclude it. See unit tests for examples.
func newShareTags(args config.Args, filePath string) shareTags {
	var s shareTags

	parts := strings.Split(filePath, ".")
	tagStr := parts[len(parts)-2]
	if len(parts) > 2 && strings.HasPrefix(tagStr, "share:") {
		for _, tag := range strings.Split(tagStr[6:], ":") {
			if strings.HasPrefix(tag, "-") {
				s.excludes = append(s.excludes, tag[1:])
			} else {
				s.includes = append(s.includes, tag)
			}
		}
	}

	if len(s.includes) == 0 && len(s.excludes) == 0 {
		// If nothing found, include all of them
		s.includes = args.Platforms
	}

	return s
}
