package presentation

import (
	"fmt"
	"gophkeeper/internal/client/tui"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
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
	lSpinner := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("170"))),
	)
	return List{
		table:        t,
		downloadFunc: download,
		loadListFunc: loadList,
		spinner:      lSpinner,
	}
}

type List struct {
	table        table.Model
	spinner      spinner.Model
	parentModel  tea.Model
	downloadFunc func(int64) error
	loadListFunc func() ([]fileInfo, error)
	err          error
	loading      bool
}

func (l List) Init() tea.Cmd { return l.spinner.Tick }

func (l List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			return l, fetchFiles(l.loadListFunc)
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
			return l, downloadFile(l.downloadFunc, l.table.SelectedRow())
		}
	case tui.SpawnMsg:
		l.parentModel = msg.Parent
		return l, tea.Sequence(l.spinner.Tick, fetchFiles(l.loadListFunc))
	case tui.ErrorMsg:
		l.err = msg.Err
		return l, tui.ClearErrorAfter(5 * time.Second)
	case tui.ClearErrorMsg:
		l.err = nil
	case tui.LoadMsg:
		l.loading = msg.InProcess
	case fetchedFiles:
		l.table.SetRows(msg.rows)
	case successMsg:
		parent, cmd := back(l.parentModel)
		return parent, tea.Sequence(
			tea.Printf("File %s saved to downloads folder", msg.fileName),
			cmd,
		)
	}
	var cmdSpin tea.Cmd
	l.spinner, cmdSpin = l.spinner.Update(msg)
	l.table, cmd = l.table.Update(msg)
	return l, tea.Batch(cmd, cmdSpin)
}

func (l List) View() string {
	var s strings.Builder
	if l.err != nil {
		s.WriteString(errorStyle.Render(l.err.Error()))
	}
	s.WriteString(baseStyle.Render("\n\n"+l.table.View()) + "\n")
	if l.loading {
		s.WriteString(fmt.Sprintf("\n\n %s Downloading file...", l.spinner.View()))
	}
	return s.String()
}
