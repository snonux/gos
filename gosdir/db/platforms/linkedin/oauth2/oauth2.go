package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"codeberg.org/snonux/gos/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

var (
	oauthConfig      *oauth2.Config
	oauthPersonId    string
	oauthAccessToken string
)

func getLinkedInID(token *oauth2.Token) (string, error) {
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

func oauthIndexHandler(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func oauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	linkedInID, err := getLinkedInID(token)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write([]byte("Successfully fetched the LinkedInID\n"))

	if err := postMessage(token, linkedInID, "test"); err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write([]byte("Successfully posted a message to LinkedIn!\n"))
}

// TODO: Fetch the access token and user ID and store it i na file in .config/gos/...
// TODO: Separate posting of the message and fetching of the userID and access token
func AccessToken(args config.Args) (config.Secrets, error) {
	if args.Secrets.LinkedInAccessToken != "" && args.Secrets.LinkedInPersonID != "" {
		// TODO: Check, whether the access token is still valid. If not, get a new one.
		return args.Secrets, nil
	}

	oauthConfig = &oauth2.Config{
		ClientID:     args.Secrets.LinkedInClientID,
		ClientSecret: args.Secrets.LinkedInSecret,
		RedirectURL:  args.Secrets.LinkedInRedirectURL,
		Scopes:       []string{"openid", "profile", "w_member_social"},
		Endpoint:     linkedin.Endpoint,
	}

	http.HandleFunc("/", oauthIndexHandler)
	http.HandleFunc("/callback", oauthCallbackHandler)

	log.Println("Listening on http://localhost:8080 for LinkedIn oauth2")
	err := http.ListenAndServe(":8080", nil)

	args.Secrets.MastodonAccessToken = oauthAccessToken
	args.Secrets.LinkedInPersonID = oauthPersonId

	return args.Secrets, err
}
