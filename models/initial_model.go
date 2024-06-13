package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	choices  []string
	cursor   int
	selected map[int]string
}

func NewInitialModel() Model {
	return Model{
		choices:  []string{"one", "two", "three", "four", "five"},
		selected: map[int]string{},
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
		case " ":
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = m.choices[m.cursor]
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	s := "what are we selecting today?\n\n"
	s += fmt.Sprintf("selected: %v\n\n", m.selected)

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	s += "\npress 'q' to quit\n"
	return s
}
