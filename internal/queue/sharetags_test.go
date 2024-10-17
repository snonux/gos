package queue

import (
	"slices"
	"testing"

	"codeberg.org/snonux/gos/internal/config"
)

func TestShareTagsPositive(t *testing.T) {
	t.Parallel()

	args := config.Args{Platforms: map[string]int{"mastodon": 100, "linkedin": 100}}
	testTable := map[string]shareTags{
		"./foo/bar.without.tags.txt.20240101-010101.queued": {
			includes: []string{"mastodon", "linkedin"},
		},
		"./foo/bar.share:linkeDin.txt.20240101-010101.queued": {
			includes: []string{"linkedin"},
		},
		"./foo/bar.share:-LinkedIn.txt.20240101-010101.queued": {
			excludes: []string{"linkedin"},
		},
		"./foo/bar.share:linkedin:mastOdon.txt.20240101-010101.queued": {
			includes: []string{"linkedin", "mastodon"},
		},
		"./foo/bar.share:linkediN:-mastodon:XCOM.txt.20240101-010101.queued": {
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

	args := config.Args{Platforms: map[string]int{"mastodon": 100, "linkedin": 100}}
	testTable := map[string]shareTags{
		"./foo/bar.without.tags.txt.20240101-010101.queued": {
			includes: []string{"linkedin"},
		},
		"./foo/bar.share:linkedIn.txt.20240101-010101.queued": {
			includes: []string{"mastodon"},
		},
		"./foo/bar.share:-liNkedin.txt.20240101-010101.queued": {
			includes: []string{"linkedin"},
		},
		"./foo/bar.share:linkedin:mastodon.txt.20240101-010101.queued": {
			includes: []string{"oups", "mastodon"},
		},
		"./foo/bar.share:linkedin:-MASTODON:xcom.txt.20240101-010101.queued": {
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
