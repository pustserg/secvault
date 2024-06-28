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

func NewShowEntryModel(prevModel tea.Model, repo repository.RepositoryInterface, entryID string, password string) ShowEntryModel {
	entry, err := repo.Get(entryID, password)
	if err != nil {
		panic(err)
	}
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
		case tea.KeyEsc.String(), "b":
			return m.prevModel, nil
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit
		case "e":
			switch m.entry.Kind {
			case repository.NoteType:
				return NewEditNoteModel(m.entry, m, m.repo, m.password), nil
			case repository.WebsiteType:
				return m, nil
			}
		case tea.KeyBackspace.String(), tea.KeyDelete.String(), "d":
			// After deleting the entry, we want to go back to the previous model (entries list)
			return NewConfirmationModel(m.prevModel, UpdateEntriesCmd, "Are you sure you want to delete this entry?", []string{"y", "n"}, func() error {
				err := m.repo.Delete(m.entry.ID, m.password)
				if err != nil {
					return err
				}
				return nil
			}), nil
		}
		return m, nil
	case string:
		switch msg {
		case UpdateEntryMsg:
			entry, err := m.repo.Get(m.entry.ID, m.password)
			if err != nil {
				return m, nil
			}
			m.entry = entry
			return m, nil
		}
	}
	return m, nil
}

func (m ShowEntryModel) View() string {
	s := "Show Entry\n\n"
	s += fmt.Sprintf("Kind: %s\n", m.entry.Kind)
	s += fmt.Sprintf("ID: %s\n", m.entry.ID)

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

	s += "\nPress 'Backspace' to delete this entry\n"
	s += "\nPress 'Esc' to go back, 'Ctrl+C' to quit\n"
	return s
}
