package mastodon

import (
	"context"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
)

func Post(ctx context.Context, args config.Args, ent entry.Entry) error {
	// 	payload := map[string]string{"status": status}
	// 	payloadBytes, err := json.Marshal(payload)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to marshal payload: %w", err)
	// 	}

	// 	// Create the HTTP request
	// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	// 	if err != nil {
	// 		return fmt.Errorf("failed to create request: %w", err)
	// 	}

	// 	// Set headers
	// 	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
	// 	req.Header.Set("Content-Type", "application/json")

	// 	// Execute the request
	// 	client := &http.Client{}
	// 	resp, err := client.Do(req)
	// 	if err != nil {
	// 		return fmt.Errorf("request failed: %w", err)
	// 	}
	// 	defer resp.Body.Close()

	// 	// Check for HTTP errors
	// 	if resp.StatusCode != http.StatusOK {
	// 		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	// 	}

	// 	fmt.Println("Message posted to Mastodon successfully")
	return nil
}
