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
		fmt.Fprint(w, health.String())
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if err := handleSubmit(w, r, *dataDir); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			health.set(critical, "submitHandler", err.Error())
			return
		}
		health.clear("submitHandler")
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		if err := handleList(w, r, *dataDir); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			health.set(critical, "listHandler", err.Error())
			return
		}
		health.clear("listHandler")
	})

	log.Println("Server is starting on ", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, nil); err != err {
		log.Fatal("Error starting server: ", err)
	}
}
