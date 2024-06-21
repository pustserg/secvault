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
	return nil
}

func (m ListEntriesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "shift+tab":
			if len(m.entries) > 0 {
				m.cursor = (m.cursor - 1 + len(m.entries)) % len(m.entries)
			}
		case "down", "tab":
			if len(m.entries) > 0 {
				m.cursor = (m.cursor + 1) % len(m.entries)
			}
		case "b":
			return m.prevModel, nil
		case "q", "ctrl+c":
			m.searchQuery.Blur()
			return m, tea.Quit
		case "enter":
			if len(m.entries) > 0 {
				return NewShowEntryModel(m, m.repo, m.entries[m.cursor], m.password), nil
			}
		default:
			m.searchQuery, cmd = m.searchQuery.Update(msg)
			m.entries = m.repo.List(m.searchQuery.Value(), m.password)
			return m, cmd
		}
	}
	return m, cmd
}

func (m ListEntriesModel) View() string {
	var cursor string

	s := "List entries\n\n"

	s += fmt.Sprintf("Total: %d\n\n", len(m.entries))

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

	s += "\n\nPress 'b' to go back, 'q' to exit\n"
	return s
}
