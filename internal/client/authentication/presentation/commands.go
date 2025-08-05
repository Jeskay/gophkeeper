package presentation

import (
	"gophkeeper/internal/client/tui"

	tea "github.com/charmbracelet/bubbletea"
)

type successMsg struct {
	login string
}

func (p *Panel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(p.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range p.inputs {
		p.inputs[i], cmds[i] = p.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func submit(login string, password string, onComplete func(string, string) error) tea.Cmd {
	cmd := func() tea.Msg {
		err := onComplete(login, password)
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		return successMsg{login: login}
	}
	return tea.Sequence(tui.SetStatus(true), cmd, tui.SetStatus(false))
}
