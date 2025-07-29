package ui

import (
	"errors"
	"gophkeeper/internal/client/tui"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

func NewFilePicker(dir string) *FilePicker {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".mod", ".sum", ".go", ".txt", ".md"}
	fp.CurrentDirectory = dir
	return &FilePicker{filePicker: fp}
}

type FilePicker struct {
	filePicker   filepicker.Model
	parentModel  tea.Model
	selectedFile string
	err          error
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
		m.selectedFile = path
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filePicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m FilePicker) View() string {
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		log.Println("error before view")
		s.WriteString(m.filePicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString("Pick a file:")
	} else {
		s.WriteString("Selected file: " + m.filePicker.Styles.Selected.Render(m.selectedFile))
	}
	s.WriteString("\n\n" + m.filePicker.View() + "\n")
	return s.String()
}
