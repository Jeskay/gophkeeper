package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type SpawnMsg struct {
	Parent tea.Model
}

type AuthMsg struct{}

type ClearErrorMsg struct{}
type ErrorMsg struct {
	Err error
}
type LoadMsg struct {
	InProcess bool
}

func ClearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return ClearErrorMsg{}
	})
}

func SetStatus(loading bool) tea.Cmd {
	return func() tea.Msg {
		return LoadMsg{InProcess: loading}
	}
}
