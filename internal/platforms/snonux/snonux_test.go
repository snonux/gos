package snonux

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/snonux/gos/internal/entry"
)

func TestInboxMarkdownPath(t *testing.T) {
	en, err := entry.New("/queue/foo.bar.share:snonux.txt.20250101-120000.queued")
	if err != nil {
		t.Fatal(err)
	}
	path, err := inboxMarkdownPath("/home/u/.gosdir", en)
	if err != nil {
		t.Fatal(err)
	}
	wantDir := filepath.FromSlash("/home/u/.gosdir/snonux/inbox")
	if dir := filepath.Dir(path); dir != wantDir {
		t.Errorf("dir = %q want %q", dir, wantDir)
	}
	base := filepath.Base(path)
	if !strings.HasSuffix(base, ".md") {
		t.Errorf("expected .md suffix, got %q", base)
	}
	if !strings.Contains(base, "foo") {
		t.Errorf("expected entry name in filename, got %q", base)
	}
}
