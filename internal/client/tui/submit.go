package tui

import (
	"fmt"

	"codeberg.org/snonux/gos/internal/config/client"
	config "codeberg.org/snonux/gos/internal/config/client"
	"codeberg.org/snonux/gos/internal/easyhttp"
	"codeberg.org/snonux/gos/internal/types"
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
	servers, err := conf.Servers()
	if err == nil {
		var entry types.Entry
		err = easyhttp.PostData("/submit", conf.APIKey, &entry, servers...)
	}

	return func() tea.Msg {
		return finishedMsg{
			callback: callback,
			err:      err,
		}
	}
}
