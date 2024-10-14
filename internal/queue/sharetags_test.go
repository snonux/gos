package queue

import (
	"slices"
	"testing"

	"codeberg.org/snonux/gos/internal/config"
)

func TestShareTagsPositive(t *testing.T) {
	t.Parallel()

	args := config.Args{Platforms: []string{"mastodon", "linkedin"}}
	testTable := map[string]shareTags{
		"./foo/bar.without.tags.txt": {
			includes: args.Platforms, // No tags: default platforms
		},
		"./foo/bar.share:linkedin.txt": {
			includes: []string{"linkedin"},
		},
		"./foo/bar.share:-linkedin.txt": {
			excludes: []string{"linkedin"},
		},
		"./foo/bar.share:linkedin:mastodon.txt": {
			includes: []string{"linkedin", "mastodon"},
		},
		"./foo/bar.share:linkedin:-mastodon:xcom.txt": {
			includes: []string{"linkedin", "xcom"},
			excludes: []string{"mastodon"},
		},
	}

	for filePath, expectedResult := range testTable {
		t.Run(filePath, func(t *testing.T) {
			shareTags := newShareTags(args, filePath)
			if !slices.Equal(shareTags.includes, expectedResult.includes) {
				t.Errorf("Expected includes to be %v but got %v", expectedResult.includes, shareTags.includes)
			}
			if !slices.Equal(shareTags.excludes, expectedResult.excludes) {
				t.Errorf("Expected excludes to be %v but got %v", expectedResult.excludes, shareTags.excludes)
			}
		})

	}
}
func TestShareTagsNegative(t *testing.T) {
	t.Parallel()

	args := config.Args{Platforms: []string{"mastodon", "linkedin"}}
	testTable := map[string]shareTags{
		"./foo/bar.without.tags.txt": {
			includes: []string{"linkedin"},
		},
		"./foo/bar.share:linkedin.txt": {
			includes: []string{"mastodon"},
		},
		"./foo/bar.share:-linkedin.txt": {
			includes: []string{"linkedin"},
		},
		"./foo/bar.share:linkedin:mastodon.txt": {
			includes: []string{"oups", "mastodon"},
		},
		"./foo/bar.share:linkedin:-mastodon:xcom.txt": {
			includes: []string{"linkedin", "xcom"},
			excludes: []string{"mastodon", "xcom"},
		},
	}

	for filePath, unexpectedResult := range testTable {
		t.Run(filePath, func(t *testing.T) {
			shareTags := newShareTags(args, filePath)
			if slices.Equal(shareTags.includes, unexpectedResult.includes) &&
				slices.Equal(shareTags.excludes, unexpectedResult.excludes) {
				t.Errorf("expected %v not to be the actual result", unexpectedResult)
			}
		})

	}
}

func TestShareTagsIsIncluded(t *testing.T) {
	t.Parallel()

	assertIncluded := func(shareTags shareTags, platforms ...string) {
		for _, platform := range platforms {
			if !shareTags.IsIncluded(platform) {
				t.Errorf("expected %s included in %v", platform, shareTags)
			}
			if shareTags.IsExcluded(platform) {
				t.Errorf("expected %s not to be excluded in %v", platform, shareTags)
			}
		}
	}
	assertExcluded := func(shareTags shareTags, platforms ...string) {
		for _, platform := range platforms {
			if shareTags.IsIncluded(platform) {
				t.Errorf("expected %s not to be included in %v", platform, shareTags)
			}
			if !shareTags.IsExcluded(platform) {
				t.Errorf("expected %s to be excluded in %v", platform, shareTags)
			}
		}
	}
	args := config.Args{Platforms: []string{"mastodon", "linkedin"}}

	filePath := "foo/bar/baz.txt"
	t.Run(filePath, func(t *testing.T) {
		assertIncluded(newShareTags(args, filePath), "mastodon", "linkedin")
	})

	filePath = "foo/bar/baz.share:mastodon.txt"
	t.Run(filePath, func(t *testing.T) {
		assertIncluded(newShareTags(args, filePath), "mastodon")
		assertExcluded(newShareTags(args, filePath), "linkedin")
	})

	filePath = "foo/bar/baz.share:mastodon.txt"
	t.Run(filePath, func(t *testing.T) {
		assertIncluded(newShareTags(args, filePath), "mastodon")
		assertExcluded(newShareTags(args, filePath), "linkedin")
	})

	filePath = "foo/bar/baz.share:linkedin:mastodon.txt"
	t.Run(filePath, func(t *testing.T) {
		assertIncluded(newShareTags(args, filePath), "mastodon", "linkedin")
	})

	filePath = "foo/bar/baz.share:-linkedin:-mastodon.txt"
	t.Run(filePath, func(t *testing.T) {
		assertExcluded(newShareTags(args, filePath), "mastodon", "linkedin")
	})

	filePath = "foo/bar/baz.share:-linkedin:mastodon.txt"
	t.Run(filePath, func(t *testing.T) {
		assertIncluded(newShareTags(args, filePath), "mastodon")
		assertExcluded(newShareTags(args, filePath), "linkedin")
	})
}
