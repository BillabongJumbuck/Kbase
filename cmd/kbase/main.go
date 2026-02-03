package main

import (
	"fmt"
	"os"

	"github.com/BillabongJumbuck/Kbase/internal/config"
	"github.com/BillabongJumbuck/Kbase/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Get config path
	configPath, err := getConfigPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting config path: %v\n", err)
		os.Exit(1)
	}

	// Initialize default config if needed
	if err := config.InitDefaultConfig(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
		os.Exit(1)
	}

	// Load commands
	commands, loadErr := config.LoadCommands(configPath)

	// Create and run the TUI application
	m := tui.NewModel(configPath, commands, loadErr)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running app: %v\n", err)
		os.Exit(1)
	}
}

// getConfigPath returns the config file path from args or default location
func getConfigPath() (string, error) {
	// Check for command-line argument
	if len(os.Args) > 1 {
		return os.Args[1], nil
	}

	// Use default path
	return config.GetDefaultConfigPath()
}
