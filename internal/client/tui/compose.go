package tui

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type composeFinishedMsg struct {
	callback func() error
	err      error
}

func composeAction(editor, dataDir string) tea.Cmd {
	// Maybe handle the error explicitly? Should be obvious anyway
	// when the editor can't open the file because directory isn't there
	_ = ensureDirectoryExists(dataDir)
	composeFile := fmt.Sprintf("%s/compose.txt", dataDir)

	return openEditor(editor, composeFile, func() error {
		timestamp := time.Now().Format("20060102-150405")
		queuedFile := fmt.Sprintf("%s/queued-%s.txt", dataDir, timestamp)
		return os.Rename(composeFile, queuedFile)
	})
}

func openEditor(editor, filePath string, callback func() error) tea.Cmd {
	return tea.ExecProcess(exec.Command(editor, filePath), func(err error) tea.Msg {
		return composeFinishedMsg{callback, err}
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
