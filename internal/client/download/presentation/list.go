package presentation

import (
	"gophkeeper/internal/client/tui"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func NewList(download func(name string) error, loadList func() ([]fileInfo, error)) List {
	columns := []table.Column{
		{Title: "Id", Width: 4},
		{Title: "Title", Width: 10},
		{Title: "Size", Width: 10},
		{Title: "Status", Width: 10},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("170")).
		Bold(false)
	t.SetStyles(s)
	return List{table: t, downloadFunc: download, loadListFunc: loadList}
}

type List struct {
	table        table.Model
	parentModel  tea.Model
	downloadFunc func(string) error //TODO: make error handling for those functions
	loadListFunc func() ([]fileInfo, error)
}

func (l List) Init() tea.Cmd { return nil }

func (l List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			l.updateRows()
		case "esc":
			if l.table.Focused() {
				l.table.Blur()
			} else {
				l.table.Focus()
			}
		case "q", "ctrl+c":
			parent, cmd := l.parentModel.Update(nil)
			return parent, cmd
		case "enter":
			l.downloadFunc(l.table.SelectedRow()[1])
			parent, cmd := l.parentModel.Update(nil)
			return parent, tea.Batch(
				tea.Printf("Selected %s to download", l.table.SelectedRow()[1]),
				cmd,
			)
		}
	case tui.SpawnMsg:
		l.parentModel = msg.Parent
		l.updateRows()
		return l, tea.ClearScreen
	}
	l.table, cmd = l.table.Update(msg)
	return l, cmd
}

func (l List) View() string {
	return baseStyle.Render(l.table.View()) + "\n"
}

func (l *List) updateRows() {
	fInfos, err := l.loadListFunc()
	if err != nil {
		return
	}
	rows := make([]table.Row, len(fInfos))
	for i, n := range fInfos {
		id := strconv.FormatInt(n.Id, 10)
		rows[i] = table.Row{id, n.Name, n.Size + "B", n.Status}
	}
	l.table.SetRows(rows)
}
