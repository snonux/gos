package scheduler

import (
	"context"
	"log"

	"codeberg.org/snonux/gos/internal/config/server"
)

// TODO: Finish implementing this
func Run(ctx context.Context, config server.ServerConfig) error {
	for _, platform := range config.SocialPlatformsEnabled {
		log.Println("TODO: implement ... posting a post now or what on", platform)
	}
	return nil
}
