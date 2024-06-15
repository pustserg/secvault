package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var DefaultConfigString = `---
password_length: 12
storage_path: ` + os.Getenv("HOME") + `/.secvault/storage.bin
`

type AppConfig struct {
	PasswordLength int    `yaml:"password_length"`
	StoragePath    string `yaml:"storage_path"`
}

func NewAppConfig(configPath string) *AppConfig {
	cfg := AppConfig{}
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	return &cfg
}
