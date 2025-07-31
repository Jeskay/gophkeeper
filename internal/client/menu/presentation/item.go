package presentation

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	id        int
	text      string
	available bool
}

func NewMenuItem(id int, text string, available bool) item {
	return item{id: id, text: text, available: available}
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.text)

	fn := conditionItemColor(itemStyle, i.available, false).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return conditionItemColor(selectedItemStyle, i.available, true).Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))

}
