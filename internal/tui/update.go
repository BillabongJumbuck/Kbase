package tui

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/BillabongJumbuck/Kbase/internal/config"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all state updates based on incoming messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tickMsg:
		// Clear status message if expired
		m.ClearStatusIfExpired()
		return m, tickCmd()

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle quit keys (available in all modes)
	if msg.String() == "ctrl+q" || msg.String() == "esc" {
		return m, tea.Quit
	}

	// Error mode: only allow 'e' to edit or quit
	if m.Mode == ErrorMode {
		if msg.String() == "e" {
			return m, m.openEditor()
		}
		return m, nil
	}

	// Normal mode key handling
	switch msg.String() {
	case "ctrl+c":
		return m, m.copyToClipboard()

	case "e":
		return m, m.openEditor()

	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if m.Cursor < len(m.FilteredCommands)-1 {
			m.Cursor++
		}

	case "backspace":
		if len(m.SearchInput) > 0 {
			m.SearchInput = m.SearchInput[:len(m.SearchInput)-1]
			m.updateFilter()
		}

	default:
		// Add character to search input (single printable characters only)
		if len(msg.String()) == 1 && msg.String()[0] >= 32 && msg.String()[0] <= 126 {
			m.SearchInput += msg.String()
			m.updateFilter()
		}
	}

	return m, nil
}

// copyToClipboard copies the selected command to the system clipboard
func (m *Model) copyToClipboard() tea.Cmd {
	if len(m.FilteredCommands) == 0 || m.Cursor >= len(m.FilteredCommands) {
		return nil
	}

	cmd := m.FilteredCommands[m.Cursor].Cmd
	err := clipboard.WriteAll(cmd)

	if err != nil {
		m.SetStatus("Clipboard not supported", 3*time.Second)
		// Fallback: print to stdout on exit
		fmt.Println(cmd)
	} else {
		m.SetStatus("Copied!", 2*time.Second)
	}

	return nil
}

// openEditor opens the config file in the user's $EDITOR
func (m *Model) openEditor() tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // Default to vim
	}

	c := exec.Command(editor, m.ConfigPath)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		// Reload commands after editor closes
		commands, parseErr := config.LoadCommands(m.ConfigPath)
		if parseErr != nil {
			m.Mode = ErrorMode
			m.ParseError = parseErr
		} else {
			m.Mode = NormalMode
			m.AllCommands = commands
			m.updateFilter()
		}
		return nil
	})
}

// updateFilter updates the filtered commands list based on search input
func (m *Model) updateFilter() {
	m.FilteredCommands = config.FilterCommands(m.AllCommands, m.SearchInput)
	// Reset cursor if out of bounds
	if m.Cursor >= len(m.FilteredCommands) {
		m.Cursor = 0
	}
}
