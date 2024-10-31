package mastodon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/prompt"
)

const mastodonTimeout = 10 * time.Second

func Post(ctx context.Context, args config.Args, sizeLimit int, ent entry.Entry) error {
	content, _, err := ent.ContentWithLimit(sizeLimit)
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
	if err := prompt.FileAction("Do you want to post this message to Mastodon?", content, ent.Path); err != nil {
		return err
	}

	newCtx, cancel := context.WithTimeout(ctx, mastodonTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(newCtx, "POST", args.Secrets.MastodonURL, bytes.NewBuffer(payloadBytes))
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d\n%s\n",
			resp.StatusCode, string(body))
	}
	return nil
}
