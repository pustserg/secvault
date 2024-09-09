package models

import tea "github.com/charmbracelet/bubbletea"

var (
	UpdateEntryMsg   tea.Msg = "UPDATE_ENTRY"
	UpdateEntriesMsg tea.Msg = "UPDATE_ENTRIES"

	UpdateEntryCmd   tea.Cmd = func() tea.Msg { return UpdateEntryMsg }
	UpdateEntriesCmd tea.Cmd = func() tea.Msg { return UpdateEntriesMsg }
)

const (
	cursorSymbol            = ">"
	ListEntriesTargetAction = "LIST_ENTRIES"
	AddEntryTargetAction    = "ADD_ENTRY"
)
