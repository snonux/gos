package queue

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/oi"
)

// Extracts the inline tags into the filepath and removes them from the content.
func extractInlineTags(filePath string) (string, error) {
	content, err := oi.SlurpAndTrim(filePath)
	if err != nil {
		return "", err
	}

	newFilePath, newContent, err := extractInlineTagsToFilePath(filePath, content)
	if err != nil {
		return "", err
	}
	if newFilePath == filePath {
		return filePath, nil
	}

	colour.Infof("Rewriting path '%s' to '%s' (inline tag extraction)", filePath, newFilePath)
	fmt.Print("\n")
	if err := oi.WriteFile(newFilePath, newContent); err != nil {
		return "", err
	}
	return newFilePath, os.Remove(filePath)
}

func extractInlineTagsToFilePath(filePath, content string) (string, string, error) {
	tags, newContent := extractInlineTagsFromContent(content)
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

func extractInlineTagsFromContent(content string) ([]string, string) {
	parts := strings.Split(content, " ")
	// If the first word of the content contains a dot or comma and there are
	// more than 2 elems, then there are inline tags!
	if strings.Contains(parts[0], ".") || strings.Contains(parts[0], ",") {
		var tags []string
		for _, elem := range strings.Split(parts[0], ".") {
			tags = append(tags, strings.Split(elem, ",")...)
		}
		if len(tags) > 1 {
			return tags, strings.Join(parts[1:], " ")
		}
	}
	return []string{}, content
}
