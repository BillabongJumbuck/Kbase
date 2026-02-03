package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BillabongJumbuck/Kbase/internal/model"
	"gopkg.in/yaml.v3"
)

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
