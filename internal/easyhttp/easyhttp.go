package easyhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func Get(uri, apiKey string) ([]byte, error) {
	var (
		client = &http.Client{}
		bytes  []byte
	)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return bytes, fmt.Errorf("%s: %w", uri, err)
	}

	req.Header.Set("X-API-KEY", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return bytes, fmt.Errorf("%s: %w", uri, err)
	}
	defer resp.Body.Close()

	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return bytes, fmt.Errorf("%s: %w", uri, err)
	}

	return bytes, nil
}

// Get data from JSON
func GetData[T any](uri, apiKey string, data *T) error {
	bytes, err := Get(uri, apiKey)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, data)
}

func Post(uri, apiKey string, data []byte) ([]byte, error) {
	// TODO: Use contexts in Post and Get requests, e.g. NewRequestWithContext
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	if err != nil {
		return []byte{}, fmt.Errorf("%s: %w", uri, err)
	}

	req.Header.Set("X-API-KEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("%s: %w", uri, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("%s: %w", uri, err)
	}

	return body, nil
}

// Submit structure as JSON to API
func PostData[T any](uri, apiKey string, data *T, servers ...string) error {
	if len(servers) == 0 {
		return fmt.Errorf("no server configured")
	}
	var errs safErrors
	var wg sync.WaitGroup

	for _, server := range servers {
		wg.Add(1)
		go func(server string) {
			defer wg.Done()
			errs.Append(postData[T](fmt.Sprintf("%s/%s", server, uri), apiKey, data))
		}(server)
	}

	wg.Wait()
	return errs.Join()
}

func postData[T any](uri, apiKey string, data *T) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = Post(uri, apiKey, jsonData)
	return err
}
