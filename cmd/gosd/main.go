package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	config "codeberg.org/snonux/gos/internal/config/server"
	"codeberg.org/snonux/gos/internal/server"
	"codeberg.org/snonux/gos/internal/server/cron"
	"codeberg.org/snonux/gos/internal/server/handler"
)

func main() {
	configFile := flag.String("cfg", "/etc/gosd.json", "The configuration file")

	conf, err := config.New(*configFile)
	if err != nil {
		log.Fatal("error building config:", err)
	}

	var (
		serv = server.New(conf)
		hand = handler.New(conf)
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := cron.Start(ctx, conf, hand); err != nil {
		panic(err)
	}

	serv.Handle("health", func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, serv.Status.String())
		return nil
	})

	serv.Handle("submit", func(w http.ResponseWriter, r *http.Request) error {
		return hand.Submit(ctx, w, r)
	})

	serv.Handle("list", func(w http.ResponseWriter, r *http.Request) error {
		return hand.List(w, r)
	})

	serv.Handle("get", func(w http.ResponseWriter, r *http.Request) error {
		return hand.Get(w, r)
	})

	serv.Handle("merge", func(w http.ResponseWriter, r *http.Request) error {
		return hand.Merge(ctx, w, r)
	})

	log.Println("Server is starting on", conf.ListenAddr)
	if err := http.ListenAndServe(conf.ListenAddr, nil); err != err {
		log.Fatal("error starting server:", err)
	}
}
