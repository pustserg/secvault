package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/repository"
)

type AskPasswordModel struct {
	password     textinput.Model
	prevModel    tea.Model
	repo         repository.RepositoryInterface
	targetAction string
	cursor       int
	choices      []string
	errorMessage string
}

func NewAskPasswordModel(prevModel tea.Model, repo repository.RepositoryInterface, targetAction string) AskPasswordModel {
	textInput := newPasswordInput()

	return AskPasswordModel{
		prevModel:    prevModel,
		password:     textInput,
		repo:         repo,
		targetAction: targetAction,
		choices:      []string{"back", "quit"},
	}
}

func newPasswordInput() textinput.Model {
	textInput := textinput.New()
	textInput.Placeholder = "password"
	textInput.Prompt = ""
	textInput.EchoMode = textinput.EchoPassword
	textInput.EchoCharacter = '*'
	textInput.Focus()
	return textInput
}

func (m AskPasswordModel) Init() tea.Cmd {
	return m.password.Focus()
}

func (m AskPasswordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	totalChoices := len(m.choices) + 1
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "shift+tab":
			if m.cursor == 0 {
				m.password.Blur()
			}
			m.cursor = (m.cursor - 1 + totalChoices) % totalChoices
			if m.cursor == 0 {
				m.password.Focus()
			}
		case "down", "tab":
			if m.cursor == 0 {
				m.password.Blur()
			}
			m.cursor = (m.cursor + 1) % (len(m.choices) + 1)
			if m.cursor == 0 {
				m.password.Focus()
			}
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			// if cursor is on password field update password value and call next model
			if m.cursor == 0 {
				m.password.Blur()
				if err := m.repo.CheckPassword(m.password.Value()); err != nil {
					m.errorMessage = err.Error()
					return m, nil
				} else {
					if m.targetAction == "add entry" {
						return NewChooseEntryKindModel(m, m.repo, m.password.Value()), nil
					} else {
						return m, nil
					}
				}
			} else if m.cursor == 1 { // back
				return m.prevModel, nil
			} else if m.cursor == 2 { // quit
				return m, tea.Quit
			}
		default:
			m.password, cmd = m.password.Update(msg)
			return m, cmd
		}
	}
	return m, cmd
}

func (m AskPasswordModel) View() string {
	s := "Enter password for storage: \n\n"

	cursor := " "

	// render password input
	if m.cursor == 0 {
		cursor = ">"
	}
	s += fmt.Sprintf("%s Enter password: %s\n", cursor, m.password.View())

	// render choices
	for i, choice := range m.choices {
		if m.cursor == i+1 {
			cursor = ">"
		} else {
			cursor = " "
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	if m.errorMessage != "" {
		s += fmt.Sprintf("\nError: %s", m.errorMessage)
	}

	return s
}
