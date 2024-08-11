package tui

import tea "github.com/charmbracelet/bubbletea"

type finishedMsg struct {
	cb  func() error
	err error
}

func (f finishedMsg) Error() string {
	return f.err.Error()
}

func finishedCmd(cb func() error, err error) tea.Cmd {
	return func() tea.Msg {
		return finishedMsg{
			cb:  cb,
			err: err,
		}
	}
}
