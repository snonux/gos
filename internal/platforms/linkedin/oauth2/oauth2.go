package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/snonux/gos/internal/colour"
	"github.com/snonux/gos/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

var (
	oauthConfig      *oauth2.Config
	oauthAccessToken string
	oauthPersonID    string
	errCh            chan error
	globalCtx        context.Context
)

func getOauthPersonID(token *oauth2.Token) (string, error) {
	const url = "https://api.linkedin.com/v2/userinfo"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating request:%w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("X-RestLi-Protocol-Version", "2.0.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making the request:%w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log the error but don't fail the operation since we've already read the data
			colour.Errorln("Error closing response body:", err)
		}
	}()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to retrieve user profile. Status: %s\n%s\n", resp.Status, string(body))
	}

	type User struct {
		Sub string `json:"sub"`
	}
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return "", fmt.Errorf("Error unmarshalling JSON: %w", err)
	}

	return user.Sub, nil
}

func upHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("I am up!\n"))
}

func oauthIndexHandler(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func oauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	defer close(errCh)
	code := r.URL.Query().Get("code")

	colour.Infoln("Exchanging OAuth2 token")
	token, err := oauthConfig.Exchange(globalCtx, code)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		errCh <- err
		return
	}
	oauthAccessToken = token.AccessToken
	_, _ = w.Write([]byte("Successfully fetched LinkedIn access token\n"))

	if oauthPersonID, err = getOauthPersonID(token); err != nil {
		_, _ = w.Write([]byte(err.Error()))
		errCh <- err
		return
	}
	_, _ = w.Write([]byte("Successfully fetched LinkedIn person ID\n"))
}

// tryCreds returns LinkedIn credentials when the access token can call the userinfo API.
func tryCreds(conf config.Config) (personID, accessToken string, ok bool) {
	if conf.LinkedInAccessToken == "" || conf.LinkedInPersonID == "" {
		return "", "", false
	}
	token := &oauth2.Token{AccessToken: conf.LinkedInAccessToken}
	if _, err := getOauthPersonID(token); err != nil {
		return "", "", false
	}
	return conf.LinkedInPersonID, conf.LinkedInAccessToken, true
}

func LinkedInCreds(ctx context.Context, args config.Args) (string, string, error) {
	if personID, accessToken, ok := tryCreds(args.Config); ok {
		return personID, accessToken, nil
	}

	// args.Config is a value copy. An earlier OAuth in this process may have
	// written fresh credentials to disk without updating the caller's Args.
	if disk, err := config.New(args.ConfigPath, false); err == nil {
		if personID, accessToken, ok := tryCreds(disk); ok {
			return personID, accessToken, nil
		}
	}

	conf := args.Config
	oauthConfig = &oauth2.Config{
		ClientID:     conf.LinkedInClientID,
		ClientSecret: conf.LinkedInSecret,
		RedirectURL:  conf.LinkedInRedirectURL,
		Scopes:       []string{"openid", "profile", "w_member_social"},
		Endpoint:     linkedin.Endpoint,
	}
	errCh = make(chan error)
	globalCtx = ctx

	// Use a dedicated mux/server so a second OAuth attempt in the same process
	// does not panic on re-registering "/" with http.DefaultServeMux (Go 1.22+).
	mux := http.NewServeMux()
	mux.HandleFunc("/", oauthIndexHandler)
	mux.HandleFunc("/callback", oauthCallbackHandler)
	mux.HandleFunc("/up", upHandler)

	server := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			colour.Errorln("Error shutting down LinkedIn OAuth server:", err)
		}
	}()

	colour.Infoln("Listening on http://localhost:8080 for LinkedIn OAuth2")

	if err := WaitUntilURLIsReachable("http://localhost:8080/up"); err != nil {
		return "", "", err
	}

	if err := openURLInFirefox(args.OAuth2Browser, "http://localhost:8080"); err != nil {
		return "", "", err
	}

	var errs error
	for err := range errCh {
		errs = errors.Join(errs, err)
	}
	if errs != nil {
		return "", "", errs
	}

	conf.LinkedInAccessToken = oauthAccessToken
	conf.LinkedInPersonID = oauthPersonID
	return oauthPersonID, oauthAccessToken, conf.WriteToDisk(args.ConfigPath)
}

func openURLInFirefox(browser, url string) error {
	colour.Infoln("Opening", url, "in", browser)
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/C", "start", browser, url)
		return cmd.Start()
	case "darwin":
		cmd := exec.Command("open", "-a", browser, url)
		return cmd.Start()
	default:
		// Linux and other Unix like (e.g. *BSDs)
		cmd := exec.Command(browser, url)
		return cmd.Start()
	}
}

func WaitUntilURLIsReachable(url string) error {
	var counter int
	for counter < 10 {
		counter++
		time.Sleep(1 * time.Second)
		resp, err := http.Get(url)

		if err != nil {
			colour.Infofln("URL is not reachable: %v", err)
		} else {
			colour.Infofln("URL is reachable: %s - Status Code: %d", url, resp.StatusCode)
			if err := resp.Body.Close(); err != nil {
				// Log the error but don't fail the operation since we've already read the data
				colour.Errorln("Error closing response body:", err)
			}
			return nil
		}
	}
	return fmt.Errorf("%s not reachable after %d tries", url, counter)
}
