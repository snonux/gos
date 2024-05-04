package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	listenAddr := flag.String("listenAddr", "localhost:8080", "The listen address")
	dataDir := flag.String("dataDir", "./data", "The data directory")

	if _, err := os.Stat(*dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(*dataDir, 700); err != nil {
			log.Fatal(*dataDir, err)
		}
	}

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've submitted somthing")
	})

	log.Println("Server is starting on ", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, nil); err != err {
		log.Fatal("Error starting server: ", err)
	}
}
