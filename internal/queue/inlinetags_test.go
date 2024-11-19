package queue

import (
	"slices"
	"strings"
	"testing"
)

func TestExtractInlineTagsToFilePath(t *testing.T) {
	const filePath = "./gosdir/foo.golang.rox.txt"

	table := map[string]string{
		"foo,bar,baz blablablabla...":              "./gosdir/foo.golang.rox.foo.bar.baz.extracted.txt",
		"foo.bar.baz blablablabla...":              "./gosdir/foo.golang.rox.foo.bar.baz.extracted.txt",
		"foo.bar,baz blablablabla...":              "./gosdir/foo.golang.rox.foo.bar.baz.extracted.txt",
		"foo,bar.baz    blablablabla...":           "./gosdir/foo.golang.rox.foo.bar.baz.extracted.txt",
		"share:li,foo this    is the main content": "./gosdir/foo.golang.rox.share:li.foo.extracted.txt",
	}

	for content, expectedFilePath := range table {
		t.Run(content, func(t *testing.T) {
			newFilePath, _, err := extractInlineTagsToFilePath(filePath, content)
			if err != nil {
				t.Error(err)
			}
			if newFilePath != expectedFilePath {
				t.Errorf("expected file path '%s' but got '%s'", expectedFilePath, newFilePath)
			}
		})
	}
}

func TestExtractInlineTagsFromContent(t *testing.T) {
	table := map[string][]string{
		"foo,bar,baz blablablabla...":                {"foo", "bar", "baz"},
		"foo.bar.baz blablablabla...":                {"foo", "bar", "baz"},
		"foo.bar,baz blablablabla...":                {"foo", "bar", "baz"},
		"foo,bar.baz    blablablabla...":             {"foo", "bar", "baz"},
		"share:li,foo this    is the main content":   {"share:li", "foo"},
		"shar()e:li,foo this    is the main content": {},
	}

	for input, expectedTags := range table {
		t.Run(input, func(t *testing.T) {
			tags, contentWithoutTags, err := extractInlineTagsFromContent(input)
			if err != nil {
				t.Error(err)
			}
			if len(tags) != len(expectedTags) {
				t.Errorf("expected %d inline tags (%v) but got %d (%v)",
					len(expectedTags), expectedTags, len(tags), tags)
			}
			for _, expectedTag := range expectedTags {
				if !slices.Contains(tags, expectedTag) {
					t.Errorf("expected '%s' to be an inline tag but got '%v'",
						expectedTag, tags)
				}
			}

			expectedMainContent := input
			parts := strings.Split(input, " ")
			if inlineTagRE.MatchString(parts[0]) {
				expectedMainContent = strings.Join(parts[1:], " ")
			}

			if contentWithoutTags != expectedMainContent {
				t.Errorf("expected the main content to be '%s' but got '%s'",
					expectedMainContent, contentWithoutTags)
			}
		})
	}
}
