package presentation

import (
	tea "github.com/charmbracelet/bubbletea"

	authPresent "gophkeeper/internal/client/authentication/presentation"
	downloadPresent "gophkeeper/internal/client/download/presentation"
	"gophkeeper/internal/client/menu/control"
	uploadPresent "gophkeeper/internal/client/upload/presentation"
)

type menuPresenter struct {
	controller control.Controller
	auth       authPresent.Presentation
	upload     uploadPresent.Presentation
	download   downloadPresent.Presentation
}

func NewPresenter(controller control.Controller, auth authPresent.Presentation, upload uploadPresent.Presentation, download downloadPresent.Presentation) *menuPresenter {
	return &menuPresenter{controller: controller, auth: auth, upload: upload, download: download}
}

func (p *menuPresenter) CreateMenu() tea.Model {
	actions := []*MenuAction{
		{Name: "Log In", AuthRequired: false, Model: p.auth.CreateLoginPanel(func(s1 string) { p.controller.SaveToken(s1) })},
		{Name: "Register", AuthRequired: false, Model: p.auth.CreateRegisterPanel()},
		{Name: "Download", AuthRequired: true, Model: p.download.CreateTable()},
		{Name: "Upload", AuthRequired: true, Model: p.upload.CreateExplorer()},
	}
	return NewActionMenu(actions, p.controller.IsAuthenticated)
}
