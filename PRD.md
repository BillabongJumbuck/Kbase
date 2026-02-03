# Product Requirements Document (PRD) - Kbase

## 1. Project Overview

**Kbase** is a lightweight, terminal-based (TUI) command knowledge base tool. It is designed to help developers quickly store, retrieve, and review frequently used commands. It adopts a "Local-First" strategy, utilizing YAML files for data persistence. It supports cross-platform execution (Linux, Windows/WSL, Android/Termux) and features millisecond-level startup with instant fuzzy search capabilities.

## 2. Core Values

- **Speed:** Instant startup and zero-latency search.
- **Minimalism:** Focus strictly on "Lookup" and "Memorization," avoiding complex social or cloud sync features.
- **Workflow Integration:** Operates entirely within the terminal, supports Git-based synchronization (user-managed), and integrates seamlessly with the system's `$EDITOR`.

## 3. Tech Stack

- **Language:** Go (Golang)
- **UI Framework:** `charmbracelet/bubbletea` (TUI Interaction)
- **Styling:** `charmbracelet/lipgloss` (UI Beautification)
- **Data Parsing:** `gopkg.in/yaml.v3`
- **Clipboard:** `github.com/atotto/clipboard` or `golang.design/x/clipboard`

## 4. Data Schema

Data will be stored in `~/.config/kbase/commands.yaml` (or a user-specified path).

### 4.1 YAML Structure Definition

YAML

```
- cmd: "kubectl get pods"             # [Required] The actual command
  desc: "List all pods in namespace"  # [Required] Short description
  tags:                               # [Optional] For classification and search weighting
    - "k8s"
    - "container"
  platform:                           # [Optional] Platform constraints. Shows on all if empty.
    - "linux"
    - "darwin"                        # Visible on Linux/macOS, hidden on Windows
  examples:                           # [Optional] Code blocks for specific usage
    - "kubectl get pods -n kube-system -o wide"
    - "kubectl get pods --watch"
```

## 5. Functional Requirements

### 5.1 Startup & Loading

- **Initialization:** Upon launch, the app reads the YAML file from the default path. If the file does not exist, a template file with example data is automatically created.
- **Platform Filtering:** The app detects the current OS (`runtime.GOOS`). If a record has a non-empty `platform` field that does not include the current OS, that record is hidden from the list.

### 5.2 List & Search (Main View)

- **UI Layout:**
  - **Top:** Search input field.
  - **Middle:** Scrollable command list (Left: `cmd`, Right: `desc` in dimmed color).
  - **Bottom:** Status bar (showing key hints).
- **Fuzzy Search:** Real-time filtering as the user types. The search scope includes `cmd`, `desc`, and `tags`.

### 5.3 Detail View

- **Interaction:** When highlighting an item, if the `examples` field exists, the detailed examples should be displayed (either in a side pane for wide screens or a bottom pane).
- **Truncation:** Long `desc` text in the list view should be truncated, with the full text visible in the detail view.

### 5.4 Clipboard Operation (Copy)

- **Trigger:** User presses `Ctrl+c`.
- **Action:** The content of the `cmd` field of the currently selected item is copied to the system clipboard.
- **Feedback:** The status bar temporarily turns green displaying "Copied!", then reverts.
- **Note:** Since `Ctrl+c` is standard for "Interrupt/Exit," the application must intercept `SIGINT` or override the default key mapping in Bubble Tea to perform the Copy action instead.

### 5.5 Edit Functionality

- **Trigger:** User presses `e`.
- **Workflow:**
  1. **Suspend:** The TUI suspends its render loop.
  2. **External Process:** The app executes the command defined in `$EDITOR` (e.g., vim, nano, code), opening the `commands.yaml` file.
  3. **Resume:** Upon closing the editor, the TUI resumes.
  4. **Reload:** The application automatically re-parses the YAML file and refreshes the list to reflect changes.

## 6. Interactions & Keybindings

| **Key**          | **Action**             | **Description**                                              |
| ---------------- | ---------------------- | ------------------------------------------------------------ |
| `Ctrl+c`         | **Copy**               | Copies the selected command to the clipboard.                |
| `Esc` / `Ctrl+q` | **Quit**               | Exits the application (Replaces standard Ctrl+c behavior).   |
| `e`              | **Edit**               | Opens the source YAML file in the default `$EDITOR`.         |
| `Enter`          | **No Action / Expand** | (Optional) Expands details or does nothing (as per user preference). |
| `↑` / `k`        | **Up**                 | Navigate up the list.                                        |
| `↓` / `j`        | **Down**               | Navigate down the list.                                      |
| `/`              | **Search**             | Focuses the search bar (if not already focused).             |

## 7. Error Handling

- **YAML Syntax Error:** If the YAML file is corrupted after editing, the app must **not crash**. Instead, it should display a red `ErrorView` indicating "Parsing Failed" and allow the user to press `e` again to fix the file.
- **Clipboard Failure:** In headless or SSH-only environments where the clipboard is inaccessible, the status bar should display "Clipboard not supported" and the command should be printed to `stdout` upon exit (optional fallback).