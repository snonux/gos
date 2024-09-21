package server

import (
	"fmt"
	"log"
	"os"

	"codeberg.org/snonux/gos/internal/config"
)

type SecretsConfig struct {
	MastodonEnable      bool   `json:"MastodonEnable,omitempty"`
	MastodonDomain      string `json:"MastodonDomain,omitempty"`
	MastodonAccessToken string `json:"MastodonAccessToken,omitempty"`
}

func newSecretsConfig(secretsFile string) (SecretsConfig, error) {
	if isWorldReadable(secretsFile) {
		return SecretsConfig{}, fmt.Errorf("config '%s' is world readable", secretsFile)
	}

	conf, err := config.FromFile[SecretsConfig](secretsFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return conf, err
		}
		log.Println("Skipping config file:", err)
	}

	conf.MastodonEnable = config.Bool("GOS_MASTODON_ENABLE", conf.MastodonEnable)
	conf.MastodonDomain = config.Str("GOS_MASTODON_DOMAIN", conf.MastodonDomain)
	conf.MastodonAccessToken = config.Str("GOS_MASTODON_ACCESS_TOKEN", conf.MastodonAccessToken)

	return conf, nil
}

func isWorldReadable(file string) bool {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return false
	}

	return fileInfo.Mode().Perm()&00004 != 0
}
