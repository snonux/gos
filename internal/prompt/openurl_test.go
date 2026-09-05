package prompt

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestExtractHTTPURLs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want []string
	}{
		{name: "empty", in: "no links here", want: nil},
		{name: "single https", in: "see https://example.com please", want: []string{"https://example.com"}},
		{name: "http and https", in: "a http://a.test b https://b.test/c", want: []string{"http://a.test", "https://b.test/c"}},
		{name: "query string", in: "https://github.com/bjarneo/ku?ref=terminaltrove", want: []string{"https://github.com/bjarneo/ku?ref=terminaltrove"}},
		{name: "trailing prose paren", in: "link https://foo.zone).", want: []string{"https://foo.zone"}},
		{name: "balanced path paren", in: "https://en.wikipedia.org/wiki/Go_(programming_language)", want: []string{"https://en.wikipedia.org/wiki/Go_(programming_language)"}},
		{name: "balanced path paren with prose", in: "see https://en.wikipedia.org/wiki/Go_(programming_language).", want: []string{"https://en.wikipedia.org/wiki/Go_(programming_language)"}},
		{name: "angle brackets", in: "see <https://example.com>", want: []string{"https://example.com"}},
		{name: "quoted", in: `"https://example.com"`, want: []string{"https://example.com"}},
		{name: "backticks", in: "see `https://example.com`", want: []string{"https://example.com"}},
		{name: "uppercase scheme", in: "HTTPS://Example.COM/Path", want: []string{"HTTPS://Example.COM/Path"}},
		{name: "ftp ignored", in: "ftp://files.example/x", want: nil},
		{name: "dedupe", in: "https://a.test and again https://a.test", want: []string{"https://a.test"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := extractHTTPURLs(tc.in)
			if !slices.Equal(got, tc.want) {
				t.Fatalf("extractHTTPURLs(%q) = %#v, want %#v", tc.in, got, tc.want)
			}
		})
	}
}

func TestOpenURLsQuestion(t *testing.T) {
	t.Parallel()
	one := openURLsQuestion([]string{"https://a.test"})
	if !strings.Contains(one, "https://a.test") {
		t.Fatalf("single question missing url: %q", one)
	}
	many := openURLsQuestion([]string{"https://a.test", "https://b.test"})
	if !strings.Contains(many, "https://a.test") || !strings.Contains(many, "https://b.test") {
		t.Fatalf("multi question missing urls: %q", many)
	}
}

func TestMaybeOpenURLsNoURLsSkipsPrompt(t *testing.T) {
	origYes := yesNoFn
	origOpen := openURL
	t.Cleanup(func() {
		yesNoFn = origYes
		openURL = origOpen
	})

	yesNoFn = func(string) (bool, error) {
		t.Fatal("YesNo should not be called when there are no URLs")
		return false, nil
	}
	openURL = func(string, string) error {
		t.Fatal("openURL should not be called when there are no URLs")
		return nil
	}

	if err := MaybeOpenURLs("firefox", "plain text only"); err != nil {
		t.Fatalf("MaybeOpenURLs: %v", err)
	}
}

func TestMaybeOpenURLsSkipOpenStillProceeds(t *testing.T) {
	origYes := yesNoFn
	origOpen := openURL
	t.Cleanup(func() {
		yesNoFn = origYes
		openURL = origOpen
	})

	calls := 0
	yesNoFn = func(string) (bool, error) {
		calls++
		if calls == 1 {
			return false, nil // decline open
		}
		return true, nil // proceed
	}
	openURL = func(string, string) error {
		t.Fatal("openURL should not be called when user declines open")
		return nil
	}

	if err := MaybeOpenURLs("firefox", "see https://example.com"); err != nil {
		t.Fatalf("MaybeOpenURLs: %v", err)
	}
	if calls != 2 {
		t.Fatalf("YesNo calls = %d, want 2", calls)
	}
}

func TestMaybeOpenURLsDeclineProceedAborts(t *testing.T) {
	origYes := yesNoFn
	origOpen := openURL
	t.Cleanup(func() {
		yesNoFn = origYes
		openURL = origOpen
	})

	yesNoFn = func(string) (bool, error) { return false, nil }
	openURL = func(string, string) error { return nil }

	err := MaybeOpenURLs("firefox", "https://example.com")
	if !errors.Is(err, ErrAborted) {
		t.Fatalf("err = %v, want ErrAborted", err)
	}
}

func TestMaybeOpenURLsYesOpensAll(t *testing.T) {
	origYes := yesNoFn
	origOpen := openURL
	t.Cleanup(func() {
		yesNoFn = origYes
		openURL = origOpen
	})

	yesNoFn = func(string) (bool, error) { return true, nil }
	var opened []string
	openURL = func(browser, url string) error {
		if browser != "firefox" {
			t.Fatalf("browser = %q, want firefox", browser)
		}
		opened = append(opened, url)
		return nil
	}

	content := "a https://one.example b https://two.example/path"
	if err := MaybeOpenURLs("firefox", content); err != nil {
		t.Fatalf("MaybeOpenURLs: %v", err)
	}
	want := []string{"https://one.example", "https://two.example/path"}
	if !slices.Equal(opened, want) {
		t.Fatalf("opened = %#v, want %#v", opened, want)
	}
}

func TestMaybeOpenURLsOpenFailureDoesNotAbort(t *testing.T) {
	origYes := yesNoFn
	origOpen := openURL
	t.Cleanup(func() {
		yesNoFn = origYes
		openURL = origOpen
	})

	yesNoFn = func(string) (bool, error) { return true, nil }
	openURL = func(string, string) error {
		return errors.New("opening browser: boom")
	}

	if err := MaybeOpenURLs("firefox", "https://example.com"); err != nil {
		t.Fatalf("MaybeOpenURLs returned error %v, want nil so posting can continue after proceed", err)
	}
}

func TestMaybeOpenURLsYesNoErrorPropagates(t *testing.T) {
	origYes := yesNoFn
	t.Cleanup(func() { yesNoFn = origYes })

	yesNoFn = func(string) (bool, error) {
		return false, errors.New("stdin closed")
	}

	err := MaybeOpenURLs("firefox", "https://example.com")
	if err == nil || !strings.Contains(err.Error(), "stdin closed") {
		t.Fatalf("err = %v, want stdin closed", err)
	}
}

func TestOpenURLStartsBrowser(t *testing.T) {
	tmp := t.TempDir()
	logPath := filepath.Join(tmp, "opened.txt")
	browserPath := filepath.Join(tmp, "browser")
	script := "#!/bin/sh\nprintf '%s\\n' \"$1\" > " + logPath + "\n"
	if err := os.WriteFile(browserPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	const url = "https://example.com/check"
	if err := openURL(browserPath, url); err != nil {
		t.Fatalf("openURL: %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for {
		data, err := os.ReadFile(logPath)
		if err == nil {
			if got := strings.TrimSpace(string(data)); got != url {
				t.Fatalf("browser got %q, want %q", got, url)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("browser log was not written: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestOpenURLStartError(t *testing.T) {
	err := openURL(filepath.Join(t.TempDir(), "missing-browser"), "https://example.com")
	if err == nil || !strings.Contains(err.Error(), "opening browser:") {
		t.Fatalf("err = %v, want opening browser error", err)
	}
}
