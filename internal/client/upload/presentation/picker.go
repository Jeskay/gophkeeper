package presentation

import (
	"errors"
	"fmt"
	"gophkeeper/internal/client/tui"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewFilePicker(dir string, onComplete func(file string) error) *FilePicker {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".md", ".txt", ".dat"}
	fp.CurrentDirectory = dir
	sp := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("170"))),
	)
	return &FilePicker{filePicker: fp, spinner: sp, onComplete: onComplete}
}

type FilePicker struct {
	filePicker  filepicker.Model
	spinner     spinner.Model
	parentModel tea.Model
	onComplete  func(string) error
	err         error
	loading     bool
}

func (m FilePicker) Init() tea.Cmd {
	return tea.Sequence(m.filePicker.Init(), tea.WindowSize(), m.spinner.Tick)
}

func (m FilePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			parent, cmd := m.parentModel.Update(nil)
			return parent, cmd
		}
	case tui.ClearErrorMsg:
		m.err = nil
	case tui.ErrorMsg:
		m.err = msg.Err
		return m, tui.ClearErrorAfter(time.Second * 2)
	case tui.LoadMsg:
		m.loading = msg.InProcess
	case tui.SpawnMsg:
		m.parentModel = msg.Parent
		return m, tea.Sequence(m.spinner.Tick, m.Init())
	case successMsg:
		newM, cmd := m.parentModel.Update(nil)
		return newM, tea.Sequence(cmd, tea.Printf("File %s uploaded to storage", msg.fileName))
	}

	var cmd tea.Cmd
	m.filePicker, cmd = m.filePicker.Update(msg)

	if didSelect, path := m.filePicker.DidSelectFile(msg); didSelect {
		return m, tea.Sequence(send(path, m.onComplete), cmd)
	}

	if didSelect, path := m.filePicker.DidSelectDisabledFile(msg); didSelect {
		m.err = errors.New(path + " is not valid.")
		return m, tea.Batch(cmd, tui.ClearErrorAfter(2*time.Second))
	}
	var cmdSpin tea.Cmd
	m.spinner, cmdSpin = m.spinner.Update(msg)

	return m, tea.Batch(cmdSpin, cmd)
}

func (m FilePicker) View() string {
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filePicker.Styles.DisabledFile.Render(m.err.Error()))
	}
	s.WriteString("\n\n" + m.filePicker.View() + "\n")
	if m.loading {
		s.WriteString(fmt.Sprintf("\n\n %s Uploading file...", m.spinner.View()))
	}
	return s.String()
}
