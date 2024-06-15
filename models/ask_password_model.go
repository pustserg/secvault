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
	targetAction string
	cursor       int
	choices      []string
}

func NewAskPasswordModel(prevModel tea.Model, repo repository.RepositoryInterface) AskPasswordModel {
	textInput := textinput.New()
	textInput.Placeholder = "password"
	textInput.Prompt = ""
	textInput.EchoMode = textinput.EchoPassword
	textInput.EchoCharacter = '*'
	textInput.Focus()
	return AskPasswordModel{
		prevModel: prevModel,
		password:  textInput,
		choices:   []string{"back", "quit"},
	}
}

func (m AskPasswordModel) Init() tea.Cmd {
	return m.password.Focus()
}

func (m AskPasswordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor == 0 {
					m.password.Focus()
				}
			}
		case "down":
			if m.cursor < 3 {
				if m.cursor == 0 {
					m.password.Blur()
				}
				m.cursor++
			}
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			return handleEnter(m, msg)
		default:
			m.password, cmd = m.password.Update(msg)
		}
	}
	return m, cmd
}

func (m AskPasswordModel) View() string {
	s := "Enter password for storage: \n\n"
	s += "Current Password: " + m.password.Value() + "\n\n"

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

	return s
}

func handleEnter(m AskPasswordModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	// if cursor on the first element, we should update password and move to the next model
	switch m.cursor {
	case 0:
		var cmd tea.Cmd
		m.password, cmd = m.password.Update(msg)
		return m, cmd
	// if coursor on the second element, we should go back to the prevModel
	case 1:
		return m.prevModel, nil
	// if cursor on the third element, we should quit the app
	case 2:
		return m, tea.Quit
	}

	return m, nil
}
