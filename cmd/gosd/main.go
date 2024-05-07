package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"codeberg.org/snonux/gos/internal/server/handle"
	"codeberg.org/snonux/gos/internal/server/health"
)

const apiKey = "banana" // for dev purposes only, will be changed to something else
const healthHandlerName = `healthHandler`

var hs = health.NewStatus()

type handlerFuncWithError func(http.ResponseWriter, *http.Request) error

func httpHandle(name string, handler handlerFuncWithError) {
	var (
		handlerPath = fmt.Sprintf("/%s", name)
		handlerName = fmt.Sprintf("%sHandler", name)
	)

	http.HandleFunc(handlerPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Someone requested", handlerName)

		// The health endpoint doesn't require an API key
		if handlerName != healthHandlerName && r.Header.Get("X-API-KEY") != apiKey {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			log.Println("Unauthorized access attempt to", handlerName)
			return
		}

		if err := handler(w, r); err != nil {
			hs.Set(health.Critical, handlerName, err.Error())
			return
		}
		hs.Clear(handlerName)
	})
}

func main() {
	listenAddr := flag.String("listenAddr", "localhost:8080", "The listen address")
	dataDir := flag.String("dataDir", "data", "The data directory")

	httpHandle("health", func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, hs.String())
		return nil
	})

	httpHandle("submit", func(w http.ResponseWriter, r *http.Request) error {
		return handle.Submit(w, r, *dataDir)
	})

	httpHandle("list", func(w http.ResponseWriter, r *http.Request) error {
		return handle.List(w, r, *dataDir)
	})

	httpHandle("get", func(w http.ResponseWriter, r *http.Request) error {
		return handle.Get(w, r, *dataDir)
	})

	log.Println("Server is starting on", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, nil); err != err {
		log.Fatal("Error starting server: ", err)
	}
}
