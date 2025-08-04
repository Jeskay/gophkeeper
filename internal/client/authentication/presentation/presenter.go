package presentation

import (
	"gophkeeper/internal/client/authentication/control"

	tea "github.com/charmbracelet/bubbletea"
)

type authPresenter struct {
	controller control.Controller
}

func NewPresenter(controller control.Controller) *authPresenter {
	return &authPresenter{controller: controller}
}

func (p *authPresenter) CreateLoginPanel(onComplete func(string)) tea.Model {
	f := func(name, password string) error {
		token, err := p.controller.LogIn(name, password)
		if err == nil {
			onComplete(token)
		}
		return err
	}
	return NewPanel(f)
}

func (p *authPresenter) CreateRegisterPanel() tea.Model {
	f := func(name, password string) error {
		return p.controller.Register(name, password)
	}
	return NewPanel(f)
}
