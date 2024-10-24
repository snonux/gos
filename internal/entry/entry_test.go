package entry

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	"codeberg.org/snonux/gos/internal/timestamp"
)

func TestEntry(t *testing.T) {
	states := []State{Queued, Posted}
	stamps := []string{"20240928-111835", "20241028-120135"}

	for _, state := range states {
		for _, stamp := range stamps {
			queuedPath := fmt.Sprintf("gosdir/db/platforms/linkedin/helloworld.txt.%s.%s", stamp, state)

			ent, err := New(queuedPath)
			if err != nil {
				t.Error(err)
			}
			if ent.Path != queuedPath {
				t.Errorf("expected path %s but got %s", queuedPath, ent.Path)
			}
			if ent.State != state {
				t.Errorf("expected state %s but got %s", state, ent.State)
			}

			expectedTime, err := timestamp.Parse(stamp)
			if err != nil {
				t.Error(err)
			}
			if ent.Time != expectedTime {
				t.Errorf("expected time to be %v but got %v", expectedTime, ent.Time)
			}
		}
	}
}

func TestExtractTwoURLs(t *testing.T) {
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

func TestExtractURLs(t *testing.T) {
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

func FuzzExtractURLs(f *testing.F) {
	f.Add("/path?myjfa=lwsr4imj&dgqeg=m3uwwsak")
	f.Add("/?amfbm=bwzqu46m&xheuh=nv588d98")
	f.Add("?tuupm=reng2p1y&cbjot=0g5qvpty")
	f.Add("/path?qmcok=f%20w4tfp7g&awsnq=sjizuore&owdix=8s2dmqsv")
	f.Add("?zwilf=868o24x1&fiwmp=1d5aqbvo&irhhr=xar7qbq7&eetpy=scmi9s8i")
	f.Add("/path?mwhbm=psinstn6&nsjic=pfu0wnk9&lbmrz=5bixkhdt")
	f.Add("/path?owbwo=67mkjiz2")
	f.Add("/path?ohvxi=esy5qvml&zlvzt=2yi4q4ef&cnich=sgc8sahs")
	f.Add("/path?codsl=fpwfto6j")
	f.Add("tvdus=fhlhlh1y")
	f.Add("/foo.txt")
	noWhitespace := regexp.MustCompile(`\s+`)

	f.Fuzz(func(t *testing.T, urlPath string) {
		urlPath = noWhitespace.ReplaceAllString(strings.TrimSpace(urlPath), "%20")
		baseURLs := []string{"https://foo.zone", "http://foo.zone", "ftp://foo.zone"}
		for _, baseURL := range baseURLs {
			fullURL := fmt.Sprintf("%s%s", baseURL, urlPath)
			text := fmt.Sprintf("Hello world %s Hello World", fullURL)
			found := extractURLs(text)
			if len(found) != 1 {
				t.Errorf("expected 1 URL '%s', but got %d for text '%s'",
					fullURL, len(found), text)
			}
			if found[0] != fullURL {
				t.Errorf("expected URL '%s', but got '%s' for text '%s'",
					fullURL, found[0], text)
			}
		}
	})
}
