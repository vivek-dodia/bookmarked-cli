package service

import (
	"fmt"
	"runtime"
)

// Install installs the service for the current platform
func Install() error {
	switch runtime.GOOS {
	case "windows":
		return installWindows()
	case "darwin":
		return installMacOS()
	case "linux":
		return installLinux()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// Uninstall removes the service for the current platform
func Uninstall() error {
	switch runtime.GOOS {
	case "windows":
		return uninstallWindows()
	case "darwin":
		return uninstallMacOS()
	case "linux":
		return uninstallLinux()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// Status checks the service status for the current platform
func Status() error {
	switch runtime.GOOS {
	case "windows":
		return statusWindows()
	case "darwin":
		return statusMacOS()
	case "linux":
		return statusLinux()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
