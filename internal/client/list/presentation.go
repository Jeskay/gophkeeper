package list

import (
	"gophkeeper/internal/client/list/ui"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func CreateList() tea.Model {
	return ui.NewList([]table.Column{
		{Title: "Id", Width: 4},
		{Title: "Title", Width: 10},
		{Title: "Path", Width: 20},
		{Title: "Size", Width: 10},
	}, []table.Row{
		{"1", "file.dat", "/home/jeskay/documents", "11Mb"},
		{"2", "notes.txt", "/home/jeskay/notes", "2Kb"},
		{"3", "image.png", "/home/jeskay/pictures", "1.2Mb"},
		{"4", "report.pdf", "/home/jeskay/reports", "800Kb"},
		{"5", "music.mp3", "/home/jeskay/music", "5Mb"},
		{"6", "video.mp4", "/home/jeskay/videos", "700Mb"},
		{"7", "archive.zip", "/home/jeskay/backups", "120Mb"},
		{"8", "todo.md", "/home/jeskay/tasks", "4Kb"},
		{"9", "slides.pptx", "/home/jeskay/presentations", "3Mb"},
	})
}
