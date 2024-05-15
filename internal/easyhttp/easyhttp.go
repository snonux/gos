package easyhttp

import (
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
		return bytes, err
	}

	req.Header.Set("X-API-KEY", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return bytes, err
	}
	defer resp.Body.Close()

	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return bytes, err
	}

	return bytes, nil
}

func GetFromJson[T any](uri, apiKey string) (T, error) {
	var data T

 	bytes, err := Get(uri, apiKey)
 	if err != nil {
 		return data, err
 	}

	err = json.Unmarshal(bytes, &data)
	return data, err
}
