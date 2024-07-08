package models

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/config"
	"golang.design/x/clipboard"
)

const (
	maxPasswordLength = 3 // 999 is maximum chars in password value, not the length of the password
)

var (
	symbolsArray   = []rune("!@#$%^&*()_+-=[]{}|;:,.<>?")
	numbersArray   = []rune("0123456789")
	uppercaseArray = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lowercaseArray = []rune("abcdefghijklmnopqrstuvwxyz")
)

type GeneratePasswordModel struct {
	cfg       *config.AppConfig
	length    textinput.Model
	password  string
	prevModel tea.Model
	options   []string
	cursor    int
	selected  map[string]bool
}

func NewGeneratePasswordModel(prevModel tea.Model, cfg *config.AppConfig) GeneratePasswordModel {
	textInput := textinput.New()
	textInput.Placeholder = "password length (1-999)"
	textInput.CharLimit = maxPasswordLength
	textInput.SetValue(strconv.Itoa(cfg.PasswordLength))
	textInput.Focus()

	m := GeneratePasswordModel{
		cfg:       cfg,
		prevModel: prevModel,
		length:    textInput,
		options:   []string{"symbols", "numbers", "uppercase", "lowercase"},
		selected:  map[string]bool{"symbols": true, "numbers": true, "uppercase": true, "lowercase": true},
	}

	m.password = generatePassword(m.length.Value(), m.selected)

	return m
}

func (m GeneratePasswordModel) Init() tea.Cmd {
	return m.length.Focus()
}

func (m GeneratePasswordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Ctrl+c":
			m.length.Blur()
			return m, tea.Quit
		case "j", "down", "tab":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "k", "up", "shift+tab":
			if m.cursor > 0 {
				m.cursor--
			}

		case "b", "esc":
			return m.prevModel, nil
		case "c":
			err := clipboard.Init()
			if err != nil {
				panic(err)
			}

			clipboard.Write(clipboard.FmtText, []byte(m.password))
		case "r":
			m.password = generatePassword(m.length.Value(), m.selected)
		case " ", "enter":
			m.selected[m.options[m.cursor]] = !m.selected[m.options[m.cursor]]
			m.password = generatePassword(m.length.Value(), m.selected)
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "backspace":
			m.length, cmd = m.length.Update(msg)
			m.password = generatePassword(m.length.Value(), m.selected)
		}
	}
	return m, cmd
}

func (m GeneratePasswordModel) View() string {
	s := "generate password\n\n"

	s += fmt.Sprintf("length: %s\n\n", m.length.View())
	for i, option := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "

		if m.selected[option] {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, option)
	}

	s += fmt.Sprintf("\npassword: %s\n", m.password)

	s += "\npress 'r' to regenerate, 'c' to copy password to clipboard, 'q' to quit or 'b' to go back\n"
	return s
}

func generatePassword(lengthValue string, selected map[string]bool) string {
	length, err := strconv.Atoi(lengthValue)
	if err != nil {
		return ""
	}
	var chars []rune

	if selected["symbols"] {
		chars = append(chars, symbolsArray...)
	}

	if selected["numbers"] {
		chars = append(chars, numbersArray...)
	}

	if selected["uppercase"] {
		chars = append(chars, uppercaseArray...)
	}

	if selected["lowercase"] {
		chars = append(chars, lowercaseArray...)
	}

	if len(chars) == 0 {
		return ""
	}

	password := make([]rune, length)
	for i := range password {
		password[i] = chars[rand.Intn(len(chars))]
	}

	return string(password)
}
