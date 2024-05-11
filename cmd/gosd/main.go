package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/server"
	"codeberg.org/snonux/gos/internal/server/handle"
)

func main() {
	configFile := flag.String("cfg", "/etc/gos.json", "The configuration file")

	conf, err := config.New(*configFile)
	if err != nil {
		log.Fatal("error building config:", err)
	}

	serv := server.New(conf)

	serv.Handle("health", func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, serv.Status.String())
		return nil
	})

	serv.Handle("submit", func(w http.ResponseWriter, r *http.Request) error {
		return handle.Submit(w, r, conf.DataDir)
	})

	serv.Handle("list", func(w http.ResponseWriter, r *http.Request) error {
		return handle.List(w, r, conf.DataDir)
	})

	serv.Handle("get", func(w http.ResponseWriter, r *http.Request) error {
		return handle.Get(w, r, conf.DataDir)
	})

	log.Println("Server is starting on", conf.ListenAddr)
	if err := http.ListenAndServe(conf.ListenAddr, nil); err != err {
		log.Fatal("error starting server:", err)
	}
}
