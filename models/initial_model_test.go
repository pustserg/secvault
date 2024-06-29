package models

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pustserg/secvault/config"
	"github.com/pustserg/secvault/repository"
)

var (
	cfg  *config.AppConfig
	repo repository.RepositoryInterface
)

func init() {
	cfg = &config.AppConfig{
		StoragePath: "/tmp",
	}
	repo = repository.NewRepository(cfg.StoragePath)
}

func TestNewInitialModel(t *testing.T) {
	initialModel := NewInitialModel(cfg, repo)

	if initialModel.cfg != cfg {
		t.Errorf("Expected %v, got %v", cfg, initialModel.cfg)
	}

	if initialModel.repo != repo {
		t.Errorf("Expected %v, got %v", repo, initialModel.repo)
	}

	if initialModel.choices == nil {
		t.Errorf("Expected %v, got %v", choices, initialModel.choices)
	}

	if initialModel.cursor != 0 {
		t.Errorf("Expected %v, got %v", 0, initialModel.cursor)
	}
}

func TestInitialModelUpdateMovements(t *testing.T) {
	initialModel := NewInitialModel(cfg, repo)

	downMessages := []tea.KeyMsg{
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
		tea.KeyMsg{Type: tea.KeyDown},
		tea.KeyMsg{Type: tea.KeyTab},
	}
	for _, msg := range downMessages {
		initialModel.cursor = 0
		// move cursor down, cursor = 1
		updatedModel, _ := initialModel.Update(msg)
		updatedInitialModel, ok := updatedModel.(InitialModel)
		if !ok {
			t.Errorf("After message %v Expected %v, got %v", msg, "InitialModel", updatedModel)
		}

		if updatedInitialModel.cursor != 1 {
			t.Errorf("Expected %v, got %v", 1, updatedInitialModel.cursor)
		}
		// move cursor down, cursor = 2
		updatedModel, _ = updatedInitialModel.Update(msg)
		updatedInitialModel, ok = updatedModel.(InitialModel)
		if updatedInitialModel.cursor != 2 {
			t.Errorf("Expected %v, got %v", 2, updatedInitialModel.cursor)
		}

		// move cursor down, cursor = 3 (0 in circular)
		updatedModel, _ = updatedInitialModel.Update(msg)
		updatedInitialModel, ok = updatedModel.(InitialModel)
		if updatedInitialModel.cursor != 0 {
			t.Errorf("Expected %v, got %v", 0, updatedInitialModel.cursor)
		}

		// move cursor down, cursor = 1 (circular)
		updatedModel, _ = updatedInitialModel.Update(msg)
		updatedInitialModel, ok = updatedModel.(InitialModel)

		if updatedInitialModel.cursor != 1 {
			t.Errorf("Expected %v, got %v", 1, updatedInitialModel.cursor)
		}
	}

	upMessages := []tea.KeyMsg{
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
		tea.KeyMsg{Type: tea.KeyUp},
		tea.KeyMsg{Type: tea.KeyShiftTab},
	}
	for _, msg := range upMessages {
		initialModel.cursor = 0
		// move cursor down, cursor = 2
		updatedModel, _ := initialModel.Update(msg)
		updatedInitialModel, ok := updatedModel.(InitialModel)
		if !ok {
			t.Errorf("After message %v Expected %v, got %v", msg, "InitialModel", updatedModel)
		}

		if updatedInitialModel.cursor != 2 {
			t.Errorf("Expected %v, got %v", 2, updatedInitialModel.cursor)
		}
		// move cursor up, cursor = 1
		updatedModel, _ = updatedInitialModel.Update(msg)
		updatedInitialModel, ok = updatedModel.(InitialModel)
		if updatedInitialModel.cursor != 1 {
			t.Errorf("Expected %v, got %v", 1, updatedInitialModel.cursor)
		}

		// move cursor up, cursor = 0
		updatedModel, _ = updatedInitialModel.Update(msg)
		updatedInitialModel, ok = updatedModel.(InitialModel)
		if updatedInitialModel.cursor != 0 {
			t.Errorf("Expected %v, got %v", 0, updatedInitialModel.cursor)
		}

		// move cursor up, cursor = 0 (circular)
		updatedModel, _ = updatedInitialModel.Update(msg)
		updatedInitialModel, ok = updatedModel.(InitialModel)

		if updatedInitialModel.cursor != 2 {
			t.Errorf("Expected %v, got %v", 2, updatedInitialModel.cursor)
		}
	}
}

func TestInitialModelUpdateQuit(t *testing.T) {
	initialModel := NewInitialModel(cfg, repo)

	// quit
	_, cmd := initialModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Errorf("Expected %v, got nil", tea.Quit())
	}
	cmdMsg := cmd()
	if _, ok := cmdMsg.(tea.QuitMsg); !ok {
		t.Errorf("expected command to be tea.QuitMsg, but got %T", cmdMsg)
	}
}

func TestInitialModelUpdateChoce(t *testing.T) {
	initialModel := NewInitialModel(cfg, repo)

	chociemsgs := []tea.KeyMsg{
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}},
		tea.KeyMsg{Type: tea.KeyEnter},
	}

	for _, msg := range chociemsgs {
		// generate password is the first choice
		initialModel.cursor = 0
		updatedModel, _ := initialModel.Update(msg)
		_, ok := updatedModel.(GeneratePasswordModel)
		if !ok {
			t.Errorf("After message %v Expected %v, got %v", msg, "GeneratePasswordModel", updatedModel)
		}

		// add entry is the second choice
		initialModel.cursor = 1
		updatedModel, _ = initialModel.Update(msg)
		nextModel, ok := updatedModel.(AskPasswordModel)
		if !ok {
			t.Errorf("After message %v Expected %v, got %v", msg, "AskPasswordModel", updatedModel)
		}
		if nextModel.targetAction != AddEntryTargetAction {
			t.Errorf("Expected %v, got %v", AddEntryTargetAction, nextModel.targetAction)
		}

		// list entries is the third choice
		initialModel.cursor = 2
		updatedModel, _ = initialModel.Update(msg)
		nextModel, ok = updatedModel.(AskPasswordModel)
		if !ok {
			t.Errorf("After message %v Expected %v, got %v", msg, "AskPasswordModel", updatedModel)
		}
		if nextModel.targetAction != ListEntriesTargetAction {
			t.Errorf("Expected %v, got %v", ListEntriesTargetAction, nextModel.targetAction)
		}
	}
}
