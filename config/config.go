package config

import (
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

func NewAppConfig(configPath string) (*AppConfig, error) {
	cfg := AppConfig{}
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
