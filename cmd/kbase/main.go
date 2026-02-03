package main

import (
	"fmt"
	"os"

	"github.com/BillabongJumbuck/Kbase/internal/config"
	"github.com/BillabongJumbuck/Kbase/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load application configuration
	appConfig, err := config.LoadAppConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading app config: %v\n", err)
		os.Exit(1)
	}

	// Load commands from all configured paths
	commands, loadErr := config.LoadAllCommands(appConfig)

	// Get first command path for editing (backward compatibility)
	var editPath string
	if len(appConfig.CommandPaths) > 0 {
		editPath = appConfig.CommandPaths[0]
	} else {
		// Fallback to default path if no paths configured
		editPath, err = config.GetDefaultConfigPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting default config path: %v\n", err)
			os.Exit(1)
		}
	}

	// Create and run the TUI application
	m := tui.NewModel(editPath, commands, loadErr)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running app: %v\n", err)
		os.Exit(1)
	}
}
