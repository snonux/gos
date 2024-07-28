package tui

import tea "github.com/charmbracelet/bubbletea"

type finishedMsg struct {
	callback func() error
	err      error
}

func (f finishedMsg) Error() string {
	return f.err.Error()
}

func finished(callback func() error, err error) tea.Cmd {
	return func() tea.Msg {
		return finishedMsg{
			callback: callback,
			err:      err,
		}
	}
}
