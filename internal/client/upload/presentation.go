package upload

import (
	"gophkeeper/internal/client/upload/ui"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
)

func CreateExplorer() tea.Model {
	dir, err := os.UserHomeDir()
	if err != nil {
		dir = ""
	}
	dir = path.Join(dir, "Downloads")
	return ui.NewFilePicker(dir)
}
