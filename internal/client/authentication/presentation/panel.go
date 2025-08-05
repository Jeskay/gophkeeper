package presentation

import (
	"fmt"
	"gophkeeper/internal/client/tui"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle
	errorStyle   = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Foreground(lipgloss.Color("161"))
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func NewPanel(onComplete func(login, password string) error) Panel {
	sp := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("170"))),
	)
	m := Panel{
		inputs:     make([]textinput.Model, 2),
		onComplete: onComplete,
		spinner:    sp,
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
	spinner     spinner.Model
	cursorMode  cursor.Mode
	onComplete  func(string, string) error
	parentModel tea.Model
	loading     bool
	err         error
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

			if s == "enter" && p.focusIndex == len(p.inputs) {
				login := p.inputs[0].Value()
				password := p.inputs[1].Value()
				return p, submit(login, password, p.onComplete)
			}

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
		return p, p.spinner.Tick
	case tui.ClearErrorMsg:
		p.err = nil
	case tui.ErrorMsg:
		p.err = msg.Err
		return p, tui.ClearErrorAfter(5 * time.Second)
	case tui.LoadMsg:
		p.loading = msg.InProcess
	case successMsg:
		m, cmd := p.parentModel.Update(tui.AuthMsg{})
		return m, tea.Batch(tea.Printf("Login: %s", msg.login), cmd)
	}

	var cmdSpin tea.Cmd
	p.spinner, cmdSpin = p.spinner.Update(msg)
	cmd := p.updateInputs(msg)

	return p, tea.Batch(cmd, cmdSpin)
}

func (p Panel) View() string {
	var s strings.Builder

	if p.err != nil {
		s.WriteString(errorStyle.Render(p.err.Error()) + "\n\n")
	}

	for i := range p.inputs {
		s.WriteString(p.inputs[i].View())
		if i < len(p.inputs)-1 {
			s.WriteRune('\n')
		}
	}

	button := &blurredButton
	if p.focusIndex == len(p.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&s, "\n\n%s\n\n", *button)

	if p.loading {
		s.WriteString(fmt.Sprintf("\n\n %s Attempting to Log in...", p.spinner.View()))
	}

	s.WriteString(helpStyle.Render("cursor mode is "))
	s.WriteString(cursorModeHelpStyle.Render(p.cursorMode.String()))
	s.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return s.String()
}
