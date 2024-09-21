package internal

import (
	"context"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/queue"
)

func Run(ctx context.Context, args config.Args) error {
	return queue.Run(args)
}
