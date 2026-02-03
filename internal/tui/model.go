package tui

import (
	"time"

	"github.com/BillabongJumbuck/Kbase/internal/model"
	tea "github.com/charmbracelet/bubbletea"
)

// ViewMode represents the current view state
type ViewMode int

const (
	NormalMode ViewMode = iota
	ErrorMode
)

// Model represents the TUI application state
type Model struct {
	ConfigPath       string
	AllCommands      []model.Command
	FilteredCommands []model.Command
	Cursor           int
	SearchInput      string
	Mode             ViewMode
	ParseError       error
	StatusMsg        string
	StatusExpiry     time.Time
	Width            int
	Height           int
}

// TickMsg is sent periodically to update the UI
type tickMsg time.Time

// tickCmd returns a command that sends a tick message every second
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tickCmd()
}

// NewModel creates a new TUI model with the given config path
func NewModel(configPath string, commands []model.Command, err error) Model {
	m := Model{
		ConfigPath:       configPath,
		AllCommands:      commands,
		FilteredCommands: commands,
		Cursor:           0,
		SearchInput:      "",
		Mode:             NormalMode,
		ParseError:       err,
	}

	if err != nil {
		m.Mode = ErrorMode
	}

	return m
}

// SetStatus sets a temporary status message
func (m *Model) SetStatus(msg string, duration time.Duration) {
	m.StatusMsg = msg
	m.StatusExpiry = time.Now().Add(duration)
}

// ClearStatusIfExpired clears the status message if it has expired
func (m *Model) ClearStatusIfExpired() {
	if !m.StatusExpiry.IsZero() && time.Now().After(m.StatusExpiry) {
		m.StatusMsg = ""
		m.StatusExpiry = time.Time{}
	}
}
