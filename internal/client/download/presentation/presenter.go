package presentation

import (
	"gophkeeper/internal/client/download/control"

	tea "github.com/charmbracelet/bubbletea"
)

type fileInfo struct {
	Id     int64
	Name   string
	Status string
	Size   string
}
type downloadPresenter struct {
	controller control.Controller
}

func NewPresenter(controller control.Controller) *downloadPresenter {
	return &downloadPresenter{controller: controller}
}

func (p *downloadPresenter) CreateTable() tea.Model {
	loadFunc := func() ([]fileInfo, error) {
		list, err := p.controller.DownloadFileList()
		if err != nil {
			return []fileInfo{}, err
		}
		names := make([]fileInfo, len(list))
		for i, n := range list {
			names[i] = fileInfo{Name: n.Name, Status: n.Status, Id: n.Id, Size: n.Size}
		}
		return names, nil
	}
	return NewList(p.controller.DownloadFile, loadFunc)
}
