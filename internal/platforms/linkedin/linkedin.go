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
	"time"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/platforms/linkedin/oauth2"
	"codeberg.org/snonux/gos/internal/prompt"
)

var errUnauthorized = errors.New("unauthorized access, refresh or create token?")

const (
	linkedInPostsURL = "https://api.linkedin.com/rest/posts"
	linkedInTimeout  = 10 * time.Second
)

func Post(ctx context.Context, args config.Args, sizeLimit int, en entry.Entry) error {
	err := post(ctx, args, sizeLimit, en)
	if errors.Is(err, errUnauthorized) {
		log.Println(err, "=> trying to refresh LinkedIn access token")
		args.Secrets.LinkedInAccessToken = "" // Reset the token
		return post(ctx, args, sizeLimit, en)
	}
	return err
}

func post(ctx context.Context, args config.Args, sizeLimit int, en entry.Entry) error {
	if args.DryRun {
		log.Println("Not posting", en, "to LinkedIn as dry-run enabled")
		return nil
	}

	newCtx, cancel := context.WithTimeout(ctx, linkedInTimeout)
	defer cancel()
	personID, accessToken, err := oauth2.LinkedInCreds(newCtx, args)
	if err != nil {
		return err
	}
	content, urls, err := en.ContentWithLimit(sizeLimit)
	if err != nil {
		return err
	}

	newCtx, cancel = context.WithTimeout(ctx, linkedInTimeout)
	defer cancel()
	prev, err := NewPreview(newCtx, urls)
	if err != nil {
		return err
	}

	var filePath string
	if prev.imageURL != "" {
		if filePath, err = prev.DownloadImage(args.CacheDir); err != nil {
			return err
		}
		log.Println("Downloaded preview image to ", filePath)
	}

	question := fmt.Sprintf("Do you want to post this message to Linkedin (%v)?", prev)
	if err := prompt.FileAction(question, content, en.Path); err != nil {
		return err
	}

	newCtx, cancel = context.WithTimeout(ctx, linkedInTimeout)
	defer cancel()
	return callLinkedInAPI(newCtx, personID, accessToken, content, prev)
}

// TODO: Also post preview images
func callLinkedInAPI(ctx context.Context, personID, accessToken, content string, prev preview) error {
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

	if !prev.Empty() {
		post["content"] = map[string]interface{}{
			"article": map[string]interface{}{
				"title":  prev.title,
				"source": prev.url,
			},
		}
	}

	payload, err := json.Marshal(post)
	fmt.Println(string(payload))
	if err != nil {
		return fmt.Errorf("Error encoding JSON:%w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", linkedInPostsURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("Error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-RestLi-Protocol-Version", "2.0.0")
	req.Header.Add("LinkedIn-Version", "202409")

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
