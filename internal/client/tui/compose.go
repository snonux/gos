package tui

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	config "codeberg.org/snonux/gos/internal/config/client"
	tea "github.com/charmbracelet/bubbletea"
)

func composeAction(conf config.ClientConfig, queue bool) tea.Cmd {
	err := ensureDirectoryExists(conf.DataDir)
	composeFile := fmt.Sprintf("%s/%s", conf.DataDir, conf.ComposeFile)

	return openEditor(conf.Editor, composeFile, func() error {
		if err != nil {
			return err
		}
		// ye
		if !queue {
			return nil
		}
		timestamp := time.Now().Format("20060102-150405")
		queuedFile := fmt.Sprintf("%s/queued-%s.txt", conf.DataDir, timestamp)
		return os.Rename(composeFile, queuedFile)
	})
}

func openEditor(editor, filePath string, cb func() error) tea.Cmd {
	return tea.ExecProcess(exec.Command(editor, filePath), func(err error) tea.Msg {
		return finishedMsg{
			cb: cb,
			err:      err,
		}
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
