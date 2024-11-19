package platforms

import (
	"fmt"
	"strings"
)

type Platform string

var aliases = map[string]string{
	"linkedin": "linkedin",
	"li":       "linkedin",
	"mastodon": "mastodon",
	"ma":       "mastodon",
	"xcom":     "xcom",
	"x":        "xcom",
	"twitter":  "xcom",
	"tw":       "xcom",
}

func New(platformStr string) (Platform, error) {
	var p Platform
	name, ok := aliases[strings.ToLower(platformStr)]
	if !ok {
		return p, fmt.Errorf("no such platform: '%s'", platformStr)
	}
	return Platform(name), nil
}

func (p Platform) String() string {
	return string(p)
}
