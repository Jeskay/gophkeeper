package menu

import (
	"gophkeeper/internal/client/menu/ui"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func CreateMenu(actions map[string]tea.Model) tea.Model {
	var items = make([]list.Item, len(actions))
	var i = 0
	for n := range actions {
		items[i] = ui.NewItem(n, true)
		i++
	}
	return ui.NewMenu(items, actions)
}
