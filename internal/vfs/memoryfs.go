package vfs

import (
	"fmt"
	"strings"
)

type MemoryFS map[string][]byte

func (fs MemoryFS) ReadFile(filePath string) ([]byte, error) {
	if bytes, ok := fs[filePath]; ok {
		return bytes, nil
	}
	return []byte{}, fmt.Errorf("no such file path: %s", filePath)
}

func (fs MemoryFS) SaveFile(filePath string, bytes []byte) error {
	fs[filePath] = bytes
	return nil
}

func (fs MemoryFS) FindFiles(dataDir, suffix string) ([]string, error) {
	var filePaths []string

	for filePath := range fs {
		if !strings.HasSuffix(filePath, suffix) {
			continue
		}
		filePaths = append(filePaths, filePath)
	}

	return filePaths, nil
}
