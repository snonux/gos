package vfs

import (
	"slices"
	"strings"
	"testing"
)

func TestMemoryFS(t *testing.T) {
	t.Parallel()
	var vfs VFS = make(MemoryFS)

	writeFiles := map[string]string{
		"/data/dir/foo.json":        "hello world",
		"/data/dir/subdir/bar.json": "hello solar system",
		"/data/dir/subdir/baz.json": "hello universe",
		"/data/dir/subdir/bay.txt":  "hello bar keeper",
	}

	for path, content := range writeFiles {
		bytes := []byte(content)
		_ = vfs.WriteFile(path, bytes)
	}

	t.Run("files are there", func(t *testing.T) {
		testFilesAreThere(t, vfs, writeFiles)
	})

	t.Run("find json files", func(t *testing.T) {
		testFindFiles(t, vfs, writeFiles, "/data/dir/subdir", ".json", 2)
	})
}

func testFindFiles(t *testing.T, vfs VFS, writeFiles map[string]string, dataDir, suffix string, count int) {
	filePaths, err := vfs.FindFiles(dataDir, suffix)
	if err != nil {
		t.Error(err)
		return
	}

	if len(filePaths) != count {
		t.Error("expected", count, "json files, but got", filePaths)
		return
	}

	for filePath := range writeFiles {
		if !strings.HasPrefix(filePath, dataDir) || !strings.HasSuffix(filePath, suffix) {
			continue
		}

		if !slices.Contains(filePaths, filePath) {
			t.Error("expected file", filePath, "to be there, but it isn't in", filePaths)
			return
		}
	}
}

func testFilesAreThere(t *testing.T, vfs VFS, writeFiles map[string]string) {
	for path, content := range writeFiles {
		bytes, err := vfs.ReadFile(path)
		if err != nil {
			t.Error(err)
			return
		}
		if content != string(bytes) {
			t.Error("expected", content, "in file", path, "but got", string(bytes))
			return
		}
	}
}
