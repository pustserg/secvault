package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/repository"
)

type AddNoteModel struct {
	repo      repository.RepositoryInterface
	password  string
	prevModel tea.Model
	cursor    int
	fields    []textinput.Model
}

func NewAddNoteModel(prevModel tea.Model, repo repository.RepositoryInterface, password string) AddNoteModel {
	nameField := textinput.New()
	nameField.Placeholder = "Name"
	nameField.Prompt = "Enter name: "
	nameField.Focus()

	noteField := textinput.New()
	noteField.Placeholder = "Note"
	noteField.Prompt = "Enter note: "

	return AddNoteModel{
		prevModel: prevModel,
		repo:      repo,
		fields:    []textinput.Model{nameField, noteField},
		password:  password,
	}
}

func (m AddNoteModel) Init() tea.Cmd {
	return m.fields[0].Focus()
}

func (m AddNoteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				if m.cursor < len(m.fields) {
					m.fields[m.cursor].Blur()
				}
				m.cursor--

				m.fields[m.cursor].Focus()
			}
		case "down":
			if m.cursor < len(m.fields)-1 {
				m.fields[m.cursor].Blur()
				m.cursor++

				m.fields[m.cursor].Focus()
			} else if m.cursor == len(m.fields)-1 {
				m.cursor++
			}
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.cursor < len(m.fields)-1 {
				m.fields[m.cursor].Blur()
				m.cursor++

				m.fields[m.cursor].Focus()
			} else if m.cursor == 2 {
				entry := repository.Entry{Name: m.fields[0].Value(), Note: m.fields[1].Value()}
				m.repo.Add(entry, m.password)
				return m.prevModel, nil
			}
		default:
			if m.cursor < len(m.fields) {
				m.fields[m.cursor], cmd = m.fields[m.cursor].Update(msg)
			}
		}
	}
	return m, cmd
}

func (m AddNoteModel) View() string {
	s := "Add note\n\n"
	s += fmt.Sprintf("Cursor: %d\n", m.cursor)

	var cursor string

	for i, field := range m.fields {
		if i == m.cursor {
			cursor = ">"
		} else {
			cursor = " "
		}

		s += fmt.Sprintf("%s %s", cursor, field.View())
		s += "\n"
	}

	if m.cursor == len(m.fields) {
		cursor = ">"
	} else {
		cursor = " "
	}

	s += fmt.Sprintf("%s Submit\n", cursor)

	return s
}
