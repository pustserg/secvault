package models

import tea "github.com/charmbracelet/bubbletea"

type ConfirmationModel struct {
	message            string
	successReturnModel tea.Model
	failureReturnModel tea.Model
	callbackCommand    tea.Cmd
	choices            []string
	action             func() error
}

func NewConfirmationModel(successReturnModel, failureReturnModel tea.Model, callbackCommand tea.Cmd, message string, choices []string, action func() error) ConfirmationModel {
	return ConfirmationModel{
		message:            message,
		successReturnModel: successReturnModel,
		failureReturnModel: failureReturnModel,
		callbackCommand:    callbackCommand,
		choices:            choices,
		action:             action,
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
			return m.successReturnModel, m.callbackCommand
		case "n":
			return m.failureReturnModel, nil
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
