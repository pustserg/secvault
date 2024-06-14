package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/config"
)

type Model struct {
	cfg     *config.AppConfig
	choices []string
	cursor  int
}

func NewInitialModel(cfg *config.AppConfig) Model {
	return Model{
		cfg:     cfg,
		choices: []string{"generate password"},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "j", "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case " ", "enter":
			switch m.choices[m.cursor] {
			case "generate password":
				return NewGeneratePasswordModel(m, m.cfg), nil
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	s := "What are we going to do today?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\npress 'q' to quit\n"
	return s
}
