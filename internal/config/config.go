package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GitHubRepo     string `yaml:"github_repo"`      // e.g., "username/bookmarks"
	GitHubToken    string `yaml:"github_token"`     // Personal access token
	GitHubBranch   string `yaml:"github_branch"`    // Branch to push to (default: main)
	DebounceMs     int    `yaml:"debounce_ms"`      // Debounce delay in milliseconds (default: 500)
	LogPath        string `yaml:"log_path"`         // Log file path (optional)
	CommitMessage  string `yaml:"commit_message"`   // Custom commit message template
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".bookmarked")
	return filepath.Join(configDir, "config.yaml"), nil
}

// Load reads and parses the config file
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w\nRun 'bookmarked init' to create a config file", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if cfg.GitHubBranch == "" {
		cfg.GitHubBranch = "main"
	}
	if cfg.DebounceMs == 0 {
		cfg.DebounceMs = 500
	}
	if cfg.CommitMessage == "" {
		cfg.CommitMessage = "Update bookmarks"
	}

	// Validate required fields
	if cfg.GitHubRepo == "" {
		return nil, fmt.Errorf("github_repo is required in config file")
	}
	if cfg.GitHubToken == "" {
		return nil, fmt.Errorf("github_token is required in config file")
	}

	return &cfg, nil
}

// InitConfig creates a new config file with template
func InitConfig() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config file already exists at: %s", configPath)
	}

	// Create template config
	template := `# Bookmarked Configuration
# GitHub repository to sync bookmarks to (format: username/repo-name)
github_repo: ""

# GitHub personal access token with repo write permissions
# Create one at: https://github.com/settings/tokens
github_token: ""

# Branch to push to (optional, default: main)
github_branch: "main"

# Debounce delay in milliseconds (optional, default: 500)
debounce_ms: 500

# Log file path (optional, logs to stdout if not set)
log_path: ""

# Commit message template (optional, default: "Update bookmarks")
commit_message: "Update bookmarks"
`

	if err := os.WriteFile(configPath, []byte(template), 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("✓ Config file created at: %s\n", configPath)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Edit the config file and add your GitHub repo and token")
	fmt.Println("2. Run 'bookmarked sync' to test the configuration")
	fmt.Println("3. Run 'bookmarked install' to set up automatic syncing")

	return nil
}

// GetRepoPath returns the local repository path
func GetRepoPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	repoPath := filepath.Join(homeDir, ".bookmarked", "repo")
	return repoPath, nil
}
