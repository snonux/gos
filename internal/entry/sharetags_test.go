package entry

import (
	"slices"
	"strings"
	"testing"

	"codeberg.org/snonux/gos/internal/config"
)

func TestShareTagsPositive(t *testing.T) {
	args := config.Args{Platforms: map[string]int{
		"mastodon": 100,
		"linkedin": 100,
	}}
	testTable := map[string]shareTags{
		"./foo/bar.without.tags.txt.20240101-010101.queued": {
			includes: []string{"mastodon", "linkedin"},
		},
		"./foo/bar.share:linkeDin.txt.20240101-010101.queued": {
			includes: []string{"linkedin"},
		},
		"./foo/bar.share:-LinkedIn.txt.20240101-010101.queued": {
			includes: []string{"mastodon"},
			excludes: []string{"linkedin"},
		},
		"./foo/bar.share:linkedin:mastOdon.txt.20240101-010101.queued": {
			includes: []string{"linkedin", "mastodon"},
		},
		"./foo/bar.share:linkediN:-mastodon:XCOM.txt.20240101-010101.queued": {
			includes: []string{"linkedin", "xcom"},
			excludes: []string{"mastodon"},
		},
		"./foo/bar/ql-e7657e8a1ab573f84ad0dbc55199e937.share:-mastodon.txt.20241018-105524.queued": {
			includes: []string{"linkedin"},
			excludes: []string{"mastodon"},
		},
	}

	for filePath, expectedResult := range testTable {
		t.Run(filePath, func(t *testing.T) {
			shareTags, err := newShareTags(args, filePathTags(filePath))
			if err != nil {
				t.Error(err)
			}
			if !sameElements(shareTags.includes, expectedResult.includes) {
				t.Errorf("Expected includes to be %v but got %v with %s",
					expectedResult.includes, shareTags.includes, filePath)
			}
			if !sameElements(shareTags.excludes, expectedResult.excludes) {
				t.Errorf("Expected excludes to be %v but got %v with %s",
					expectedResult.excludes, shareTags.excludes, filePath)
			}
		})

	}
}
func TestShareTagsNegative(t *testing.T) {
	args := config.Args{Platforms: map[string]int{
		string("mastodon"): 100,
		string("linkedin"): 100,
	}}
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
			shareTags, err := newShareTags(args, filePathTags(filePath))
			if err != nil {
				t.Error(err)
			}
			if sameElements(shareTags.includes, unexpectedResult.includes) &&
				sameElements(shareTags.excludes, unexpectedResult.excludes) {
				t.Errorf("expected %v not to be the actual result with %s",
					unexpectedResult, filePath)
			}
		})

	}
}

// Can't use slices.Equal as order of elements may be different.
func sameElements(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, elem := range a {
		if !slices.Contains(b, elem) {
			return false
		}
	}
	return true
}

func filePathTags(filePath string) map[string]struct{} {
	tags := make(map[string]struct{})
	for _, tag := range strings.Split(filePath, ".") {
		tags[tag] = struct{}{}
	}
	return tags
}
