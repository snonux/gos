package server

import (
	"fmt"
	"log"
	"net/http"

	config "codeberg.org/snonux/gos/internal/config/server"
	"codeberg.org/snonux/gos/internal/server/health"
)

type Server struct {
	Status health.Status
	Conf   config.ServerConfig
}

type HandlerFuncWithError func(http.ResponseWriter, *http.Request) error

func New(conf config.ServerConfig, status health.Status) Server {
	return Server{Conf: conf, Status: status}
}

func (serv Server) Handle(name string, handler HandlerFuncWithError) {
	var (
		handlerPath = fmt.Sprintf("/%s", name)
		handlerName = fmt.Sprintf("%sHandler", name)
	)

	http.HandleFunc(handlerPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Someone requested", handlerName)

		// The health endpoint doesn't require an API key
		if handlerName != "healthHandler" {
			accessHealthStatusKey := "server.Handler.Access"
			if r.Header.Get("X-API-KEY") != serv.Conf.APIKey {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				serv.Status.Set(health.Critical, accessHealthStatusKey, fmt.Errorf("Unauthorized access attempt to %s", handlerName))
				return
			}
			serv.Status.Clear(accessHealthStatusKey)
		}

		if err := handler(w, r); err != nil {
			serv.Status.Set(health.Critical, handlerName, err)
			return
		}
		serv.Status.Clear(handlerName)
	})
}
