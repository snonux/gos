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

	title := fmt.Sprintf("Posts for %s", strings.Join(args.GeminiSummaryFor, " "))
	gemtext, err := fmt.Print(generateGemtext(args, entries, title))
	if err != nil {
		return err
	}
	fmt.Println(gemtext)

	return nil
}

func generateGemtext(args config.Args, entries []entry.Entry, title string) (string, error) {
	var (
		sb             strings.Builder
		currentDateStr string
	)

	sb.WriteString("# ")
	sb.WriteString(title)
	if args.GemtexterEnable {
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
				sb.WriteString(gemtextLink(args.GeminiCapsules, url, 30))
			}
		}
	}

	if args.GemtexterEnable {
		sb.WriteString("\n\nOther related posts:")
		sb.WriteString("\n\n<< template::inline::index posts-from")
		sb.WriteString("\n\n")
	}

	return sb.String(), nil
}

func matchingEntries(args config.Args) iter.Seq2[entry.Entry, error] {
	return func(yield func(entry.Entry, error) bool) {
		for _, dateStr := range args.GeminiSummaryFor {
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

func gemtextLink(geminiCapsules []string, url string, maxLen int) string {
	url = strings.TrimSpace(url)
	var (
		urlNoProto = regexp.MustCompile(`^[a-zA-Z]+://`).ReplaceAllString(url, "")
		links      []string
	)

	// Check whether any element of the slice starts with prefix.
	hasPrefix := func(str string, prefixes []string) bool {
		for _, element := range prefixes {
			if strings.HasPrefix(str, element) {
				return true
			}
		}
		return false
	}

	// Shorten the link description if too long.
	shorten := func(url, urlNoProto string, maxLen int) string {
		if len(urlNoProto) <= maxLen {
			return "=> " + url + " " + urlNoProto
		}
		halfLen := (maxLen - 3) / 2
		shortened := urlNoProto[:halfLen] + "..." + urlNoProto[len(urlNoProto)-halfLen:]
		return "=> " + url + " " + shortened
	}

	// Is this a Gemini link? If so, add it to the link list.
	if hasPrefix(urlNoProto, geminiCapsules) && hasPrefix(url, []string{"http://", "https://"}) {
		urlNoProto := strings.ReplaceAll(urlNoProto, ".html", ".gmi")
		url := "gemini://" + urlNoProto
		links = append(links, shorten(url, urlNoProto, maxLen)+" (Gemini)")
	}

	links = append(links, shorten(url, urlNoProto, maxLen))
	return strings.Join(links, "\n")
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
