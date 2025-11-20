//go:build darwin

package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.bookmarked.sync</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.ExePath}}</string>
        <string>start</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <false/>
    <key>StandardOutPath</key>
    <string>{{.LogPath}}</string>
    <key>StandardErrorPath</key>
    <string>{{.LogPath}}</string>
</dict>
</plist>
`

func installMacOS() error {
	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create LaunchAgents directory
	launchAgentsDir := filepath.Join(homeDir, "Library", "LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create LaunchAgents directory: %w", err)
	}

	// Create log directory
	logDir := filepath.Join(homeDir, ".bookmarked")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	plistPath := filepath.Join(launchAgentsDir, "com.bookmarked.sync.plist")
	logPath := filepath.Join(logDir, "bookmarked.log")

	// Create plist file
	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(plistPath)
	if err != nil {
		return fmt.Errorf("failed to create plist file: %w", err)
	}
	defer file.Close()

	data := struct {
		ExePath string
		LogPath string
	}{
		ExePath: exePath,
		LogPath: logPath,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}

	// Load the launch agent
	cmd := exec.Command("launchctl", "load", plistPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to load launch agent: %w\nOutput: %s", err, string(output))
	}

	fmt.Println("✓ Service installed successfully")
	fmt.Println("  The service will start automatically when you log in")
	fmt.Printf("  Logs: %s\n", logPath)
	fmt.Println("\nTo start now, run: bookmarked start")
	fmt.Println("To check status, run: bookmarked status")

	return nil
}

func uninstallMacOS() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	plistPath := filepath.Join(homeDir, "Library", "LaunchAgents", "com.bookmarked.sync.plist")

	// Unload the launch agent
	cmd := exec.Command("launchctl", "unload", plistPath)
	cmd.Run() // Ignore errors if not loaded

	// Remove plist file
	if err := os.Remove(plistPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove plist file: %w", err)
	}

	fmt.Println("✓ Service uninstalled successfully")
	return nil
}

func statusMacOS() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	plistPath := filepath.Join(homeDir, "Library", "LaunchAgents", "com.bookmarked.sync.plist")

	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		fmt.Println("✗ Service is not installed")
		return nil
	}

	// Check if loaded
	cmd := exec.Command("launchctl", "list", "com.bookmarked.sync")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("✓ Service is installed but not running")
	} else {
		fmt.Println("✓ Service is installed and running")
		fmt.Println("\nService details:")
		fmt.Println(string(output))
	}

	return nil
}
