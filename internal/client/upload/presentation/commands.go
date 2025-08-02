package presentation

import (
	"gophkeeper/internal/client/tui"

	tea "github.com/charmbracelet/bubbletea"
)

type successMsg struct {
	fileName string
}

func send(path string, onComplete func(string) error) tea.Cmd {
	cmd := func() tea.Msg {
		err := onComplete(path)
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		return successMsg{fileName: path}
	}
	return tea.Sequence(tui.SetStatus(true), cmd, tui.SetStatus(false))
}
