package linkedin

import (
	"regexp"
	"strings"
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

// extractURLs finds all occurrences of URLs starting with "http://" or "https://" in a given string.
func extractURLs(input string) []string {
	// Regular expression pattern to match URLs starting with http:// or https://
	urlPattern := `(http://|https://)[^\s]+`

	// Compile the regular expression
	re := regexp.MustCompile(urlPattern)

	// Find all matches in the input string
	urls := re.FindAllString(input, -1)

	return urls
}
