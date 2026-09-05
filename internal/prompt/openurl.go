package prompt

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/snonux/gos/internal/colour"
)

// Adapted from tasksamurai's urlRegex (https?://\S+); (?i) so HTTPS:// matches.
var urlRegex = regexp.MustCompile(`(?i)https?://\S+`)

// openURL starts browser with url and reaps the process in a goroutine
// (tasksamurai openURLCmd pattern) so the caller does not wait on the browser.
// Darwin uses `open -a`; Windows/Linux invoke the browser binary directly so
// query strings with '&' are not re-parsed by cmd.exe.
var openURL = func(browser, url string) error {
	colour.Infoln("Opening", url, "in", browser)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// Launch the browser binary directly (same as Linux). Going through
		// `cmd /C start` re-parses the command line and breaks query URLs that
		// contain '&'.
		cmd = exec.Command(browser, url)
	case "darwin":
		cmd = exec.Command("open", "-a", browser, url)
	default:
		cmd = exec.Command(browser, url)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("opening browser: %w", err)
	}
	go func() {
		_ = cmd.Wait()
	}()
	return nil
}

// yesNoFn is swapped in tests; production uses YesNo.
var yesNoFn = YesNo

// extractHTTPURLs returns http(s) URLs from text, trimming trailing prose
// punctuation that often abuts a link while preserving balanced path parens
// (e.g. Wikipedia "Title_(disambiguation)" URLs).
func extractHTTPURLs(text string) []string {
	raw := urlRegex.FindAllString(text, -1)
	if len(raw) == 0 {
		return nil
	}
	out := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for _, u := range raw {
		u = trimURLJunk(u)
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		out = append(out, u)
	}
	return out
}

// trimURLJunk removes trailing wrapper/prose punctuation. Closing ')' is only
// stripped while there are more ')' than '(' (prose like "see https://x.com).").
func trimURLJunk(u string) string {
	for u != "" {
		last := u[len(u)-1]
		switch last {
		case '.', ',', ';', ':', '!', '?', '"', '\'', '>', ']', '`', '*':
			u = u[:len(u)-1]
		case ')':
			if strings.Count(u, "(") >= strings.Count(u, ")") {
				return u
			}
			u = u[:len(u)-1]
		default:
			return u
		}
	}
	return u
}

// MaybeOpenURLs runs after the user has confirmed the post text. When the
// content contains http(s) URLs it offers to open them for verification, then
// requires an explicit proceed confirmation so posting does not race ahead of
// the browser. Declining to open skips verification and still allows posting;
// declining to proceed aborts with ErrAborted.
func MaybeOpenURLs(browser, content string) error {
	urls := extractHTTPURLs(content)
	if len(urls) == 0 {
		return nil
	}

	ok, err := yesNoFn(openURLsQuestion(urls))
	if err != nil {
		return err
	}
	if ok {
		for _, u := range urls {
			if err := openURL(browser, u); err != nil {
				colour.Errorln(err)
			}
		}
	}

	proceed, err := yesNoFn("Proceed with posting?")
	if err != nil {
		return err
	}
	if !proceed {
		return fmt.Errorf("%w after link check", ErrAborted)
	}
	return nil
}

func openURLsQuestion(urls []string) string {
	if len(urls) == 1 {
		return fmt.Sprintf("Open %s in the browser to verify the link still works?", urls[0])
	}
	listed := strings.Join(urls, "\n    ")
	return fmt.Sprintf("Open these %d URLs in the browser to verify the links still work?\n    %s\n ",
		len(urls), listed)
}
