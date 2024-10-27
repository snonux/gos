package linkedin

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

type preview struct {
	title, url string
}

func NewPreview(ctx context.Context, urls []string) (preview, error) {
	if len(urls) == 0 {
		return preview{}, nil
	}
	title, err := fetchTitle(ctx, urls[0])
	if err == nil && title == "" {
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

// fetchTitle fetches the HTML page at the given URL and returns the content of the <title> tag.
func fetchTitle(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get a successful response: %v", resp.StatusCode)
	}

	// Parse the HTML document
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
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
		return "", fmt.Errorf("no title element found")
	}

	return title, nil
}
