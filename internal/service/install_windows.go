//go:build windows

package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func installWindows() error {
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

	// Create log directory
	logDir := filepath.Join(homeDir, ".bookmarked")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logPath := filepath.Join(logDir, "bookmarked.log")

	// Create a scheduled task that runs at login
	taskName := "Bookmarked"

	// Delete existing task if it exists
	exec.Command("schtasks", "/Delete", "/TN", taskName, "/F").Run()

	// Create new task
	cmd := exec.Command("schtasks", "/Create",
		"/TN", taskName,
		"/TR", fmt.Sprintf("\"%s\" start", exePath),
		"/SC", "ONLOGON",
		"/RL", "HIGHEST",
		"/F",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create scheduled task: %w\nOutput: %s", err, string(output))
	}

	fmt.Println("✓ Service installed successfully")
	fmt.Println("  The service will start automatically when you log in")
	fmt.Printf("  Logs: %s\n", logPath)
	fmt.Println("\nTo start now, run: bookmarked start")
	fmt.Println("To check status, run: bookmarked status")

	return nil
}

func uninstallWindows() error {
	taskName := "Bookmarked"

	cmd := exec.Command("schtasks", "/Delete", "/TN", taskName, "/F")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete scheduled task: %w\nOutput: %s", err, string(output))
	}

	fmt.Println("✓ Service uninstalled successfully")
	return nil
}

func statusWindows() error {
	taskName := "Bookmarked"

	cmd := exec.Command("schtasks", "/Query", "/TN", taskName, "/FO", "LIST", "/V")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("✗ Service is not installed")
		return nil
	}

	fmt.Println("✓ Service is installed")
	fmt.Println("\nTask details:")
	fmt.Println(string(output))

	return nil
}


// Stub functions for other platforms
func installMacOS() error {
	return fmt.Errorf("macOS installation not available on Windows")
}

func installLinux() error {
	return fmt.Errorf("Linux installation not available on Windows")
}

func uninstallMacOS() error {
	return fmt.Errorf("macOS uninstallation not available on Windows")
}

func uninstallLinux() error {
	return fmt.Errorf("Linux uninstallation not available on Windows")
}

func statusMacOS() error {
	return fmt.Errorf("macOS status not available on Windows")
}

func statusLinux() error {
	return fmt.Errorf("Linux status not available on Windows")
}
