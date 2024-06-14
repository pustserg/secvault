package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/config"
	"github.com/pustserg/secvault/models"
)

var (
	cfg         *config.AppConfig
	workdir     string = os.Getenv("HOME") + "/.secvault"
	config_path string = workdir + "/config.yaml"
)

func main() {
	ensureAppDirExists(workdir)
	ensureConfigFileExists(config_path)

	cfg = config.NewAppConfig(config_path)

	model := models.NewInitialModel(cfg)
	fmt.Println("Password length:", cfg.PasswordLength)
	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		fmt.Println("Error starting program:", err)
	}
}

func ensureAppDirExists(workdir string) {
	if _, err := os.Stat(workdir); os.IsNotExist(err) {
		os.Mkdir(workdir, 0755)
	} else if err != nil {
		fmt.Println("Error checking app dir:", err)
		os.Exit(1)
	}
}

func ensureConfigFileExists(configFilePath string) {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fmt.Println("Config file not found, would you like to create a new one? (y/n)")
		var answer string
		fmt.Scanln(&answer)
		if answer == "y" {
			saveDefaultConfig(configFilePath)
		} else {
			os.Exit(0)
		}
	}
}

func saveDefaultConfig(configFilePath string) {
	file, err := os.Create(configFilePath)
	if err != nil {
		fmt.Println("Error creating config file:", err)
		os.Exit(1)
	}

	defer file.Close()

	_, err = file.WriteString(config.DefaultConfigString)
	if err != nil {
		fmt.Println("Error writing to config file:", err)
		os.Exit(1)
	}
}
