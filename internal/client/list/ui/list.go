package ui

import (
	"gophkeeper/internal/client/tui"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func NewList(columns []table.Column, rows []table.Row) List {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
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
	return List{table: t}
}

type List struct {
	table       table.Model
	parentModel tea.Model
}

func (l List) Init() tea.Cmd { return nil }

func (l List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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
			return l, tea.Batch(
				tea.Printf("Selected %s to download", l.table.SelectedRow()[1]),
			)
		}
	case tui.SpawnMsg:
		l.parentModel = msg.Parent
		return l, tea.ClearScreen
	}
	l.table, cmd = l.table.Update(msg)
	return l, cmd
}

func (l List) View() string {
	return baseStyle.Render(l.table.View()) + "\n"
}
