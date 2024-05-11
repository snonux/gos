package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/server/handle"
	"codeberg.org/snonux/gos/internal/server/health"
)

const healthHandlerName = `healthHandler`

type server struct {
	hs   health.Status
	conf config.Config
}

type handlerFuncWithError func(http.ResponseWriter, *http.Request) error

func (s server) httpHandle(name string, handler handlerFuncWithError) {
	var (
		handlerPath = fmt.Sprintf("/%s", name)
		handlerName = fmt.Sprintf("%sHandler", name)
	)

	http.HandleFunc(handlerPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Someone requested", handlerName)

		// The health endpoint doesn't require an API key
		if handlerName != healthHandlerName && r.Header.Get("X-API-KEY") != s.conf.ApiKey {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			log.Println("Unauthorized access attempt to", handlerName)
			return
		}

		if err := handler(w, r); err != nil {
			s.hs.Set(health.Critical, handlerName, err.Error())
			return
		}
		s.hs.Clear(handlerName)
	})
}

func main() {
	configFile := flag.String("cfg", "/etc/gos.json", "The configuration file")

	conf, err := config.New(*configFile)
	if err != nil {
		log.Fatal("error building config:", err)
	}

	serv := server{
		conf: conf,
		hs:   health.NewStatus(),
	}

	serv.httpHandle("health", func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, serv.hs.String())
		return nil
	})

	serv.httpHandle("submit", func(w http.ResponseWriter, r *http.Request) error {
		return handle.Submit(w, r, serv.conf.DataDir)
	})

	serv.httpHandle("list", func(w http.ResponseWriter, r *http.Request) error {
		return handle.List(w, r, serv.conf.DataDir)
	})

	serv.httpHandle("get", func(w http.ResponseWriter, r *http.Request) error {
		return handle.Get(w, r, serv.conf.DataDir)
	})

	log.Println("Server is starting on", serv.conf.ListenAddr)
	if err := http.ListenAndServe(serv.conf.ListenAddr, nil); err != err {
		log.Fatal("error starting server:", err)
	}
}
