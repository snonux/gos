package linkedin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

var (
	errNoTitleElementFound = errors.New("no title element found")
	errNoImageElementFound = errors.New("no image element found")
)

type preview struct {
	title, imageURL, url string
}

func NewPreview(ctx context.Context, urls []string) (preview, error) {
	if len(urls) == 0 {
		return preview{}, nil
	}
	title, imageURL, err := extractFromURL(ctx, urls[0])
	if errors.Is(err, errNoTitleElementFound) || title == "" {
		log.Println("Setting title to", urls[0])
		title = urls[0]
	}
	return preview{title: title, imageURL: imageURL, url: urls[0]}, err
}

func (p preview) String() string {
	return fmt.Sprintf("Title: %s; URL: %s", p.title, p.url)
}

func (p preview) Empty() bool {
	return p.url == ""
}

func findTitle(n *html.Node) (string, error) {
	var title string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(n)
	if title == "" {
		return "", errNoTitleElementFound
	}
	return title, nil
}

func findFirstImage(n *html.Node) (string, error) {
	var imageURL string
	var traverse func(*html.Node) bool
	traverse = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					imageURL = attr.Val
					return true // Stop searching when the first image's URL is found
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if traverse(c) {
				return true
			}
		}
		return false
	}
	if !traverse(n) {
		return "", errNoImageElementFound
	}
	return imageURL, nil
}

func resolveURL(baseURL, rawURL string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse raw URL: %w", err)
	}
	return base.ResolveReference(u).String(), nil
}

func extractFromURL(ctx context.Context, url string) (string, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to get URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to get a successful response: %v", resp.StatusCode)
	}

	return extract(url, resp.Body)
}

func extract(url string, htmlBody io.Reader) (string, string, error) {
	doc, err := html.Parse(htmlBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var errs error
	title, err := findTitle(doc)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	imageURL, err := findFirstImage(doc)
	if err != nil {
		errs = errors.Join(errs, err)
	} else if imageURL != "" {
		if imageURL, err = resolveURL(url, imageURL); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return title, imageURL, errs
}
