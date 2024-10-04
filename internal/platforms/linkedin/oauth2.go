package linkedin

import (
	"bytes"
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

var oauthConfig *oauth2.Config

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

func postMessage(token *oauth2.Token, linkedInID, message string) error {
	const url = "https://api.linkedin.com/v2/posts"

	post := map[string]interface{}{
		"author":     fmt.Sprintf("urn:li:person:%s", linkedInID),
		"commentary": message,
		"visibility": "PUBLIC",
		"distribution": map[string]interface{}{
			"feedDistribution":               "MAIN_FEED",
			"targetEntities":                 []string{},
			"thirdPartyDistributionChannels": []string{},
		},
		"lifecycleState":            "PUBLISHED",
		"isReshareDisabledByAuthor": false,
	}

	payload, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("Error encoding JSON:%w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("Error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-RestLi-Protocol-Version", "2.0.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Failed to post to LinkedIn. Status: %s\n%s\n\n", resp.Status, body)
	}
	return nil
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

// TODO: Check for how logn the access token is valid for
// TODO: Fetch the access token and user ID and store it i na file in .config/gos/...
// TODO: Refresh access token when it is about to expire or expired
// TODO: Separate posting of the message and fetching of the userID and access token
func oauth(args config.Args) error {
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
	return http.ListenAndServe(":8080", nil)
}
