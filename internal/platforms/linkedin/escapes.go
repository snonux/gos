package linkedin

import (
	"strings"
)

// https://learn.microsoft.com/en-us/linkedin/marketing/community-management/shares/little-text-format?view=li-lms-2024-01#language-grammar
func escapeLinkedInText(input string) string {
	var builder strings.Builder

	// TODO: '"' escapes don't work correctly yet, they show up as \"...\" in LinkedIn posts
	reservedChars := map[rune]string{
		'"': "\\\"", // Just remove this line to fix the above?
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
