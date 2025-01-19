package summary

import "testing"

func TestGemtextLink(t *testing.T) {
	geminiCapsules := []string{"foo.zone"}

	table := map[string]string{
		"http://example.com":  "=> http://example.com example.com",
		"https://example.org": "=> https://example.org example.org",
		"https://example.org/some/very/long/link/here?with=a&free=of&parameters=here": "=> https://example.org/some/very/long/link/here?with=a&free=of&parameters=here example.org/s...rameters=here",

		"beer://foo.zone":             "=> beer://foo.zone foo.zone",
		"http://foo.zone":             "=> gemini://foo.zone foo.zone (Gemini)\n=> http://foo.zone foo.zone",
		"https://foo.zone/index.html": "=> gemini://foo.zone/index.gmi foo.zone/index.gmi (Gemini)\n=> https://foo.zone/index.html foo.zone/index.html",

		"https://foo.zone/gemtext/this-is-awesome.html": "=> gemini://foo.zone/gemtext/this-is-awesome.gmi foo.zone/gemt...s-awesome.gmi (Gemini)\n=> https://foo.zone/gemtext/this-is-awesome.html foo.zone/gemt...-awesome.html",
	}

	for url, expected := range table {
		if result := gemtextLink(geminiCapsules, url, 30); result != expected {
			t.Errorf("Expected '%s' but got '%s' with input '%s'", expected, result, url)
		}
	}
}
