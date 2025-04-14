package linkedin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/oi"
	"golang.org/x/net/html"
)

var (
	errNoTitleElementFound = errors.New("no title element found")
	errNoImageElementFound = errors.New("no image element found")
)

type preview struct {
	title, thumbnailURL, thumbnailDownloadPath, url string
}

func NewPreview(ctx context.Context, args config.Args, urls []string) (preview, error) {
	var (
		p   preview
		err error
	)
	if len(urls) == 0 {
		return p, nil
	}
	p.url = urls[0]

	if p.title, p.thumbnailURL, err = extractFromURL(ctx, urls[0]); err != nil {
		if errors.Is(err, errNoTitleElementFound) || p.title == "" {
			colour.Infoln("Setting title to", urls[0])
			p.title = urls[0]
		}
		if errors.Is(err, errNoImageElementFound) {
			colour.Infoln("URL", urls[0], "is without any image, that's fine, though.")
		}
		if !errors.Is(err, errNoTitleElementFound) && !errors.Is(err, errNoImageElementFound) {
			return p, err
		}
	}

	if p.thumbnailURL != "" {
		if p.thumbnailDownloadPath, err = p.DownloadImage(args.CacheDir); err != nil {
			return p, err
		}
		colour.Infoln("Downloaded preview image to ", p.thumbnailDownloadPath)
	}
	return p, nil
}

func (p preview) String() string {
	if p.thumbnailURL != "" {
		return fmt.Sprintf("Title: %s; URL: %s, Image: %s", p.title, p.url, p.thumbnailURL)
	}
	return fmt.Sprintf("Title: %s; URL: %s", p.title, p.url)
}

func (p preview) TitleAndURL() (string, string, bool) {
	return p.title, p.url, p.url != ""
}

func (p preview) Thumbnail() (string, bool) {
	return p.thumbnailDownloadPath, p.thumbnailDownloadPath != ""
}

func (p preview) DownloadImage(destPath string) (string, error) {
	if err := oi.EnsureDir(destPath); err != nil {
		return "", err
	}
	resp, err := http.Get(p.thumbnailURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status while trying to download image: %s", resp.Status)
	}

	destFile := fmt.Sprintf("%s/%s", destPath, filepath.Base(p.thumbnailURL))
	out, err := os.Create(destFile)
	if err != nil {
		return destFile, fmt.Errorf("%s: %w", destFile, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return destFile, err
	}

	return destFile, nil
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
