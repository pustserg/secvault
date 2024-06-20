package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/repository"
)

type ChooseEntryKindModel struct {
	repo      repository.RepositoryInterface
	prevModel tea.Model
	cursor    int
	choices   []string
	password  string
}

func NewChooseEntryKindModel(prevModel tea.Model, repo repository.RepositoryInterface, password string) ChooseEntryKindModel {
	return ChooseEntryKindModel{
		prevModel: prevModel,
		repo:      repo,
		password:  password,
		choices:   []string{"Note", "Website"},
	}
}

func (m ChooseEntryKindModel) Init() tea.Cmd {
	return nil
}

func (m ChooseEntryKindModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 1 {
				m.cursor++
			}
		case "b":
			return m.prevModel, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.choices[m.cursor] {
			case "Note":
				return NewAddNoteModel(m.prevModel, m.repo, m.password), nil
			case "Website":
				return NewAddWebsiteModel(m.prevModel, m.repo, m.password), nil
			}
		}
	}
	return m, cmd
}

func (m ChooseEntryKindModel) View() string {
	s := "Choose entry kind\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "b Back\n"
	return s
}
