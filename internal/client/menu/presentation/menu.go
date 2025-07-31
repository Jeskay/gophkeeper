package presentation

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

func NewActionMenu(actions []*MenuAction, checkAuth func() bool) Menu {
	var items = make([]list.Item, len(actions))
	for i, a := range actions {
		items[i] = NewMenuItem(i, a.Name, !a.AuthRequired)
	}
	const defaultWidth = 20
	const listHeight = 14
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Gophkeeper UI"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return Menu{list: l, actions: actions, checkAuth: checkAuth}
}

type Menu struct {
	list      list.Model
	choiceId  int
	quitting  bool
	actions   []*MenuAction
	checkAuth func() bool
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
	case tui.AuthMsg:
		if m.checkAuth() {
			for i, v := range m.list.Items() {
				itemCopy := v.(item)
				if m.actions[itemCopy.id].AuthRequired {
					itemCopy.available = true
				} else {
					itemCopy.available = false
				}
				m.list.SetItem(i, itemCopy)
			}
		}
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choiceId = i.id
			}
			if i.available {
				nxtM := m.actions[m.choiceId]
				newModel, cmd := nxtM.Model.Update(tui.SpawnMsg{Parent: m})
				return newModel, cmd
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
