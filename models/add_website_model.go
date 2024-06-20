package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/pustserg/secvault/repository"
)

type AddWebsiteModel struct {
	repo      repository.RepositoryInterface
	password  string
	prevModel tea.Model
	cursor    int
	fields    []textinput.Model
}

func NewAddWebsiteModel(prevModel tea.Model, repo repository.RepositoryInterface, password string) AddWebsiteModel {
	nameField := textinput.New()
	nameField.Placeholder = "Name"
	nameField.Prompt = "Enter name: "
	nameField.Focus()

	urlField := textinput.New()
	urlField.Placeholder = "URL"
	urlField.Prompt = "Enter URL: "

	userNameField := textinput.New()
	userNameField.Placeholder = "User name"
	userNameField.Prompt = "Enter user name: "

	passwordField := textinput.New()
	passwordField.Placeholder = "Password"
	passwordField.Prompt = "Enter password: "

	totpTokenField := textinput.New()
	totpTokenField.Placeholder = "TOTP token"
	totpTokenField.Prompt = "Enter TOTP token: "

	return AddWebsiteModel{
		prevModel: prevModel,
		repo:      repo,
		password:  password,
		fields:    []textinput.Model{nameField, urlField, userNameField, passwordField, totpTokenField},
	}
}

func (m AddWebsiteModel) Init() tea.Cmd {
	return m.fields[0].Focus()
}

func (m AddWebsiteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	totalLen := len(m.fields) + 1
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "shift+tab":
			if m.cursor < len(m.fields) {
				m.fields[m.cursor].Blur()
			}
			m.cursor = (m.cursor - 1 + totalLen) % totalLen
			if m.cursor < len(m.fields) {
				m.fields[m.cursor].Focus()
			}
		case "down", "tab":
			if m.cursor < len(m.fields) {
				m.fields[m.cursor].Blur()
			}
			m.cursor = (m.cursor + 1) % totalLen
			if m.cursor < len(m.fields) {
				m.fields[m.cursor].Focus()
			}
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.cursor < len(m.fields)-1 {
				m.fields[m.cursor].Blur()
				m.cursor++
				m.fields[m.cursor].Focus()
			} else {
				entry := repository.Entry{
					ID:        uuid.New().String(),
					Kind:      "Website",
					Name:      m.fields[0].Value(),
					URL:       m.fields[1].Value(),
					UserName:  m.fields[2].Value(),
					Password:  m.fields[3].Value(),
					TotpToken: m.fields[4].Value(),
				}
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

func (m AddWebsiteModel) View() string {
	s := "Add website\n\n"

	cursor := " "
	for i, field := range m.fields {
		if i == m.cursor {
			cursor = ">"
		} else {
			cursor = " "
		}
		s += fmt.Sprintf("%s %s\n", cursor, field.View())
		if i < len(m.fields)-1 {
			s += "\n"
		}
	}
	if m.cursor == len(m.fields) {
		cursor = ">"
	} else {
		cursor = " "
	}
	s += fmt.Sprintf("\n%s Submit\n", cursor)

	return s
}
