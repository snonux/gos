package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	listenAddr := flag.String("listenAddr", "localhost:8080", "The listen address")
	dataDir := flag.String("dataDir", "./data", "The data directory")

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if err := handleSubmit(w, r, *dataDir); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
		}
	})

	log.Println("Server is starting on ", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, nil); err != err {
		log.Fatal("Error starting server: ", err)
	}
}

func handleSubmit(w http.ResponseWriter, r *http.Request, dataDir string) error {
	if r.Method != "POST" {
		return fmt.Errorf("Expexted POST request")
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/%s/%x.txt", dataDir,
		time.Now().Format("2006"), sha256.Sum256(data))

	if err := saveFile(filePath, data); err != nil {
		return err
	}

	return nil
}

func saveFile(filePath string, bytes []byte) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(filePath, bytes, 0644)
}
