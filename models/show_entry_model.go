package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/repository"
)

type ShowEntryModel struct {
	repo      repository.RepositoryInterface
	prevModel tea.Model
	entry     repository.Entry
	password  string
}

func NewShowEntryModel(prevModel tea.Model, repo repository.RepositoryInterface, entry repository.Entry, password string) ShowEntryModel {
	return ShowEntryModel{
		prevModel: prevModel,
		repo:      repo,
		entry:     entry,
		password:  password,
	}
}

func (m ShowEntryModel) Init() tea.Cmd {
	return nil
}

func (m ShowEntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			return m.prevModel, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ShowEntryModel) View() string {
	s := "Show Entry\n\n"

	if m.entry.Name == "" {
		s += "Unnamed entry\n"
	} else {
		s += fmt.Sprintf("Name: %s\n", m.entry.Name)
	}

	if m.entry.UserName != "" {
		s += fmt.Sprintf("Username: %s\n", m.entry.UserName)
	}

	if m.entry.Password != "" {
		s += fmt.Sprintf("Password: %s\n", m.entry.Password)
	}

	if m.entry.URL != "" {
		s += fmt.Sprintf("URL: %s\n", m.entry.URL)
	}

	if m.entry.Note != "" {
		s += fmt.Sprintf("Note: %s\n", m.entry.Note)
	}

	if m.entry.TotpToken != "" {
		s += fmt.Sprintf("TOTP Token: %s\n", m.entry.TotpToken)
	}

	s += "\nPress 'b' to go back, 'q' to quit\n"
	return s
}
