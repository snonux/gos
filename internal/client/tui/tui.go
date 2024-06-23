package tui

import (
	"fmt"

	config "codeberg.org/snonux/gos/internal/config/client"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(40)

func Run(conf config.ClientConfig) error {
	p := tea.NewProgram(initModel(conf))
	_, err := p.Run()
	return err
}

type model struct {
	choices         []string
	cursor          int
	conf            config.ClientConfig
	altscreenActive bool
	err             error
}

const (
	composePostCursorPos = 0
)

func initModel(conf config.ClientConfig) model {
	return model{
		choices: []string{"Compose post", "Schedule post"},
		conf:    conf,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			switch m.cursor {
			case composePostCursorPos:
				tmpFile := fmt.Sprintf("%s/compose.md", m.conf.DataDir)
				return m, openEditor(m.conf.Editor, tmpFile)
			}

		case "a":
			m.altscreenActive = !m.altscreenActive
			cmd := tea.EnterAltScreen
			if !m.altscreenActive {
				cmd = tea.ExitAltScreen
			}
			return m, cmd
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Please choose your destiny\n\n"

	for i, choice := range m.choices {
		cursor := "   " // no cursor
		if m.cursor == i {
			cursor = "==>"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quiet.\n"
	return style.Render(s)
}
