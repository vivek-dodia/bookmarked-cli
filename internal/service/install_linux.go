//go:build linux

package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const systemdTemplate = `[Unit]
Description=Bookmarked - Chrome Bookmark Sync Service
After=network.target

[Service]
Type=simple
ExecStart={{.ExePath}} start
Restart=on-failure
RestartSec=10

[Install]
WantedBy=default.target
`

func installLinux() error {
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

	// Create systemd user directory
	systemdDir := filepath.Join(homeDir, ".config", "systemd", "user")
	if err := os.MkdirAll(systemdDir, 0755); err != nil {
		return fmt.Errorf("failed to create systemd directory: %w", err)
	}

	// Create log directory
	logDir := filepath.Join(homeDir, ".bookmarked")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	servicePath := filepath.Join(systemdDir, "bookmarked.service")
	logPath := filepath.Join(logDir, "bookmarked.log")

	// Create service file
	tmpl, err := template.New("service").Parse(systemdTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(servicePath)
	if err != nil {
		return fmt.Errorf("failed to create service file: %w", err)
	}
	defer file.Close()

	data := struct {
		ExePath string
	}{
		ExePath: exePath,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	// Reload systemd
	exec.Command("systemctl", "--user", "daemon-reload").Run()

	// Enable the service
	cmd := exec.Command("systemctl", "--user", "enable", "bookmarked.service")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to enable service: %w\nOutput: %s", err, string(output))
	}

	fmt.Println("✓ Service installed successfully")
	fmt.Println("  The service will start automatically when you log in")
	fmt.Printf("  Logs: %s\n", logPath)
	fmt.Println("\nTo start now, run: systemctl --user start bookmarked")
	fmt.Println("To check status, run: systemctl --user status bookmarked")
	fmt.Println("Or use: bookmarked status")

	return nil
}

func uninstallLinux() error {
	// Stop the service
	exec.Command("systemctl", "--user", "stop", "bookmarked.service").Run()

	// Disable the service
	exec.Command("systemctl", "--user", "disable", "bookmarked.service").Run()

	// Remove service file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	servicePath := filepath.Join(homeDir, ".config", "systemd", "user", "bookmarked.service")
	if err := os.Remove(servicePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	// Reload systemd
	exec.Command("systemctl", "--user", "daemon-reload").Run()

	fmt.Println("✓ Service uninstalled successfully")
	return nil
}

func statusLinux() error {
	cmd := exec.Command("systemctl", "--user", "status", "bookmarked.service")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("✗ Service is not installed or not running")
		return nil
	}

	fmt.Println("Service status:")
	fmt.Println(string(output))
	return nil
}
