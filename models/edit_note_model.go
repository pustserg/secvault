package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/repository"
)

type EditNoteModel struct {
	repo      repository.RepositoryInterface
	password  string
	prevModel tea.Model
	cursor    int
	fields    []textinput.Model
	note      repository.Entry
}

func NewEditNoteModel(note repository.Entry, prevModel tea.Model, repo repository.RepositoryInterface, password string) EditNoteModel {
	nameField := textinput.New()
	nameField.Placeholder = "Name"
	nameField.Prompt = "Enter name: "
	nameField.SetValue(note.Name)
	nameField.Focus()

	noteField := textinput.New()
	noteField.Placeholder = "Note"
	noteField.SetValue(note.Note)
	noteField.Prompt = "Enter note: "

	return EditNoteModel{
		prevModel: prevModel,
		repo:      repo,
		fields:    []textinput.Model{nameField, noteField},
		password:  password,
		note:      note,
	}
}

func (m EditNoteModel) Init() tea.Cmd {
	return m.fields[0].Focus()
}

func (m EditNoteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "shift+tab":
			if m.cursor > 0 {
				if m.cursor < len(m.fields) {
					m.fields[m.cursor].Blur()
				}
				m.cursor--

				m.fields[m.cursor].Focus()
			}
		case "down", "tab":
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
			} else {
				m.note.Name = m.fields[0].Value()
				m.note.Note = m.fields[1].Value()
				m.repo.Update(m.note, m.password)
				var msg tea.Msg = "UPDATE_ENTRY"
				callbackCommand := func() tea.Msg { return msg }

				return m.prevModel, callbackCommand
			}
		default:
			if m.cursor < len(m.fields) {
				m.fields[m.cursor], cmd = m.fields[m.cursor].Update(msg)
			}
		}
	}
	return m, cmd
}

func (m EditNoteModel) View() string {
	s := "Edit note\n\n"
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
