package config

import (
	"os"
	"testing"
)

func TestNewAppConfigExitCode(t *testing.T) {
	// with invalid config path
	path := "invalid_path"
	cfg, err := NewAppConfig(path)

	if err == nil {
		t.Fatalf("Expected error with path %s, but no error returned", path)
	}

	if cfg != nil {
		t.Fatalf("Expected nil config with path %s, but got %v", path, cfg)
	}

	// with valid config path
	path = "./config.yaml"
	configContent := "---\npassword_length: 12\nstorage_path: /tmp\n"
	tmpFile, err := os.Create(path)

	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Error writing to temp file: %v", err)
	}

	cfg, err = NewAppConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error reading config: %v", err)
	}

	if cfg.PasswordLength != 12 {
		t.Fatalf("Expected password length 12, but got %d", cfg.PasswordLength)
	}

	if cfg.StoragePath != "/tmp" {
		t.Fatalf("Expected storage path /tmp, but got %s", cfg.StoragePath)
	}
}
