package presentation

import (
	"gophkeeper/internal/client/download/control"

	tea "github.com/charmbracelet/bubbletea"
)

type downloadPresenter struct {
	controller control.Controller
}

func NewPresenter(controller control.Controller) *downloadPresenter {
	return &downloadPresenter{controller: controller}
}

func (p *downloadPresenter) CreateTable() tea.Model {
	loadFunc := func() []string {
		list, err := p.controller.DownloadFileList()
		if err != nil {
			return []string{err.Error()}
		}
		names := make([]string, len(list))
		for i, n := range list {
			names[i] = n.Name
		}
		return names
	}
	return NewList(p.controller.DownloadFile, loadFunc)
}
