package tui

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	config "codeberg.org/snonux/gos/internal/config/client"
	tea "github.com/charmbracelet/bubbletea"
)

type composePostAction int

const (
	noPostAction composePostAction = iota
	queueAfterCompose
	submitAfterCompose
)

func composeAction(ctx context.Context, conf config.ClientConfig, postAction composePostAction) tea.Cmd {
	err := ensureDirectoryExists(conf.DataDir)
	composeFile := fmt.Sprintf("%s/%s", conf.DataDir, conf.ComposeFile)
	log.Println("Composing", composeFile)

	return openEditor(conf.Editor, composeFile, func() error {
		if err != nil {
			return err
		}

		switch postAction {
		case submitAfterCompose:
			return submitEntryNoCmd(ctx, conf, composeFile)
		case queueAfterCompose:
			timestamp := time.Now().Format("20060102-150405")
			queuedFile := fmt.Sprintf("%s/queued-%s.txt", conf.DataDir, timestamp)
			return os.Rename(composeFile, queuedFile)
		}

		return nil
	})
}

func openEditor(editor, filePath string, cb func() error) tea.Cmd {
	return tea.ExecProcess(exec.Command(editor, filePath), func(err error) tea.Msg {
		return finishedMsg{
			cb:  cb,
			err: err,
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
