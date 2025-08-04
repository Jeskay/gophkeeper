package abstraction

import tea "github.com/charmbracelet/bubbletea"

type Repository interface {
	SaveAvailability(value map[string]bool)
	GetAvailability() map[string]bool
	SaveToken(value string)
	GetToken() string
}

type ActionList map[string]tea.Model
