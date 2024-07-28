package tui

import (
	"context"
	"fmt"

	"codeberg.org/snonux/gos/internal/config/client"
	config "codeberg.org/snonux/gos/internal/config/client"
	"codeberg.org/snonux/gos/internal/easyhttp"
	"codeberg.org/snonux/gos/internal/types"
	tea "github.com/charmbracelet/bubbletea"
)

func submitAction(ctx context.Context, conf config.ClientConfig) tea.Cmd {
	composeFile := fmt.Sprintf("%s/%s", conf.DataDir, conf.ComposeFile)

	return submitEntry(ctx, conf, composeFile, func() error {
		// This is the callback to call
		return nil
	})
}

func submitEntry(ctx context.Context, conf client.ClientConfig, filePath string, callback func() error) tea.Cmd {
	servers, err := conf.Servers()
	if err == nil {
		var entry types.Entry
		entry, err = types.NewEntryFromFile(filePath)
		if err == nil {
			err = easyhttp.PostData(ctx, "submit", conf.APIKey, &entry, servers...)
		}
	}

	return func() tea.Msg {
		return finishedMsg{
			callback: callback,
			err:      err,
		}
	}
}
