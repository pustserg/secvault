package models

import tea "github.com/charmbracelet/bubbletea"

var (
	UpdateEntryMsg   tea.Msg = "UPDATE_ENTRY"
	UpdateEntriesMsg tea.Msg = "UPDATE_ENTRIES"

	UpdateEntryCmd   tea.Cmd = func() tea.Msg { return UpdateEntryMsg }
	UpdateEntriesCmd tea.Cmd = func() tea.Msg { return UpdateEntriesMsg }
)
