package server

import (
	"fmt"
	"log"
	"net/http"

	config "codeberg.org/snonux/gos/internal/config/server"
	"codeberg.org/snonux/gos/internal/server/health"
)

const HealthHandlerName = `healthHandler`

type Server struct {
	Status health.Status
	Conf   config.ServerConfig
}

type HandlerFuncWithError func(http.ResponseWriter, *http.Request) error

func New(conf config.ServerConfig) Server {
	serv := Server{
		Conf:   conf,
		Status: health.NewStatus(),
	}

	return serv
}

func (serv Server) Handle(name string, handler HandlerFuncWithError) {
	var (
		handlerPath = fmt.Sprintf("/%s", name)
		handlerName = fmt.Sprintf("%sHandler", name)
	)

	http.HandleFunc(handlerPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Someone requested", handlerName)

		// The health endpoint doesn't require an API key
		if handlerName != HealthHandlerName && r.Header.Get("X-API-KEY") != serv.Conf.APIKey {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			log.Println("Unauthorized access attempt to", handlerName)
			return
		}

		if err := handler(w, r); err != nil {
			log.Println(err)
			serv.Status.Set(health.Critical, handlerName, err.Error())
			return
		}
		serv.Status.Clear(handlerName)
	})
}
