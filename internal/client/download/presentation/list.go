package presentation

import (
	"gophkeeper/internal/client/tui"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var errorStyle = baseStyle.Foreground(lipgloss.Color("161"))

func NewList(download func(id int64) error, loadList func() ([]fileInfo, error)) List {
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
	downloadFunc func(int64) error
	loadListFunc func() ([]fileInfo, error)
	err          error
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
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
			selectedId, err := strconv.ParseInt(l.table.SelectedRow()[0], 10, 64)
			if err != nil {
				l.err = err
				return l, clearErrorAfter(5 * time.Second)
			}
			l.err = l.downloadFunc(selectedId)
			if l.err != nil {
				return l, clearErrorAfter(5 * time.Second)
			}
			parent, cmd := l.parentModel.Update(nil)
			return parent, tea.Sequence(
				tea.Printf("File %s saved to downloads folder", l.table.SelectedRow()[1]),
				cmd,
			)
		}
	case tui.SpawnMsg:
		l.parentModel = msg.Parent
		l.err = l.updateRows()
		if l.err != nil {
			return l, clearErrorAfter(5 * time.Second)
		}
		return l, tea.ClearScreen
	case clearErrorMsg:
		l.err = nil
	}
	l.table, cmd = l.table.Update(msg)
	return l, cmd
}

func (l List) View() string {
	var s strings.Builder
	if l.err != nil {
		s.WriteString(errorStyle.Render(l.err.Error()))
	}
	s.WriteString(baseStyle.Render("\n\n"+l.table.View()) + "\n")
	return s.String()
}

func (l *List) updateRows() error {
	fInfos, err := l.loadListFunc()
	if err != nil {
		return err
	}
	rows := make([]table.Row, len(fInfos))
	for i, n := range fInfos {
		id := strconv.FormatInt(n.Id, 10)
		rows[i] = table.Row{id, n.Name, n.Size + "B", n.Status}
	}
	l.table.SetRows(rows)
	return nil
}
