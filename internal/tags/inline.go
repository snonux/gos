package tags

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"codeberg.org/snonux/gos/internal/oi"
	"codeberg.org/snonux/gos/internal/platforms"
)

var inlineTagRE = regexp.MustCompile(`[.,:]`)

// Extracts the inline tags from the content ant inserts them into the file path.
func InlineExtract(filePath string) (string, error) {
	content, err := oi.SlurpAndTrim(filePath)
	if err != nil {
		return "", err
	}

	newFilePath, newContent, err := inlineExtractTagsToFilePath(filePath, content)
	if err != nil {
		return "", err
	}
	if newFilePath == filePath {
		return filePath, nil
	}

	fmt.Print("\n")
	if err := oi.WriteFile(newFilePath, newContent); err != nil {
		return "", err
	}
	return newFilePath, os.Remove(filePath)
}

func inlineExtractTagsToFilePath(filePath, content string) (string, string, error) {
	tags, newContent, err := inlineExtractTagsFromContent(content)
	if err != nil {
		return filePath, content, err
	}
	if len(tags) == 0 {
		return filePath, content, nil
	}

	parts := strings.Split(strings.TrimSuffix(filePath, filepath.Ext(filePath)), ".")
	parts = append(parts, tags...)
	parts = append(parts, "extracted")
	parts = append(parts, "txt")

	newFilePath := strings.Join(parts, ".")
	return newFilePath, newContent, nil
}

func inlineExtractTagsFromContent(content string) ([]string, string, error) {
	isShare := func(tag string) bool {
		return strings.HasPrefix(tag, "share:")
	}
	parts := strings.Split(content, " ")
	// First word must contain certain symbols to clarify as (inline) tags.
	if inlineTagRE.MatchString(parts[0]) {
		var tags []string
		// String separator either a dot or a comma. Each element will be a tag.
		for _, elem := range strings.Split(parts[0], ".") {
			tags = append(tags, strings.Split(elem, ",")...)
		}
		if len(tags) > 0 {
			for i := range len(tags) {
				if isShare(tags[i]) {
					var err error
					if tags[i], err = platforms.ExpandAliases(tags[i]); err != nil {
						return []string{}, content, err
					}
				}
			}
			return tags, strings.Join(parts[1:], " "), nil
		}
	}
	return []string{}, content, nil
}
