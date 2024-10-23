package linkedin

import (
	"fmt"
	"slices"
	"testing"
)

func TestLinkedInEscapes(t *testing.T) {
	var (
		input    = `This is a test message with special characters: " {} @ [] () <> # \ * _ ~ |`
		expected = `This is a test message with special characters: \" \{\} @ \[\] \(\) \<\> # \\ \* \_ \~ \|`
	)
	if escaped := escapeLinkedInText(input); escaped != expected {
		t.Errorf("expected '%s' but got '%s'", expected, escaped)
	}
}

func TestLinkedInTwoURLsExtract(t *testing.T) {
	text := `Hello world https://foo.zone
	Hello universe http://world.universe test 123`

	urls := extractURLs(text)
	if len(urls) != 2 {
		t.Errorf("expected 2 URLs, but got %d", len(urls))
	}

	if !slices.Contains(urls, "https://foo.zone") {
		t.Errorf("expected 'https://foo.zone' in the URL list, but got %v", urls)
	}
	if !slices.Contains(urls, "http://world.universe") {
		t.Errorf("expected 'http://world.universe' in the URL list, but got %v", urls)
	}
}

// TODO: Use Fuzzing here!
func TestLinkedInURLExtract(t *testing.T) {
	urls := []string{
		"http://foo.zone",
		"http://foo.zone/",
		"http://foo.zone?foo=bar",
		"http://foo.zone/?foo=bar",
		"http://foo.zone/?foo=bar",
		"http://foo.zone/hurs?foo=bar",
		"http://foo.zone?foo=bar&baz=bay",
	}

	for _, url := range urls {
		text := fmt.Sprintf("Hello world %s Hello World", url)
		found := extractURLs(text)
		if len(found) != 1 {
			t.Errorf("expected 1 URL, but got %d for text '%s'", len(found), text)
		}
		if found[0] != url {
			t.Errorf("expected URL '%s', but got '%s' for text '%s'", url, found[0], text)
		}
	}
}
