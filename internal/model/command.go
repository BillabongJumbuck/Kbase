package model

// Command represents a single command entry in the knowledge base
type Command struct {
	Cmd      string   `yaml:"cmd"`      // The actual command (required)
	Desc     string   `yaml:"desc"`     // Short description (required)
	Tags     []string `yaml:"tags,omitempty"`     // Tags for classification and search
	Platform []string `yaml:"platform,omitempty"` // Platform constraints (e.g., "linux", "darwin", "windows")
	Examples []string `yaml:"examples,omitempty"` // Usage examples
}
