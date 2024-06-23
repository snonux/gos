package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

type editorFinishedMsg struct{ err error }

func openEditor(editor, filePath string) tea.Cmd {
	// Maybe handle the error explicitly? Should be obvious anyway
	// when the editor can't open the file because directory isn't there
	_ = ensureDirectoryExists(filepath.Dir(filePath))

	c := exec.Command(editor, filePath) //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

func ensureDirectoryExists(dir string) error {
	info, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	if info.IsDir() {
		return nil
	}
	return fmt.Errorf("path %s is not a directory", dir)
}
