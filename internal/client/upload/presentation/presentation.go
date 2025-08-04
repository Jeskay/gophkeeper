package presentation

import tea "github.com/charmbracelet/bubbletea"

type Presentation interface {
	CreateExplorer() tea.Model
}
