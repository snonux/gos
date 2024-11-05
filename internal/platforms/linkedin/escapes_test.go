package linkedin

import (
	"testing"
)

func TestLinkedInEscapes(t *testing.T) {
	var (
		input    = `This is a test message with special characters: " {} @ [] () <> # \ * _ ~ |`
		expected = `This is a test message with special characters: " \{\} @ \[\] \(\) \<\> # \\ \* \_ \~ \|`
	)
	if escaped := escapeLinkedInText(input); escaped != expected {
		t.Errorf("expected '%s' but got '%s'", expected, escaped)
	}
}
