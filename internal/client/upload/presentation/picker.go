package presentation

import (
	"errors"
	"gophkeeper/internal/client/tui"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

func NewFilePicker(dir string, onComplete func(file string) error) *FilePicker {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".md", ".txt", ".dat"}
	fp.CurrentDirectory = dir
	return &FilePicker{filePicker: fp, onComplete: onComplete}
}

type FilePicker struct {
	filePicker  filepicker.Model
	parentModel tea.Model
	onComplete  func(string) error
	err         error
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m FilePicker) Init() tea.Cmd {
	return tea.Sequence(m.filePicker.Init(), tea.WindowSize())
}

func (m FilePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			parent, cmd := m.parentModel.Update(nil)
			return parent, cmd
		}
	case clearErrorMsg:
		m.err = nil
	case tui.SpawnMsg:
		m.parentModel = msg.Parent
		return m, m.Init()
	}

	var cmd tea.Cmd
	m.filePicker, cmd = m.filePicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filePicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.err = m.onComplete(path)
		if m.err != nil {
			return m, tea.Batch(cmd, clearErrorAfter(5*time.Second))
		}
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filePicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m FilePicker) View() string {
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filePicker.Styles.DisabledFile.Render(m.err.Error()))
	}
	s.WriteString("\n\n" + m.filePicker.View() + "\n")
	return s.String()
}
