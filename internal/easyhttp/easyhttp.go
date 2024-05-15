package easyhttp

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
)

func Get(uri, apiKey string) ([]byte, error) {
	var (
		client = &http.Client{}
		bytes []byte
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
