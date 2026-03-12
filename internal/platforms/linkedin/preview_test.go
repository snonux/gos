package linkedin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"codeberg.org/snonux/gos/internal/config"
)

func TestPreviewExtract(t *testing.T) {
	var (
		expectedTitle    = "Baz"
		expectedImageURL = "https://free.beer:666/bar/foo.jpg"
		mockHTML         = strings.NewReader(`
<!DOCTYPE html>
<html>
<head>
    <title>Baz</title>
</head>
<body>
    <img src="./foo.jpg" alt="Foo">
</body>
</html>
`)
	)

	title, imageURL, err := extract("https://free.beer:666/bar/", mockHTML)
	if err != nil {
		t.Error(err)
	}
	if title != expectedTitle {
		t.Errorf("expected title '%s' but got '%s'", expectedTitle, title)
	}
	if imageURL != expectedImageURL {
		t.Errorf("expected imageURL '%s' but got '%s'", expectedImageURL, imageURL)
	}
}

func TestNewPreviewIgnoresForbiddenPage(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "blocked", http.StatusForbidden)
	}))
	defer srv.Close()

	prev, err := NewPreview(context.Background(), config.Args{
		CacheDir: t.TempDir(),
	}, []string{srv.URL})
	if err != nil {
		t.Fatalf("expected preview fetch failure to be non-fatal, got %v", err)
	}

	title, sourceURL, ok := prev.TitleAndURL()
	if !ok {
		t.Fatal("expected preview to keep URL fallback")
	}
	if title != srv.URL {
		t.Fatalf("expected fallback title %q, got %q", srv.URL, title)
	}
	if sourceURL != srv.URL {
		t.Fatalf("expected source URL %q, got %q", srv.URL, sourceURL)
	}
	if _, ok := prev.Thumbnail(); ok {
		t.Fatal("expected no thumbnail for forbidden preview page")
	}
}

func TestNewPreviewIgnoresForbiddenImage(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			_, _ = w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Blocked image</title>
</head>
<body>
    <img src="/blocked.jpg" alt="blocked">
</body>
</html>
`))
		case "/blocked.jpg":
			http.Error(w, "blocked", http.StatusForbidden)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	cacheDir := t.TempDir()
	prev, err := NewPreview(context.Background(), config.Args{
		CacheDir: cacheDir,
	}, []string{srv.URL})
	if err != nil {
		t.Fatalf("expected image download failure to be non-fatal, got %v", err)
	}

	title, sourceURL, ok := prev.TitleAndURL()
	if !ok {
		t.Fatal("expected preview title and URL to be preserved")
	}
	if title != "Blocked image" {
		t.Fatalf("expected title %q, got %q", "Blocked image", title)
	}
	if sourceURL != srv.URL {
		t.Fatalf("expected source URL %q, got %q", srv.URL, sourceURL)
	}
	if thumbnailPath, ok := prev.Thumbnail(); ok {
		t.Fatalf("expected no thumbnail after blocked download, got %q", thumbnailPath)
	}

	matches, err := filepath.Glob(filepath.Join(cacheDir, "*"))
	if err != nil {
		t.Fatalf("glob failed: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected cache dir to stay empty, got %v", matches)
	}
}
