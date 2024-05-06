package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	listenAddr := flag.String("listenAddr", "localhost:8080", "The listen address")
	dataDir := flag.String("dataDir", "data", "The data directory")
	health := newHealthStatus()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Someone requested /health")
		fmt.Fprint(w, health.String())
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Someone requested /submit")
		if err := handleSubmit(w, r, *dataDir); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			health.set(critical, "submitHandler", err.Error())
			return
		}
		health.clear("submitHandler")
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Someone requested /list")
		if err := handleList(w, r, *dataDir); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			health.set(critical, "listHandler", err.Error())
			return
		}
		health.clear("listHandler")
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Someone requested /get")
		if err := handleGet(w, r, *dataDir); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			health.set(critical, "getHandler", err.Error())
			return
		}
		health.clear("getHandler")
	})

	log.Println("Server is starting on ", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, nil); err != err {
		log.Fatal("Error starting server: ", err)
	}
}
