package linkedin

import (
	"strings"
	"testing"
)

func TestPreviewExtract(t *testing.T) {

	expectedTitle := "Baz"
	expectedImageURL := "https://free.beer:666/bar/foo.jpg"
	mockHTML := strings.NewReader(`
<!DOCTYPE html>
<html>
<head>
    <title>Baz</title>
</head>
<body>
    <img src="./foo.jpg" alt="Foo">
</body>
</html>
`)

	title, imageURL, err := extract("https://free.beer:666/bar/", mockHTML)
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
