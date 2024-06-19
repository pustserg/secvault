package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/repository"
)

type AddEntryModel struct {
	repo      repository.RepositoryInterface
	password  string
	entry     repository.Entry
	prevModel tea.Model
	cursor    int
	fields    []textinput.Model
}

func NewAddEntryModel(prevModel tea.Model, repo repository.RepositoryInterface, password string) AddEntryModel {
	nameField := textinput.New()
	nameField.Prompt = ""
	nameField.Placeholder = "Name"
	nameField.Focus()

	urlField := textinput.New()
	urlField.Prompt = ""
	nameField.Placeholder = "Name"
	urlField.Placeholder = "URL"

	userNameField := textinput.New()
	userNameField.Prompt = ""
	userNameField.Placeholder = "User Name"

	passwordField := textinput.New()
	passwordField.Prompt = ""
	passwordField.Placeholder = "Password"

	totpTokenField := textinput.New()
	totpTokenField.Prompt = ""
	totpTokenField.Placeholder = "TOTP"

	noteField := textinput.New()
	noteField.Prompt = ""
	noteField.Placeholder = "Note"

	return AddEntryModel{
		repo:      repo,
		prevModel: prevModel,
		password:  password,
		fields:    []textinput.Model{nameField, urlField, userNameField, passwordField, totpTokenField, noteField},
	}
}

func (m AddEntryModel) Init() tea.Cmd {
	return m.fields[0].Focus()
}

func (m AddEntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				if m.cursor < len(m.fields) {
					m.fields[m.cursor].Focus()
				}
				if m.cursor < len(m.fields) {
					m.fields[m.cursor].Blur()
				}
				m.cursor--
				if m.cursor < len(m.fields) {
					m.fields[m.cursor].Focus()
				}
			}
		case "down", "enter":
			if m.cursor < len(m.fields)-1 {
				if m.cursor < len(m.fields) {
					m.fields[m.cursor].Blur()
				}
				m.cursor++
				if m.cursor < len(m.fields) {
					m.fields[m.cursor].Focus()
				}
			}
		case "ctrl+c":
			if m.cursor < len(m.fields)-1 {
				m.fields[m.cursor].Blur()
			}
			return m, tea.Quit
		default:
			m.fields[m.cursor], cmd = m.fields[m.cursor].Update(msg)
		}
	}
	return m, cmd
}

func (m AddEntryModel) View() string {
	s := "Add Entry\n\n"
	s += fmt.Sprintf("Current model: %T\n", m)
	s += "Current Password: " + m.password + "\n"
	for i, field := range m.fields {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
			field.Focus()
		}
		s += fmt.Sprintf("%s %s\n", cursor, field.View())
	}

	s += "\n"
	s += "Ctrl+C to quit\n\n"

	return s
}
