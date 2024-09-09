package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/config"
	"github.com/pustserg/secvault/models"
	"github.com/pustserg/secvault/repository"
)

var (
	repo       repository.RepositoryInterface
	workdir    = os.Getenv("HOME") + "/.secvault"
	configPath = workdir + "/config.yaml"
)

func main() {
	ensureAppDirExists(workdir)
	ensureConfigFileExists(configPath)

	cfg, err := config.NewAppConfig(configPath)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	ensureStoragePathExists(cfg.StoragePath)
	repo = repository.NewRepository(cfg.StoragePath)

	model := models.NewInitialModel(cfg, repo)

	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		fmt.Println("Error starting program:", err)
	}
}

func ensureAppDirExists(workdir string) {
	if _, err := os.Stat(workdir); os.IsNotExist(err) {
		err := os.Mkdir(workdir, 0755)
		if err != nil {
			fmt.Println("Error creating app dir:", err)
			os.Exit(1)
		}
	} else if err != nil {
		fmt.Println("Error checking app dir:", err)
		os.Exit(1)
	}
}

func ensureConfigFileExists(configFilePath string) {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fmt.Println("Config file not found, would you like to create a new one? (y/n)")
		var answer string
		_, err := fmt.Scanln(&answer)
		if err != nil {
			fmt.Println("Error reading user input:", err)
			os.Exit(1)
		}
		if answer == "y" {
			saveDefaultConfig(configFilePath)
		} else {
			os.Exit(0)
		}
	}
}

func ensureStoragePathExists(storagePath string) {
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		fmt.Println("Storage file not found, would you like to create a new one? (y/n)")
		var answer string
		_, err := fmt.Scanln(&answer)
		if err != nil {
			fmt.Println("Error reading user input:", err)
			os.Exit(1)
		}
		if answer == "y" {
			fmt.Println("Creating storage file at", storagePath)
			file, err := os.OpenFile(storagePath, os.O_RDONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println("Error creating storage file:", err)
				os.Exit(1)
			}
			err = file.Close()
			if err != nil {
				fmt.Println("Error closing storage file:", err)
				os.Exit(1)
			}
		} else {
			os.Exit(0)
		}
	} else if err != nil {
		fmt.Println("Error checking storage path:", err)
		os.Exit(1)
	}
}

func saveDefaultConfig(configFilePath string) {
	file, err := os.Create(configFilePath)
	if err != nil {
		fmt.Println("Error creating config file:", err)
		os.Exit(1)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing config file:", err)
			os.Exit(1)
		}
	}(file)

	_, err = file.WriteString(config.DefaultConfigString)
	if err != nil {
		fmt.Println("Error writing to config file:", err)
		os.Exit(1)
	}
}
