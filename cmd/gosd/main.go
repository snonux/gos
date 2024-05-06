package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"codeberg.org/snonux/gos/internal/server/handle"
	"codeberg.org/snonux/gos/internal/server/health"
)

func main() {
	listenAddr := flag.String("listenAddr", "localhost:8080", "The listen address")
	dataDir := flag.String("dataDir", "data", "The data directory")
	hs := health.NewStatus()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Someone requested /health")
		fmt.Fprint(w, hs.String())
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Someone requested /submit")
		if err := handle.Submit(w, r, *dataDir); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			hs.Set(health.Critical, "submitHandler", err.Error())
			return
		}
		hs.Clear("submitHandler")
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Someone requested /list")
		if err := handle.List(w, r, *dataDir); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			hs.Set(health.Critical, "listHandler", err.Error())
			return
		}
		hs.Clear("listHandler")
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Someone requested /get")
		if err := handle.Get(w, r, *dataDir); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			hs.Set(health.Critical, "getHandler", err.Error())
			return
		}
		hs.Clear("getHandler")
	})

	log.Println("Server is starting on ", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, nil); err != err {
		log.Fatal("Error starting server: ", err)
	}
}
