package linkedin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/platforms/linkedin/oauth2"
	"codeberg.org/snonux/gos/internal/prompt"
)

var errUnauthorized = errors.New("unauthorized access, refresh or create token?")

// TODO: Why are no previews of links shown then posted?
func Post(ctx context.Context, args config.Args, sizeLimit int, ent entry.Entry) error {
	err := post(ctx, args, sizeLimit, ent)
	if errors.Is(err, errUnauthorized) {
		log.Println(err, "=> trying to refresh LinkedIn access token")
		args.Secrets.LinkedInAccessToken = "" // Reset the token
		return post(ctx, args, sizeLimit, ent)
	}
	return err
}

func post(ctx context.Context, args config.Args, sizeLimit int, ent entry.Entry) error {
	if args.DryRun {
		log.Println("Not posting", ent, "to LinkedIn as dry-run enabled")
		return nil
	}
	personID, accessToken, err := oauth2.LinkedInCreds(ctx, args)
	if err != err {
		return err
	}
	content, err := ent.ContentWithLimit(sizeLimit)
	if err != nil {
		return err
	}
	if err := prompt.DoYouWantThis("Do you want to post this message to Linkedin?", content); err != nil {
		if errors.Is(err, prompt.ErrEditContent) {
			if err := ent.Edit(); err != nil {
				return err
			}
			return post(ctx, args, sizeLimit, ent)
		}
		return err
	}
	return callLinkedInAPI(ctx, personID, accessToken, content)
}

func callLinkedInAPI(ctx context.Context, personID, accessToken, content string) error {
	const url = "https://api.linkedin.com/v2/posts"

	post := map[string]interface{}{
		"author":     fmt.Sprintf("urn:li:person:%s", personID),
		"commentary": escapeLinkedInText(content),
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
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("failed to post to LinkedIn. Status: %s\n%s\n",
			resp.Status, string(body))
		if resp.StatusCode == http.StatusUnauthorized {
			err = errors.Join(err, errUnauthorized)
		}
	}
	return err
}
