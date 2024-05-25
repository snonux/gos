package internal

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// virtual file system - useful for testing as well
type VFS interface {
	ReadFile(name string) ([]byte, error)
	SaveFile(filePath string, bytes []byte) error
	FindFiles(dataPath string) ([]string, error)
}

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

func (RealFS) FindFiles(dataDir string) ([]string, error) {
	var filePaths []string

	visit := func() filepath.WalkFunc {
		return func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
				return nil
			}
			if info.IsDir() || !strings.HasSuffix(path, ".json") {
				return nil
			}
			filePaths = append(filePaths, path)
			/*
				entry, err := types.NewEntryFromFile(path)
				if err != err {
					return err
				}
				r.add(entry)
			*/
			return nil
		}
	}

	return filePaths, filepath.Walk(dataDir, visit())
}
