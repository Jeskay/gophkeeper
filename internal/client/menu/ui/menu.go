package ui

import (
	"gophkeeper/internal/client/tui"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle         = lipgloss.NewStyle().MarginLeft(2)
	paginationStyle    = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle          = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle      = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	itemStyle          = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle  = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	conditionItemColor = func(s lipgloss.Style, available bool, selected bool) lipgloss.Style {
		if !available {
			return s.Foreground(lipgloss.Color("243"))
		} else if selected {
			return s.Foreground(lipgloss.Color("170"))
		}
		return s
	}
)

func NewMenu(items []list.Item, actions map[string]tea.Model) Menu {
	const defaultWidth = 20
	const listHeight = 14
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Gophkeeper UI"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return Menu{list: l, actions: actions}
}

type Menu struct {
	list     list.Model
	choice   string
	quitting bool
	actions  map[string]tea.Model
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) View() string {
	if m.quitting {
		return quitTextStyle.Render("Quitting...")
	}
	return "\n" + m.list.View()
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.text
			}
			if i.available {
				nxtM := m.actions[m.choice]
				newModel, cmd := nxtM.Update(tui.SpawnMsg{Parent: m})
				return newModel, cmd
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
