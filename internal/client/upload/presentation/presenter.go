package presentation

import (
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"

	"gophkeeper/internal/client/upload/control"
)

type uploadPresenter struct {
	controller control.Controller
}

func NewPresenter(controller control.Controller) *uploadPresenter {
	return &uploadPresenter{controller: controller}
}

func (p *uploadPresenter) CreateExplorer() tea.Model {
	dir, err := os.UserHomeDir()
	if err != nil {
		dir = ""
	}
	dir = path.Join(dir, "Downloads")
	return NewFilePicker(dir, p.controller.UploadFile)
}
