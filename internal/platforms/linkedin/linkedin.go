package linkedin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/platforms/linkedin/oauth2"
	"codeberg.org/snonux/gos/internal/prompt"
)

var errUnauthorized = errors.New("unauthorized access, refresh or create token?")

const linkedInTimeout = 10 * time.Second

func Post(ctx context.Context, args config.Args, sizeLimit int, en entry.Entry) error {
	err := post(ctx, args, sizeLimit, en)
	if errors.Is(err, errUnauthorized) {
		colour.Infoln(err, "=> trying to refresh LinkedIn access token")
		args.Secrets.LinkedInAccessToken = "" // Reset the token
		return post(ctx, args, sizeLimit, en)
	}
	return err
}

func post(ctx context.Context, args config.Args, sizeLimit int, en entry.Entry) error {
	if args.DryRun {
		colour.Infoln("Not posting", en, "to LinkedIn as dry-run enabled")
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

	prev, err := NewPreview(newCtx, args, urls)
	if err != nil {
		return err
	}

	// TODO: Refactor this. Make it so that in a loop we can also check for the content with limit.
	// Maybe pass an interface en.ContentWithLimit and en.Path() to prompt.FileAction
	question := fmt.Sprintf("Do you want to post this message to Linkedin (%v)?", prev)
	if content, err = prompt.FileAction(question, content, en.Path); err != nil {
		return err
	}

	newCtx, cancel = context.WithTimeout(ctx, linkedInTimeout)
	defer cancel()
	return postMessageToLinkedInAPI(newCtx, personID, accessToken, content, prev)
}

// https://learn.microsoft.com/en-us/linkedin/marketing/community-management/shares/posts-api
func postMessageToLinkedInAPI(ctx context.Context, personID, accessToken, content string, prev preview) error {
	const linkedInPostsURL = "https://api.linkedin.com/rest/posts"

	personURN := fmt.Sprintf("urn:li:person:%s", personID)
	post := map[string]interface{}{
		"author":     personURN,
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

	article := map[string]interface{}{}
	if thumbnailPath, ok := prev.Thumbnail(); ok {
		thumbnailURN, err := postImageToLinkedInAPI(ctx, personURN, accessToken, thumbnailPath)
		if err != nil {
			return err
		}
		article["thumbnail"] = thumbnailURN
	}
	if title, url, ok := prev.TitleAndURL(); ok {
		article["title"] = title
		article["source"] = url
		post["content"] = map[string]interface{}{"article": article}
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

// https://learn.microsoft.com/en-us/linkedin/marketing/community-management/shares/images-api
func postImageToLinkedInAPI(ctx context.Context, personURN, accessToken, imagePath string) (string, error) {
	uploadURL, imageURN, err := initializeImageUpload(ctx, personURN, accessToken)
	if err != nil {
		return imageURN, err
	}
	return imageURN, performImageUpload(ctx, imagePath, uploadURL, accessToken)
}

func initializeImageUpload(ctx context.Context, personURN, accessToken string) (string, string, error) {
	const linkedInAPIURL = "https://api.linkedin.com/rest/images?action=initializeUpload"

	type InitializeUploadRequest struct {
		Owner string `json:"owner"`
	}
	requestBody, err := json.Marshal(map[string]interface{}{
		"initializeUploadRequest": InitializeUploadRequest{Owner: personURN},
	})

	if err != nil {
		return "", "", fmt.Errorf("error creating request body: %w", err)
	}

	// Initialize image upload
	req, err := http.NewRequestWithContext(ctx, "POST", linkedInAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("LinkedIn-Version", "202409")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	type InitializeUploadResponse struct {
		Value struct {
			UploadURL string `json:"uploadUrl"`
			Image     string `json:"image"`
		} `json:"value"`
	}
	var response InitializeUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", "", fmt.Errorf("error decoding response: %w", err)
	}

	return response.Value.UploadURL, response.Value.Image, nil
}

func performImageUpload(ctx context.Context, imagePath, uploadURL, accessToken string) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", uploadURL, bytes.NewBuffer(imageData))
	if err != nil {
		return fmt.Errorf("error creating upload request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending upload request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upload failed with status %s: %s", resp.Status, string(body))
	}
	return nil
}
