package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/config"
	"github.com/pustserg/secvault/repository"
)

var (
	choices = []string{"generate password", "add entry", "list entries"}
)

type InitialModel struct {
	repo    repository.RepositoryInterface
	cfg     *config.AppConfig
	choices []string
	cursor  int
}

func NewInitialModel(cfg *config.AppConfig, repo repository.RepositoryInterface) InitialModel {
	return InitialModel{
		cfg:     cfg,
		repo:    repo,
		choices: choices,
	}
}

func (m InitialModel) Init() tea.Cmd {
	return nil
}

func (m InitialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "j", "down", "tab":
			m.cursor = (m.cursor + 1) % len(m.choices)
		case "k", "up", "shift+tab":
			m.cursor = (m.cursor - 1 + len(m.choices)) % len(m.choices)
		case " ", "enter":
			switch m.choices[m.cursor] {
			case "generate password":
				return NewGeneratePasswordModel(m, m.cfg), nil
			case "add entry":
				// need to ask the password for storage and then go to the next model
				return NewAskPasswordModel(m, m.repo, AddEntryTargetAction), nil
			case "list entries":
				return NewAskPasswordModel(m, m.repo, ListEntriesTargetAction), nil
			}
		}
	}

	return m, nil
}

func (m InitialModel) View() string {
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
