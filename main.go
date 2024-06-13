package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/models"
)

func main() {
	model := models.NewInitialModel()
	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		fmt.Println("Error starting program:", err)
	}
}
