package models

import tea "github.com/charmbracelet/bubbletea"

type ConfirmationModel struct {
	message         string
	prevModel       tea.Model
	callbackCommand tea.Cmd
	choices         []string
	action          func() error
}

func NewConfirmationModel(prevModel tea.Model, callbackCommand tea.Cmd, message string, choices []string, action func() error) ConfirmationModel {
	return ConfirmationModel{
		message:         message,
		prevModel:       prevModel,
		callbackCommand: callbackCommand,
		choices:         choices,
		action:          action,
	}
}

func (m ConfirmationModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			// err := m.action()
			// if err != nil {
			// 	return NewErrorModel(m, err.Error()), nil
			// }
			m.action()
			return m.prevModel, m.callbackCommand
		case "n":
			return m.prevModel, nil
		}
	}
	return m, cmd
}

func (m ConfirmationModel) View() string {
	s := m.message + "\n\n"

	for i, choice := range m.choices {
		s += choice
		if i < len(m.choices)-1 {
			s += " / "
		}
	}

	return s
}
