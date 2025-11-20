package bookmarks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetBookmarkPath returns the Chrome bookmarks file path for the current OS
func GetBookmarkPath() (string, error) {
	var bookmarkPath string

	switch runtime.GOOS {
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return "", fmt.Errorf("LOCALAPPDATA environment variable not set")
		}
		bookmarkPath = filepath.Join(localAppData, "Google", "Chrome", "User Data", "Default", "Bookmarks")

	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		bookmarkPath = filepath.Join(homeDir, "Library", "Application Support", "Google", "Chrome", "Default", "Bookmarks")

	case "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		bookmarkPath = filepath.Join(homeDir, ".config", "google-chrome", "Default", "Bookmarks")

	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Check if file exists
	if _, err := os.Stat(bookmarkPath); os.IsNotExist(err) {
		return "", fmt.Errorf("Chrome bookmarks file not found at: %s", bookmarkPath)
	}

	return bookmarkPath, nil
}

// FormatBookmarks reads the Chrome bookmarks file and returns formatted JSON
func FormatBookmarks(bookmarkPath string) ([]byte, error) {
	// Read the raw bookmarks file
	data, err := os.ReadFile(bookmarkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read bookmarks file: %w", err)
	}

	// Parse JSON
	var bookmarks interface{}
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		return nil, fmt.Errorf("failed to parse bookmarks JSON: %w", err)
	}

	// Format with indentation for better diffs
	formatted, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to format bookmarks: %w", err)
	}

	return formatted, nil
}

// CopyToRepo copies and formats bookmarks to the target repository path
func CopyToRepo(bookmarkPath, repoPath string) error {
	formatted, err := FormatBookmarks(bookmarkPath)
	if err != nil {
		return err
	}

	targetPath := filepath.Join(repoPath, "Bookmarks.json")
	if err := os.WriteFile(targetPath, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write formatted bookmarks: %w", err)
	}

	return nil
}
