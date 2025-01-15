package tags

import (
	"slices"
	"strings"
	"testing"
)

func TestInlineExtractTagsToFilePath(t *testing.T) {
	const filePath = "./gosdir/foo.golang.rox.txt"

	table := map[string]string{
		"foo,bar,baz blablablabla...":              "./gosdir/foo.golang.rox.foo.bar.baz.extracted.txt",
		"foo.bar.baz blablablabla...":              "./gosdir/foo.golang.rox.foo.bar.baz.extracted.txt",
		"foo.bar,baz blablablabla...":              "./gosdir/foo.golang.rox.foo.bar.baz.extracted.txt",
		"foo,bar.baz    blablablabla...":           "./gosdir/foo.golang.rox.foo.bar.baz.extracted.txt",
		"share:li,foo this    is the main content": "./gosdir/foo.golang.rox.share:linkedin.foo.extracted.txt",
		"share:li:ma this is the main content":     "./gosdir/foo.golang.rox.share:linkedin:mastodon.extracted.txt",
		"share:li:ma,now this is the main content": "./gosdir/foo.golang.rox.share:linkedin:mastodon.now.extracted.txt",
		"share,soon this will be shared soon":      "./gosdir/foo.golang.rox.share.soon.extracted.txt",
	}

	for content, expectedFilePath := range table {
		t.Run(content, func(t *testing.T) {
			newFilePath, _, err := inlineExtractTagsToFilePath(filePath, content)
			if err != nil {
				t.Error(err)
			}
			if newFilePath != expectedFilePath {
				t.Errorf("expected file path '%s' but got '%s'", expectedFilePath, newFilePath)
			}
		})
	}
}

func TestInlineExtractTagsFromContent(t *testing.T) {
	table := map[string][]string{
		"foo,bar,baz blablablabla...":                {"foo", "bar", "baz"},
		"foo.bar.baz blablablabla...":                {"foo", "bar", "baz"},
		"foo.bar,baz blablablabla...":                {"foo", "bar", "baz"},
		"foo,bar.baz    blablablabla...":             {"foo", "bar", "baz"},
		"share:li this    is the main content":       {"share:linkedin"},
		"share:li,foo this    is the main content":   {"share:linkedin", "foo"},
		"shar()e:li,foo this    is the main content": {"shar()e:li", "foo"},
		"share this post":                            {},
		"share,soon the main content here":           {"share", "soon"},
		`share,soon
			
			the main content here
			#foo`: {"share", "soon"},
	}

	for input, expectedTags := range table {
		t.Run(input, func(t *testing.T) {
			tags, contentWithoutTags, err := inlineExtractTagsFromContent(input)
			if err != nil {
				t.Error(err)
			}
			t.Log(expectedTags, tags)
			if len(tags) != len(expectedTags) {
				t.Errorf("expected %d inline tags (%v) but got %d (%v) for input '%v'",
					len(expectedTags), expectedTags, len(tags), tags, input)
			}
			for _, expectedTag := range expectedTags {
				if !slices.Contains(tags, expectedTag) {
					t.Errorf("expected '%s' to be an inline tag but got '%v'",
						expectedTag, tags)
				}
			}

			expectedMainContent := input
			parts := strings.Fields(input)
			if inlineTagRE.MatchString(parts[0]) {
				expectedMainContent = strings.TrimPrefix(expectedMainContent, parts[0])
			}

			if contentWithoutTags != strings.TrimSpace(expectedMainContent) {
				t.Errorf("expected the main content to be '%s' but got '%s'",
					expectedMainContent, strings.TrimSpace(contentWithoutTags))
			}
		})
	}
}
