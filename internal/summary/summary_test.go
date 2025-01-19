package summary

import "testing"

func TestGemtextLink(t *testing.T) {
	const geminiCapsule = "foo.zone"

	table := map[string]string{
		"http://example.com":  "=> http://example.com example.com",
		"https://example.org": "=> https://example.org example.org",
		"https://example.org/some/very/long/link/here?with=a&free=of&parameters=here": "=> https://example.org/some/very/long/link/here?with=a&free=of&parameters=here example.org/s...rameters=here",
		"http://foo.zone":  "=> gemini://foo.zone foo.zone",
		"https://foo.zone": "=> gemini://foo.zone foo.zone",
		"beer://foo.zone":  "=> beer://foo.zone foo.zone",
		"https://foo.zone/some/very/long/link/here?with=a&free=of&parameters=here": "=> gemini://foo.zone/some/very/long/link/here?with=a&free=of&parameters=here foo.zone/some...rameters=here",
	}

	for url, expected := range table {
		if result := gemtextLink(geminiCapsule, url, 30); result != expected {
			t.Errorf("Expected '%s' but got '%s' with input '%s'", expected, result, url)
		}
	}
}
