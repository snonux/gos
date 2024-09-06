package tui

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"codeberg.org/snonux/gos/internal/config/client"
	config "codeberg.org/snonux/gos/internal/config/client"
	"codeberg.org/snonux/gos/internal/easyhttp"
	"codeberg.org/snonux/gos/internal/types"
	tea "github.com/charmbracelet/bubbletea"
)

func submitActionCmd(ctx context.Context, conf config.ClientConfig) tea.Cmd {
	composeFile := fmt.Sprintf("%s/%s", conf.DataDir, conf.ComposeFile)
	log.Println("Submitting", composeFile)

	return submitEntryCmd(ctx, conf, composeFile, func() error {
		// This is the cb to call when the entry was submitted succesfully
		return nil
	})
}

func submitEntryCmd(ctx context.Context, conf client.ClientConfig, composeFile string, cb func() error) tea.Cmd {
	return finishedCmd(cb, submitEntry(ctx, conf, composeFile))
}

func submitEntry(ctx context.Context, conf client.ClientConfig, composeFile string) error {
	if len(conf.Servers) == 0 {
		return errors.New("no server configured")
	}

	entry, err := types.NewEntryFromTextFile(composeFile)
	if err != nil {
		return err
	}

	if err := easyhttp.PostData(ctx, "submit", conf.APIKey, &entry, conf.Servers...); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102-150405")
	submittedFile := fmt.Sprintf("%s/submitted-%s.txt", conf.DataDir, timestamp)

	log.Println("Renaming", composeFile, "to", submittedFile)
	return os.Rename(composeFile, submittedFile)
}
