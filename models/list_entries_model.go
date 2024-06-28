package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/repository"
)

type ListEntriesModel struct {
	repo        repository.RepositoryInterface
	password    string
	prevModel   tea.Model
	searchQuery textinput.Model
	cursor      int
	entries     []repository.Entry
}

func NewListEntriesModel(prevModel tea.Model, repo repository.RepositoryInterface, password string) ListEntriesModel {
	searchInput := textinput.New()
	searchInput.Placeholder = "search"
	searchInput.Prompt = "Filter entries: "
	searchInput.Focus()

	return ListEntriesModel{
		searchQuery: searchInput,
		prevModel:   prevModel,
		repo:        repo,
		password:    password,
		entries:     repo.List("", password),
	}
}

func (m ListEntriesModel) Init() tea.Cmd {
	m.entries = m.repo.List(m.searchQuery.Value(), m.password)
	return nil
}

func (m ListEntriesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp, tea.KeyShiftTab:
			if len(m.entries) > 0 {
				m.cursor = (m.cursor - 1 + len(m.entries)) % len(m.entries)
			}
		case tea.KeyDown, tea.KeyTab:
			if len(m.entries) > 0 {
				m.cursor = (m.cursor + 1) % len(m.entries)
			}
		case tea.KeyEsc:
			return m.prevModel, nil
		case tea.KeyCtrlC:
			m.searchQuery.Blur()
			return m, tea.Quit
		case tea.KeyEnter:
			if len(m.entries) > 0 {
				return NewShowEntryModel(m, m.repo, m.entries[m.cursor].ID, m.password), nil
			}
		default:
			m.searchQuery, cmd = m.searchQuery.Update(msg)
			m.entries = m.repo.List(m.searchQuery.Value(), m.password)
			m.cursor = 0
			return m, cmd
		}
	case string:
		switch msg {
		case "UPDATE_ENTRIES":
			m.entries = m.repo.List(m.searchQuery.Value(), m.password)
			return m, nil
		}
	}
	return m, cmd
}

func (m ListEntriesModel) View() string {
	var cursor string

	s := "List entries\n\n"

	s += fmt.Sprintf("%s\n\n", m.searchQuery.View())

	for i, entry := range m.entries {
		if i == m.cursor {
			cursor = ">"
		} else {
			cursor = " "
		}

		var entryName string
		if entry.Name == "" {
			entryName = "Unnamed"
		} else {
			entryName = entry.Name
		}
		s += fmt.Sprintf("%s %s\n", cursor, entryName)
	}

	s += "\n\nPress 'Esc' to go back, 'Ctrl+c' to exit\n"
	return s
}
