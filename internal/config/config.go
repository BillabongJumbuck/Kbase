package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BillabongJumbuck/Kbase/internal/model"
	"gopkg.in/yaml.v3"
)

// AppConfig represents the application configuration
type AppConfig struct {
	CommandPaths []string `yaml:"command_paths"` // List of paths to command YAML files
}

// LoadCommands reads and parses the YAML file, filters by platform, and returns commands
func LoadCommands(path string) ([]model.Command, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var commands []model.Command
	if err := yaml.Unmarshal(data, &commands); err != nil {
		return nil, err
	}

	return FilterByPlatform(commands), nil
}

// FilterByPlatform filters commands based on current OS
func FilterByPlatform(commands []model.Command) []model.Command {
	currentOS := runtime.GOOS
	filtered := make([]model.Command, 0, len(commands))

	for _, cmd := range commands {
		if len(cmd.Platform) == 0 {
			// No platform restriction, include it
			filtered = append(filtered, cmd)
		} else {
			// Check if current OS is in the platform list
			for _, p := range cmd.Platform {
				if p == currentOS {
					filtered = append(filtered, cmd)
					break
				}
			}
		}
	}

	return filtered
}

// InitDefaultConfig creates the config directory and a default YAML file if they don't exist
func InitDefaultConfig(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); err == nil {
		return nil // File already exists
	}

	// Create directory
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create default YAML template
	defaultYAML := `- cmd: "kubectl get pods"
  desc: "List all pods in namespace"
  tags:
    - "k8s"
    - "container"
  platform:
    - "linux"
    - "darwin"
  examples:
    - "kubectl get pods -n kube-system -o wide"
    - "kubectl get pods --watch"

- cmd: "docker ps -a"
  desc: "List all containers"
  tags:
    - "docker"
    - "container"

- cmd: "git log --oneline --graph --all"
  desc: "Show git commit graph"
  tags:
    - "git"
    - "vcs"

- cmd: "find . -name '*.go' -type f"
  desc: "Find all Go files in current directory"
  tags:
    - "shell"
    - "find"

- cmd: "ps aux | grep <process>"
  desc: "Search for running processes"
  tags:
    - "shell"
    - "process"
  platform:
    - "linux"
    - "darwin"
`

	return os.WriteFile(path, []byte(defaultYAML), 0644)
}

// GetDefaultConfigPath returns the default config path: ~/.config/kbase/commands.yaml
func GetDefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "kbase", "commands.yaml"), nil
}

// GetAppConfigPath returns the platform-specific application config file path
func GetAppConfigPath() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%\Kbase\config.yaml
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		configDir = filepath.Join(appData, "Kbase")
	case "darwin":
		// macOS: ~/Library/Application Support/Kbase/config.yaml
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, "Library", "Application Support", "Kbase")
	default:
		// Linux and others: ~/.config/Kbase/config.yaml
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, ".config", "Kbase")
	}

	return filepath.Join(configDir, "config.yaml"), nil
}

// LoadAppConfig loads the application configuration from the config file
func LoadAppConfig() (*AppConfig, error) {
	configPath, err := GetAppConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, create default one
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := InitAppConfig(configPath); err != nil {
			return nil, fmt.Errorf("failed to initialize config: %w", err)
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Expand environment variables and home directory in paths
	for i, path := range config.CommandPaths {
		config.CommandPaths[i] = expandPath(path)
	}

	return &config, nil
}

// InitAppConfig creates a default application config file
func InitAppConfig(configPath string) error {
	// Create directory
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Get default command path
	defaultCmdPath, err := GetDefaultConfigPath()
	if err != nil {
		return err
	}

	// Create default config
	defaultConfig := AppConfig{
		CommandPaths: []string{defaultCmdPath},
	}

	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return err
	}

	// Add comments to the YAML file
	header := `# Kbase Application Configuration
# This file defines where to load command data from

# command_paths: List of paths to command YAML files or directories
# You can specify multiple paths, and commands from all files will be loaded
# - If a path is a file, it will be loaded directly
# - If a path is a directory, all .yaml and .yml files in it will be loaded (non-recursive)
# Paths support:
#   - Absolute paths: /home/user/.config/kbase/commands.yaml
#   - Directory paths: /home/user/.config/kbase/commands/
#   - Home directory expansion: ~/.config/kbase/commands.yaml
#   - Environment variables: $HOME/.config/kbase/commands.yaml

`

	return os.WriteFile(configPath, append([]byte(header), data...), 0644)
}

// LoadAllCommands loads commands from all configured paths
// Paths can be either files or directories
// If a directory is specified, all .yaml and .yml files in it will be loaded
func LoadAllCommands(config *AppConfig) ([]model.Command, error) {
	var allCommands []model.Command

	for _, path := range config.CommandPaths {
		info, err := os.Stat(path)
		if err != nil {
			// If path doesn't exist, try to initialize it as a file
			if os.IsNotExist(err) {
				if err := InitDefaultConfig(path); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to initialize %s: %v\n", path, err)
					continue
				}
				// Try to load it now
				commands, err := LoadCommands(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to load commands from %s: %v\n", path, err)
					continue
				}
				allCommands = append(allCommands, commands...)
				continue
			}
			fmt.Fprintf(os.Stderr, "Warning: failed to access %s: %v\n", path, err)
			continue
		}

		if info.IsDir() {
			// Load all YAML files from directory
			commands, err := LoadCommandsFromDirectory(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to load commands from directory %s: %v\n", path, err)
				continue
			}
			allCommands = append(allCommands, commands...)
		} else {
			// Load single file
			commands, err := LoadCommands(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to load commands from %s: %v\n", path, err)
				continue
			}
			allCommands = append(allCommands, commands...)
		}
	}

	return allCommands, nil
}

// LoadCommandsFromDirectory loads all YAML files from a directory
func LoadCommandsFromDirectory(dirPath string) ([]model.Command, error) {
	var allCommands []model.Command

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip subdirectories
		}

		// Check if file has .yaml or .yml extension
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".yaml") && !strings.HasSuffix(strings.ToLower(name), ".yml") {
			continue
		}

		filePath := filepath.Join(dirPath, name)
		commands, err := LoadCommands(filePath)
		if err != nil {
			// Log warning but continue with other files
			fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", filePath, err)
			continue
		}

		allCommands = append(allCommands, commands...)
	}

	return allCommands, nil
}

// expandPath expands ~ and environment variables in path
func expandPath(path string) string {
	// Expand environment variables
	path = os.ExpandEnv(path)

	// Expand home directory
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}

	return path
}

// FuzzyMatch checks if query matches cmd, desc, or tags
func FuzzyMatch(cmd model.Command, query string) bool {
	query = strings.ToLower(query)

	// Match against cmd
	if strings.Contains(strings.ToLower(cmd.Cmd), query) {
		return true
	}

	// Match against desc
	if strings.Contains(strings.ToLower(cmd.Desc), query) {
		return true
	}

	// Match against tags
	for _, tag := range cmd.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}

	return false
}

// FilterCommands returns commands that match the query
func FilterCommands(commands []model.Command, query string) []model.Command {
	if query == "" {
		return commands
	}

	filtered := make([]model.Command, 0)
	for _, cmd := range commands {
		if FuzzyMatch(cmd, query) {
			filtered = append(filtered, cmd)
		}
	}

	return filtered
}
