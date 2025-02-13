package platforms

import (
	"context"
	"fmt"
	"strings"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/platforms/linkedin"
	"codeberg.org/snonux/gos/internal/platforms/mastodon"
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

func (p Platform) Post(ctx context.Context, args config.Args, sizeLimit int, en entry.Entry) (err error) {
	colour.Infoln("Posting", en)
	switch p.String() {
	case "mastodon":
		err = mastodon.Post(ctx, args, sizeLimit, en)
	case "linkedin":
		err = linkedin.Post(ctx, args, sizeLimit, en)
	default:
		err = fmt.Errorf("Platform '%s' (not yet) implemented", p)
	}

	if err != nil {
		return err
	}
	if err := en.MarkPosted(); err != nil {
		return err
	}

	colour.Successfln("Successfully posted message to %s", p)
	return nil
}

func ExpandAliases(shareTag string) (string, error) {
	parts := strings.Split(shareTag, ":")
	if parts[0] != "share" {
		return "", fmt.Errorf("expected share tag, but got '%s' in '%s'", parts[0], shareTag)
	}

	elems := []string{"share"}
	dedup := make(map[string]struct{}, len(aliases))

	for _, alias := range parts[1:] {
		platformStr, ok := aliases[alias]
		if !ok {
			return "", fmt.Errorf("invalid platform alias '%s' in '%s'", alias, shareTag)
		}
		if _, ok := dedup[platformStr]; ok {
			// Duplicate, ignore
			continue
		}
		elems = append(elems, platformStr)
		dedup[platformStr] = struct{}{}
	}

	return strings.Join(elems, ":"), nil
}
