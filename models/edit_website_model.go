package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/repository"
)

type EditWebsiteModel struct {
	repo      repository.RepositoryInterface
	password  string
	prevModel tea.Model
	cursor    int
	fields    []textinput.Model
	website   repository.Entry
}

func NewEditWebsiteModel(entry repository.Entry, prevModel tea.Model, repo repository.RepositoryInterface, password string) EditWebsiteModel {
	nameField := textinput.New()
	nameField.Placeholder = "Name"
	nameField.Prompt = "Enter name: "
	nameField.SetValue(entry.Name)
	nameField.Focus()

	urlField := textinput.New()
	urlField.Placeholder = "URL"
	urlField.Prompt = "Enter URL: "
	urlField.SetValue(entry.URL)

	userNameField := textinput.New()
	userNameField.Placeholder = "Username"
	userNameField.Prompt = "Enter username: "
	userNameField.SetValue(entry.UserName)

	passwordField := textinput.New()
	passwordField.Placeholder = "Password"
	passwordField.Prompt = "Enter password: "
	passwordField.SetValue(entry.Password)

	totpTokenField := textinput.New()
	totpTokenField.Placeholder = "TOTP Token"
	totpTokenField.Prompt = "Enter TOTP token: "
	totpTokenField.SetValue(entry.TotpToken)

	return EditWebsiteModel{
		prevModel: prevModel,
		repo:      repo,
		fields:    []textinput.Model{nameField, urlField, userNameField, passwordField, totpTokenField},
		password:  password,
		website:   entry,
	}
}

func (m EditWebsiteModel) Init() tea.Cmd {
	return m.fields[0].Focus()
}

func (m EditWebsiteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.website.Name = m.fields[0].Value()
				m.website.URL = m.fields[1].Value()
				m.website.UserName = m.fields[2].Value()
				m.website.Password = m.fields[3].Value()
				m.website.TotpToken = m.fields[4].Value()
				m.repo.Update(m.website, m.password)

				return m.prevModel, UpdateEntryCmd
			}
		default:
			if m.cursor < len(m.fields) {
				m.fields[m.cursor], cmd = m.fields[m.cursor].Update(msg)
			}
		}
	}
	return m, cmd
}

func (m EditWebsiteModel) View() string {
	s := "Edit Website\n\n"
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
