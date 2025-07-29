package client

import (
	"gophkeeper/internal/client/authentication"
	"gophkeeper/internal/client/list"
	"gophkeeper/internal/client/menu"
	"gophkeeper/internal/client/upload"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start() {
	m := menu.CreateMenu(map[string]tea.Model{
		"Log In":   authentication.CreateAuthPanel(),
		"Download": list.CreateList(),
		"Upload":   upload.CreateExplorer(),
	})
	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal("unexpected error", err)
	}
}
