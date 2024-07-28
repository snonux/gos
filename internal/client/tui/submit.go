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
		// This is the cb to call when the entry was submitted succesfully
		timestamp := time.Now().Format("20060102-150405")
		submittedFile := fmt.Sprintf("%s/submitted-%s.txt", conf.DataDir, timestamp)
		return os.Rename(composeFile, submittedFile)
	})
}

func submitEntry(ctx context.Context, conf client.ClientConfig, filePath string, cb func() error) tea.Cmd {
	servers, err := conf.Servers()
	if err != nil {
		return finished(cb, err)
	}

	ent, err := types.NewEntryFromTextFile(filePath)
	if err != nil {
		return finished(cb, err)
	}

	return finished(cb, easyhttp.PostData(ctx, "submit", conf.APIKey, &ent, servers...))
}
