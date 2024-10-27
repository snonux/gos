package linkedin

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// https://learn.microsoft.com/en-us/linkedin/marketing/community-management/shares/little-text-format?view=li-lms-2024-01#language-grammar
func escapeLinkedInText(input string) string {
	var builder strings.Builder

	reservedChars := map[rune]string{
		'"': "\\\"",
		'|': "\\|",
		'{': "\\{",
		'}': "\\}",
		// '@': "\\@",
		'[': "\\[",
		']': "\\]",
		'(': "\\(",
		')': "\\)",
		'<': "\\<",
		'>': "\\>",
		//'#':  "\\#",
		'\\': "\\\\",
		'*':  "\\*",
		'_':  "\\_",
		'~':  "\\~",
	}

	for _, char := range input {
		if escapeSeq, ok := reservedChars[char]; ok {
			builder.WriteString(escapeSeq)
		} else {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}

// fetchTitle fetches the HTML page at the given URL and returns the content of the <title> tag.
func fetchTitle(url string) (string, error) {
	// Send a GET request to the URL
	resp, err := http.Get(url)
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
