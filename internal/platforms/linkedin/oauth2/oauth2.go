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

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/config"
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
	defer resp.Body.Close()

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

func LinkedInCreds(ctx context.Context, args config.Args) (string, string, error) {
	secrets := args.Secrets
	if secrets.LinkedInAccessToken != "" && secrets.LinkedInPersonID != "" {
		return secrets.LinkedInPersonID, secrets.LinkedInAccessToken, nil
	}

	oauthConfig = &oauth2.Config{
		ClientID:     secrets.LinkedInClientID,
		ClientSecret: secrets.LinkedInSecret,
		RedirectURL:  secrets.LinkedInRedirectURL,
		Scopes:       []string{"openid", "profile", "w_member_social"},
		Endpoint:     linkedin.Endpoint,
	}
	errCh = make(chan error)
	globalCtx = ctx

	http.HandleFunc("/", oauthIndexHandler)
	http.HandleFunc("/callback", oauthCallbackHandler)
	http.HandleFunc("/up", upHandler)

	colour.Infoln("Listening on http://localhost:8080 for LinkedIn OAuth2")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			errCh <- err
		}
	}()

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

	secrets.LinkedInAccessToken = oauthAccessToken
	secrets.LinkedInPersonID = oauthPersonID
	return oauthPersonID, oauthAccessToken, secrets.WriteToDisk(args.SecretsConfigPath)
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
			colour.Infof("URL is not reachable: %v", err)
			fmt.Print("\n")
		} else {
			colour.Infof("URL is reachable: %s - Status Code: %d", url, resp.StatusCode)
			fmt.Print("\n")
			resp.Body.Close()
			return nil
		}
	}
	return fmt.Errorf("%s not reachable after %d tries", url, counter)
}
