package tui

import (
	"context"
	"fmt"
	"os"
	"time"

	"codeberg.org/snonux/gos/internal/config/client"
	config "codeberg.org/snonux/gos/internal/config/client"
	"codeberg.org/snonux/gos/internal/easyhttp"
	"codeberg.org/snonux/gos/internal/types"
	tea "github.com/charmbracelet/bubbletea"
)

func submitAction(ctx context.Context, conf config.ClientConfig) tea.Cmd {
	composeFile := fmt.Sprintf("%s/%s", conf.DataDir, conf.ComposeFile)

	return submitEntry(ctx, conf, composeFile, func() error {
		// This is the callback to call when the entry was submitted succesfully
		timestamp := time.Now().Format("20060102-150405")
		submittedFile := fmt.Sprintf("%s/submitted-%s.txt", conf.DataDir, timestamp)
		return os.Rename(composeFile, submittedFile)
	})
}

func submitEntry(ctx context.Context, conf client.ClientConfig, filePath string, callback func() error) tea.Cmd {
	servers, err := conf.Servers()
	if err == nil {
		var entry types.Entry
		if entry, err = types.NewEntryFromFile(filePath); err == nil {
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
