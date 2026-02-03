package tui

import "github.com/charmbracelet/lipgloss"

var (
	// SearchBoxStyle defines the style for the search input box
	SearchBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	// SelectedItemStyle defines the style for the currently selected item
	SelectedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	// NormalItemStyle defines the style for normal (unselected) items
	NormalItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	// DescStyle defines the style for command descriptions (dimmed)
	DescStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// StatusBarStyle defines the default status bar style
	StatusBarStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	// StatusBarSuccessStyle defines the status bar style for success messages
	StatusBarSuccessStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("42")).
		Foreground(lipgloss.Color("230")).
		Bold(true).
		Padding(0, 1)

	// ErrorStyle defines the style for error messages
	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("196")).
		Padding(1, 2)

	// DetailBoxStyle defines the style for the detail/examples pane
	DetailBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		MarginTop(1)

	// ExampleStyle defines the style for example commands
	ExampleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("114"))
)
