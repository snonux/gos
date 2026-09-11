package oauth2

import (
	"net/http"
	"testing"

	"github.com/snonux/gos/internal/config"
)

func TestTryCredsEmpty(t *testing.T) {
	t.Parallel()
	if _, _, ok := tryCreds(config.Config{}); ok {
		t.Fatal("expected empty config credentials to be rejected")
	}
	if _, _, ok := tryCreds(config.Config{LinkedInAccessToken: "tok"}); ok {
		t.Fatal("expected missing person ID to be rejected")
	}
	if _, _, ok := tryCreds(config.Config{LinkedInPersonID: "id"}); ok {
		t.Fatal("expected missing access token to be rejected")
	}
}

// Go 1.22+ panics when the same pattern is registered twice on a ServeMux.
// LinkedInCreds must use a fresh mux per OAuth attempt; this guards the
// DefaultServeMux footgun that caused the delete-retry panic.
func TestDefaultServeMuxRejectsDuplicatePatterns(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on duplicate pattern registration")
		}
	}()
	mux.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})
}
