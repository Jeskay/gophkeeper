package ui

import (
	"fmt"
	"gophkeeper/internal/client/tui"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func NewPanel(onComplete func(login, password string)) Panel {
	m := Panel{
		inputs:     make([]textinput.Model, 2),
		onComplete: onComplete,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		t.Width = 20

		switch i {
		case 0:
			t.Placeholder = "Login"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

type Panel struct {
	focusIndex  int
	inputs      []textinput.Model
	cursorMode  cursor.Mode
	onComplete  func(string, string)
	parentModel tea.Model
}

func (p Panel) Init() tea.Cmd { return textinput.Blink }

func (p Panel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			pModel, cmd := p.parentModel.Update(nil)
			return pModel, cmd

		// Change cursor mode
		case "ctrl+r":
			p.cursorMode++
			if p.cursorMode > cursor.CursorHide {
				p.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(p.inputs))
			for i := range p.inputs {
				cmds[i] = p.inputs[i].Cursor.SetMode(p.cursorMode)
			}
			return p, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && p.focusIndex == len(p.inputs) {
				p.onComplete(p.inputs[0].Value(), p.inputs[1].Value())
				return p, tea.ClearScreen
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				p.focusIndex--
			} else {
				p.focusIndex++
			}

			if p.focusIndex > len(p.inputs) {
				p.focusIndex = 0
			} else if p.focusIndex < 0 {
				p.focusIndex = len(p.inputs)
			}

			cmds := make([]tea.Cmd, len(p.inputs))
			for i := 0; i <= len(p.inputs)-1; i++ {
				if i == p.focusIndex {
					// Set focused state
					cmds[i] = p.inputs[i].Focus()
					p.inputs[i].PromptStyle = focusedStyle
					p.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				p.inputs[i].Blur()
				p.inputs[i].PromptStyle = noStyle
				p.inputs[i].TextStyle = noStyle
			}

			return p, tea.Batch(cmds...)
		}
	case tui.SpawnMsg:
		p.parentModel = msg.Parent
		return p, tea.ClearScreen
	}

	// Handle character input and blinking
	cmd := p.updateInputs(msg)

	return p, cmd
}

func (p Panel) View() string {
	var b strings.Builder

	for i := range p.inputs {
		b.WriteString(p.inputs[i].View())
		if i < len(p.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if p.focusIndex == len(p.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(p.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func (p *Panel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(p.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range p.inputs {
		p.inputs[i], cmds[i] = p.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
