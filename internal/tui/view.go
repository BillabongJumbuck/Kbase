package tui

import (
	"fmt"
	"strings"
)

// View renders the TUI
func (m Model) View() string {
	if m.Mode == ErrorMode {
		return m.renderErrorView()
	}
	return m.renderNormalView()
}

// renderErrorView renders the error state view
func (m Model) renderErrorView() string {
	var b strings.Builder

	errorMsg := ErrorStyle.Render(
		fmt.Sprintf("⚠ Parsing Failed\n\n%v\n\nPress 'e' to edit config or Esc to quit", m.ParseError),
	)

	b.WriteString("\n")
	b.WriteString(errorMsg)
	b.WriteString("\n")

	return b.String()
}

// renderNormalView renders the normal mode view
func (m Model) renderNormalView() string {
	var b strings.Builder

	// Search box
	searchBox := SearchBoxStyle.Render(fmt.Sprintf("Search: %s_", m.SearchInput))
	b.WriteString(searchBox)
	b.WriteString("\n\n")

	// Command list
	listHeight := m.Height - 8 // Reserve space for search, detail, and status
	if listHeight < 5 {
		listHeight = 5
	}

	// Calculate visible range centered on cursor
	start := m.Cursor - listHeight/2
	if start < 0 {
		start = 0
	}
	end := start + listHeight
	if end > len(m.FilteredCommands) {
		end = len(m.FilteredCommands)
		start = end - listHeight
		if start < 0 {
			start = 0
		}
	}

	// Render command list
	for i := start; i < end && i < len(m.FilteredCommands); i++ {
		cmd := m.FilteredCommands[i]

		// Truncate description for list view
		desc := cmd.Desc
		maxDescLen := 50
		if len(desc) > maxDescLen {
			desc = desc[:maxDescLen] + "..."
		}

		cursor := " "
		if i == m.Cursor {
			cursor = ">"
		}

		cmdText := cmd.Cmd
		descText := DescStyle.Render(desc)

		if i == m.Cursor {
			line := fmt.Sprintf("%s %s  %s", cursor, cmdText, descText)
			b.WriteString(SelectedItemStyle.Render(line))
		} else {
			line := fmt.Sprintf("%s %s  %s", cursor, cmdText, descText)
			b.WriteString(NormalItemStyle.Render(line))
		}
		b.WriteString("\n")
	}

	// Show detail/examples for selected item
	if len(m.FilteredCommands) > 0 && m.Cursor < len(m.FilteredCommands) {
		selected := m.FilteredCommands[m.Cursor]

		if len(selected.Examples) > 0 {
			b.WriteString("\n")
			detailContent := fmt.Sprintf("Examples:\n%s", strings.Join(selected.Examples, "\n"))
			b.WriteString(DetailBoxStyle.Render(ExampleStyle.Render(detailContent)))
		} else if len(selected.Desc) > 50 {
			// Show full description if truncated
			b.WriteString("\n")
			b.WriteString(DetailBoxStyle.Render(selected.Desc))
		}
	}

	b.WriteString("\n\n")

	// Status bar
	statusText := "Ctrl+C: Copy | E: Edit | ↑/↓ or k/j: Navigate | Esc/Ctrl+Q: Quit"
	if m.StatusMsg != "" {
		b.WriteString(StatusBarSuccessStyle.Render(m.StatusMsg))
	} else {
		b.WriteString(StatusBarStyle.Render(statusText))
	}

	return b.String()
}
