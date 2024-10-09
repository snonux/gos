package mastodon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/prompt"
)

func Post(ctx context.Context, args config.Args, ent entry.Entry) error {
	content, err := ent.Content()
	if err != nil {
		return err
	}
	payload := map[string]string{"status": content}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	if args.DryRun {
		log.Println("Not posting", ent, "to Mastodon as dry-run enabled")
		return nil
	}
	if !prompt.Yes(fmt.Sprintf("%s\nDo you want to post this message to Mastodon?", string(payloadBytes))) {
		return prompt.ErrAborted
	}

	req, err := http.NewRequest("POST", args.Secrets.MastodonURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+args.Secrets.MastodonAccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
