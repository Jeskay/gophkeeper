package authentication

import (
	"gophkeeper/internal/client/authentication/ui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func CreateAuthPanel() tea.Model {
	return ui.NewPanel(func(login, password string) {
		log.Printf("Login:%s\nPassword:%s", login, password)
	})
}
