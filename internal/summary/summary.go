package summary

import (
	"context"
	"fmt"
	"iter"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
)

func Run(ctx context.Context, args config.Args) error {
	entries, err := deduppedEntries(args)
	if err != nil {
		return err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Time.Before(entries[j].Time)
	})

	title := fmt.Sprintf("Posts for %s", strings.Join(args.SummaryFor, " "))
	gemtext, err := fmt.Print(generateGemtext(entries, title, args.GemtexterEnable))
	if err != nil {
		return err
	}
	fmt.Print(gemtext)

	return nil
}

// TODO: Fix the Gemtexter inline toc when there are tags in it
func generateGemtext(entries []entry.Entry, title string, gemtexterEnable bool) (string, error) {
	var (
		sb             strings.Builder
		currentDateStr string
	)

	sb.WriteString("# ")
	sb.WriteString(title)
	if gemtexterEnable {
		sb.WriteString("\n\n<< template::inline::toc")
	}

	for _, en := range entries {
		dateStr := en.Time.Format("January 2006")
		if currentDateStr != dateStr {
			currentDateStr = dateStr
			sb.WriteString("\n\n## ")
			sb.WriteString(currentDateStr)
		}

		content, urls, err := en.Content()
		if err != nil {
			return "", err
		}

		content = prepare(content)
		sb.WriteString("\n\n### ")
		sb.WriteString(firstFewWords(content, 50))

		if err != err {
			return "", err
		}
		sb.WriteString("\n\n")
		sb.WriteString(content)

		if len(urls) > 0 {
			sb.WriteString("\n")
			for _, url := range urls {
				sb.WriteString("\n")
				sb.WriteString(gemtextLink(url, 70))
			}
		}
	}

	if gemtexterEnable {
		sb.WriteString("\n\nOther related posts:")
		sb.WriteString("\n\n<< template::inline::index posts-for")
		sb.WriteString("\n\n")
	}

	return sb.String(), nil
}

func matchingEntries(args config.Args) iter.Seq2[entry.Entry, error] {
	return func(yield func(entry.Entry, error) bool) {
		for _, dateStr := range args.SummaryFor {
			glob := filepath.Join(args.GosDir,
				"db/platforms/*/", fmt.Sprintf("*%s*-??????.posted", dateStr))
			paths, err := filepath.Glob(glob)
			if err != nil && !yield(entry.Zero, err) {
				return
			}
			for _, path := range paths {
				en, err := entry.New(path)
				if !yield(en, err) {
					return
				}
			}
		}
	}
}

func deduppedEntries(args config.Args) ([]entry.Entry, error) {
	dedup := make(map[string]entry.Entry)

	for en, err := range matchingEntries(args) {
		if err != nil {
			return entry.Zeroes, err
		}
		if other, ok := dedup[en.Name()]; ok {
			// If two conflicting entries (e.g. same post for mastodon and linkedin)
			// select the one which was modified latest.
			after, err := other.After(en)
			if err != nil {
				return entry.Zeroes, err
			}
			if after {
				continue
			}
		}
		dedup[en.Name()] = en
	}

	var entries []entry.Entry
	for _, val := range dedup {
		entries = append(entries, val)
	}
	return entries, nil
}

var (
	newlineRegex    = regexp.MustCompile(`\n`)
	urlRegex        = regexp.MustCompile(`https?://\S+`)
	multiSpaceRegex = regexp.MustCompile(`\s{2,}`)
	tagRegex        = regexp.MustCompile(`\B#\w+\b`)
)

func prepare(content string) string {
	content = newlineRegex.ReplaceAllString(content, " ")
	content = urlRegex.ReplaceAllString(content, "")
	content = multiSpaceRegex.ReplaceAllString(content, " ")
	content = strings.TrimSpace(content)
	content = tagRegex.ReplaceAllString(content, "`$0`")
	return content
}

func gemtextLink(url string, maxLen int) string {
	url = strings.TrimSpace(url)
	if len(url) <= maxLen {
		return "=> " + url
	}
	halfLen := (maxLen - 3) / 2
	shorten := url[:halfLen] + "..." + url[len(url)-halfLen:]
	return "=> " + url + " " + shorten
}

func firstFewWords(content string, maxLen int) string {
	words := strings.Fields(content)
	result := ""
	for _, word := range words {
		if len(result)+len(word)+len(" ...") > maxLen {
			break
		}
		if result != "" {
			result += " "
		}
		result += word
	}
	return result + " ..."
}
