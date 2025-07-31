package presentation

import tea "github.com/charmbracelet/bubbletea"

type Presentation interface {
	CreateLoginPanel(onComplete func(string)) tea.Model
	CreateRegisterPanel() tea.Model
}
