package presentation

import (
	"errors"
	"gophkeeper/internal/client/tui"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type successMsg struct {
	fileName string
}

type fetchedFiles struct {
	rows []table.Row
}

func fetchFiles(fetch func() ([]fileInfo, error)) tea.Cmd {
	cmd := func() tea.Msg {
		fInfos, err := fetch()
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		rows := make([]table.Row, len(fInfos))
		for i, n := range fInfos {
			id := strconv.FormatInt(n.Id, 10)
			rows[i] = table.Row{id, n.Name, n.Size + "B", n.Status}
		}
		return fetchedFiles{rows: rows}
	}
	return tea.Sequence(tui.SetStatus(true), cmd, tui.SetStatus(false))
}

func back(parent tea.Model) (tea.Model, tea.Cmd) {
	return parent.Update(nil)
}

func downloadFile(download func(int64) error, row table.Row) tea.Cmd {
	cmd := func() tea.Msg {
		id, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			return tui.ErrorMsg{Err: errors.New("failed to parse file id")}
		}
		if err = download(id); err != nil {
			return tui.ErrorMsg{Err: err}
		}
		return successMsg{fileName: row[1]}
	}
	return tea.Sequence(tui.SetStatus(true), cmd, tui.SetStatus(false))
}
