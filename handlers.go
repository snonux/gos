package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func handleSubmit(w http.ResponseWriter, r *http.Request, dataDir string) error {
	if r.Method != "POST" {
		return fmt.Errorf("expexted POST request")
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	entry, err := newEntry(bytes)
	if err != nil {
		return err
	}
	filePath := fmt.Sprintf("%s/%s/%s.json", dataDir, time.Now().Format("2006"), entry.id)

	jsonStr, err := entry.serialize()
	if err != nil {
		return err
	}

	if err := saveFile(filePath, jsonStr); err != nil {
		return err
	}

	return nil
}

func handleList(w http.ResponseWriter, r *http.Request, dataDir string) error {
	if r.Method != "GET" {
		return fmt.Errorf("expexted GET request")
	}

	repository := newRepository(dataDir)
	ids, err := repository.list()
	if err != nil {
		return err
	}
	for _, id := range ids {
		fmt.Fprintf(w, id)
	}
	return nil
}
