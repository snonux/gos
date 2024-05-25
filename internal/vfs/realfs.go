package vfs

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type RealFS struct{}

func (RealFS) ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func (RealFS) SaveFile(filePath string, bytes []byte) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(filePath, bytes, 0644)
}

func (RealFS) FindFiles(dataDir, suffix string) ([]string, error) {
	var filePaths []string

	visit := func() filepath.WalkFunc {
		return func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
				return nil
			}
			if info.IsDir() || !strings.HasSuffix(path, suffix) {
				return nil
			}
			filePaths = append(filePaths, path)
			return nil
		}
	}

	return filePaths, filepath.Walk(dataDir, visit())
}
