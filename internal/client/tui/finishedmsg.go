package tui

type finishedMsg struct {
	callback func() error
	err      error
}

func (f finishedMsg) Error() string {
	return f.err.Error()
}
