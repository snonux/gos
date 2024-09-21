package social

import "codeberg.org/snonux/gos/internal/config/server"

type Mastodon struct {
	conf server.ServerConfig
}

func NewMastodon(conf server.ServerConfig) Mastodon {
	return Mastodon{conf}
}

func Post(content string) error {
	return nil
}
