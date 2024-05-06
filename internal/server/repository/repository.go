package repository

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Repository struct {
	dataDir string
}

func New(dataDir string) Repository {
	return Repository{dataDir}
}

func (r Repository) List() ([]string, error) {
	var ids []string

	visit := func(files *[]string) filepath.WalkFunc {
		return func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
				return nil
			}
			if !info.IsDir() && strings.HasSuffix(path, ".json") {
				*files = append(*files, strings.TrimPrefix(path, r.dataDir))
			}
			return nil
		}
	}

	err := filepath.Walk(r.dataDir, visit(&ids))
	return ids, err
}
