package tui

import (
	"fmt"

	"codeberg.org/snonux/gos/internal/config/client"
	config "codeberg.org/snonux/gos/internal/config/client"
	tea "github.com/charmbracelet/bubbletea"
)

func submitAction(conf config.ClientConfig) tea.Cmd {
	composeFile := fmt.Sprintf("%s/%s", conf.DataDir, conf.ComposeFile)

	return submitMessage(conf, composeFile, func() error {
		// This is the callback to call
		return nil
	})
}

func submitMessage(conf client.ClientConfig, filePath string, callback func() error) tea.Cmd {
	return func() tea.Msg {
		return finishedMsg{
			callback: callback,
			err:      fmt.Errorf("This is a sample error"),
		}
	}
}
