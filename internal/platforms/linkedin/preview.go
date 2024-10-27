package linkedin

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

var errNoTitleElementFound = errors.New("no title element found")

type preview struct {
	title, url string
}

func NewPreview(ctx context.Context, urls []string) (preview, error) {
	if len(urls) == 0 {
		return preview{}, nil
	}
	title, err := fetchHTMLTitle(ctx, urls[0])
	if errors.Is(err, errNoTitleElementFound) || (err == nil && title == "") {
		log.Println("Setting title to", urls[0])
		title = urls[0]
	}
	return preview{title: title, url: urls[0]}, err
}

func (p preview) String() string {
	return fmt.Sprintf("Title: %s; URL: %s", p.title, p.url)
}

func (p preview) Empty() bool {
	return p.url == ""
}

func fetchHTMLTitle(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get a successful response: %v", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Traverse the document and find the <title> tag
	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	// Call the function to search for the title
	f(doc)

	if title == "" {
		return "", errNoTitleElementFound
	}
	return title, nil
}
