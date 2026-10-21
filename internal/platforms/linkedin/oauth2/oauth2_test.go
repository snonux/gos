package oauth2

import (
	"net/http"
	"net/http/httptest"
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

// TestCallbackRejectsWrongState checks that a callback which does not carry
// the state of the flow in progress is refused before any code exchange.
func TestCallbackRejectsWrongState(t *testing.T) {
	oauthState = "expected-state"
	errCh = make(chan error, 1)
	defer func() { oauthState = "" }()

	for _, target := range []string{"/callback?code=x", "/callback?code=x&state=forged"} {
		rec := httptest.NewRecorder()
		oauthCallbackHandler(rec, httptest.NewRequest(http.MethodGet, target, nil))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("%s: got status %d, want %d", target, rec.Code, http.StatusBadRequest)
		}
	}
	if oauthState != "expected-state" {
		t.Fatal("a rejected callback must not consume the state")
	}
}

func TestNewOAuthStateIsRandom(t *testing.T) {
	a, err := newOAuthState()
	if err != nil {
		t.Fatal(err)
	}
	b, err := newOAuthState()
	if err != nil {
		t.Fatal(err)
	}
	if a == b || len(a) != 64 {
		t.Fatalf("unexpected states %q and %q", a, b)
	}
}
