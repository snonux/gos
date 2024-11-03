package linkedin

import (
	"context"
	"testing"
)

// TODO: Mock the http request, and serve a local HTML page!
func TestFetchHTMLTitleAndFirstImage(t *testing.T) {
	var (
		url              = "https://foo.zone/about/"
		expectedTitle    = "About"
		expectedImageURL = "https://foo.zone/about/paul.jpg"
	)

	title, imageURL, err := fetchHTMLTitleAndFirstImage(context.Background(), url)
	if err != nil {
		t.Error(err)
	}
	if title != expectedTitle {
		t.Errorf("expected title '%s' but got '%s'", expectedTitle, title)
	}
	if imageURL != expectedImageURL {
		t.Errorf("expected imageURL '%s' but got '%s'", expectedImageURL, imageURL)
	}
}
