package presentation

import tea "github.com/charmbracelet/bubbletea"

type Presentation interface {
	CreateMenu(actions map[string]tea.Model) tea.Model
}

type MenuAction struct {
	Name         string
	AuthRequired bool
	Model        tea.Model
}
