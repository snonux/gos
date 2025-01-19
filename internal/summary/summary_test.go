package summary

import "testing"

func TestGemtextLink(t *testing.T) {
	const geminiCapsule = "foo.zone"

	table := map[string]string{
		"http://example.com":  "=> http://example.com example.com",
		"https://example.org": "=> https://example.org example.org",
		"https://example.org/some/very/long/link/here?with=a&free=of&parameters=here": "=> https://example.org/some/very/long/link/here?with=a&free=of&parameters=here example.org/s...rameters=here",

		// TODO: Refactor, so that internal links refer to ./ and that the suffix .gmi and/or .html is hidden if possible
		"beer://foo.zone":                               "=> beer://foo.zone foo.zone",
		"http://foo.zone":                               "=> one foo.zone",
		"https://foo.zone/index.html":                   "=> gemini://foo.zone/index.gmi foo.zone/index.gmi",
		"https://foo.zone/gemtext/this-is-awesome.html": "=> gemini://foo.zone/gemtext/this-is-awesome.gmi foo.zone/gemt...s-awesome.gmi",
	}

	for url, expected := range table {
		if result := gemtextLink(geminiCapsule, url, 30); result != expected {
			t.Errorf("Expected '%s' but got '%s' with input '%s'", expected, result, url)
		}
	}
}
