package internal

import (
	"os"
	"path/filepath"
)

func SaveFile(filePath string, bytes []byte) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(filePath, bytes, 0644)
}
