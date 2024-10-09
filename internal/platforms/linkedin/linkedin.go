package linkedin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/platforms/linkedin/oauth2"
)

// TODO: Also implemebt a Text Platform output, which then laster can be
// processed by Gemtexter as a page
func Post(ctx context.Context, args config.Args, ent entry.Entry) error {
	content, err := ent.Content()
	if err != nil {
		return nil
	}

	personID, accessToken, err := oauth2.LinkedInCreds(args)
	if err != err {
		return err
	}

	return post(personID, accessToken, content)
}

func post(personID, accessToken, message string) error {
	const url = "https://api.linkedin.com/v2/posts"

	post := map[string]interface{}{
		"author":     fmt.Sprintf("urn:li:person:%s", personID),
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

	req.Header.Add("Authorization", "Bearer "+accessToken)
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
