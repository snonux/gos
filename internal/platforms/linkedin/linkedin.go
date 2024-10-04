package linkedin

import (
	"context"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
)

func Post(ctx context.Context, args config.Args, ent entry.Entry) error {
	return oauth(args)
}
