package tui

import (
	"fmt"

	config "codeberg.org/snonux/gos/internal/config/client"
	tea "github.com/charmbracelet/bubbletea"
)

func submitAction(conf config.ClientConfig) tea.Cmd {
	// The composed file is now the file to be submitted.
	//submitFile := fmt.Sprintf("%s/%s", conf.DataDir, conf.ComposeFile)

	/*
		c := &http.Client{
			Timeout: 10 * time.Second,
		}
		res, err := c.Get(url)
		if err != nil {
			return sub{err}
		}
		defer res.Body.Close()
	*/

	return func() tea.Msg {
		err := fmt.Errorf("this is a test errorrrrrr!")
		return finishedMsg{
			err: err,
		}
	}

}
