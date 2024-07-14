package tui

import (
	"context"
	"fmt"

	config "codeberg.org/snonux/gos/internal/config/client"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFFFF")).
	Background(lipgloss.Color("#0000FF")).
	PaddingTop(2).PaddingLeft(4).PaddingRight(4).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#FFFFFF")).
	BorderBackground(lipgloss.Color("#0000FF"))

var errroStyle = lipgloss.NewStyle().
	Bold(true).
	Align(lipgloss.Center).
	Foreground(lipgloss.Color("#FFFFFF")).
	Background(lipgloss.Color("#FF0000")).
	PaddingTop(0).PaddingBottom(0).PaddingLeft(2).PaddingRight(2).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#FFFFFF")).
	BorderBackground(lipgloss.Color("#FF0000"))

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
	ctx             context.Context
	err             error
}

const (
	composeNewPostCursor = iota
	submitPostCursor
)

func initModel(conf config.ClientConfig) model {
	return model{
		choices: []string{"Compose post", "Submit post"},
		ctx:     context.Background(),
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
			case composeNewPostCursor:
				return m, composeAction(m.conf, false)
			case submitPostCursor:
				return m, submitAction(m.ctx, m.conf)
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
	case finishedMsg:
		m.err = msg.err
		if m.err != nil {
			return m, nil
		}
		if err := msg.callback(); err != nil {
			m.err = err
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

	if m.err != nil {
		s += "\n"
		s += errroStyle.Render(fmt.Sprintf("\nERROR: %s\n", m.err))
		s += "\n"
	}

	s += "\nPress q to quit.\n"
	return style.Render(s)
}
