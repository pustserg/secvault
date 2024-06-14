package config

import (
	"fmt"
	"os"

	"github.com/num30/config"
)

const DefaultConfigString = `---
password_length: 12
storage_path: ~/.secvault/storage.bin
`

type AppConfig struct {
	PasswordLength int    `yaml:"password_length" default:"12" validate:"min=1,max=999,required"`
	StoragePath    string `yaml:"storage_path" default:"~/.secvault/storage.bin" validate:"required"`
}

func NewAppConfig(configFilePath string) *AppConfig {
	cfg := AppConfig{}
	err := config.NewConfReader(configFilePath).Read(&cfg)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	return &cfg
}

// func ensureConfigFileExists() {
// 	ensureAppDirExists()
// 	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
// 		saveDefaultConfig()
// 	}
// }

// func ensureAppDirExists() {
// 	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
// 		os.Mkdir(configDirPath, 0755)
// 	} else if err != nil {
// 		fmt.Println("Error checking app dir:", err)
// 		os.Exit(1)
// 	}
// }

// func saveDefaultConfig() {
// 	file, err := os.Create(configFilePath)
// 	if err != nil {
// 		fmt.Println("Error creating config file:", err)
// 		os.Exit(1)
// 	}

// 	defer file.Close()

// 	_, err = file.WriteString(defaultConfigString)
// 	if err != nil {
// 		fmt.Println("Error writing to config file:", err)
// 		os.Exit(1)
// 	}
// }
