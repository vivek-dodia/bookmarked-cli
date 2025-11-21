# Bookmarked CLI

A minimal, cross-platform background service that automatically syncs your Chrome bookmarks to a GitHub repository.

![Platform Support](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-blue)
![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- **Automatic Sync**: Watches Chrome bookmarks and syncs changes to GitHub in real-time
- **Cross-Platform**: Works seamlessly on Windows, macOS, and Linux
- **Minimal**: Single binary with zero runtime dependencies (only Git required)
- **Formatted JSON**: Bookmarks saved as pretty-printed JSON for readable diffs
- **Smart Debouncing**: 500ms debounce prevents excessive commits during bulk operations
- **Secure**: Uses GitHub tokens, private repository recommended
- **Background Service**: Runs silently with automatic restart on failure
- **Version Control**: Full history of bookmark changes via Git

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
  - [Windows](#windows)
  - [macOS](#macos)
  - [Linux](#linux)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage](#usage)
- [How It Works](#how-it-works)
- [Troubleshooting](#troubleshooting)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Prerequisites

- **Go 1.21+** (for building from source)
- **Git** (required for syncing)
- **Chrome/Chromium browser** with at least one bookmark
- **GitHub account** with a repository for storing bookmarks
- **GitHub Personal Access Token** with `repo` scope

## Installation

### Windows

#### Option 1: Download Pre-built Binary (Recommended)

1. Download the latest Windows binary from [Releases](https://github.com/vivek-dodia/bookmarked-cli/releases)
2. Extract `bookmarked.exe` to a folder (e.g., `C:\Program Files\Bookmarked\`)
3. Add the folder to your PATH or use full path

#### Option 2: Build from Source

```powershell
# Install Go from https://golang.org/dl/ or using winget
winget install GoLang.Go

# Clone the repository
git clone git@github.com:vivek-dodia/bookmarked-cli.git
cd bookmarked-cli

# Build the binary
go build -o bookmarked.exe ./cmd/bookmarked

# Optionally, move to a directory in PATH
move bookmarked.exe C:\Windows\System32\
```

#### Quick Install Script

```powershell
# Download and run the install script
Invoke-WebRequest -Uri https://raw.githubusercontent.com/vivek-dodia/bookmarked-cli/main/install.ps1 -OutFile install.ps1
.\install.ps1
```

### macOS

#### Option 1: Download Pre-built Binary (Recommended)

```bash
# Download latest release
curl -LO https://github.com/vivek-dodia/bookmarked-cli/releases/latest/download/bookmarked-darwin-amd64

# Make executable
chmod +x bookmarked-darwin-amd64

# Move to PATH
sudo mv bookmarked-darwin-amd64 /usr/local/bin/bookmarked
```

#### Option 2: Build from Source

```bash
# Install Go (if not already installed)
brew install go

# Clone the repository
git clone git@github.com:vivek-dodia/bookmarked-cli.git
cd bookmarked-cli

# Build the binary
go build -o bookmarked ./cmd/bookmarked

# Move to PATH
sudo mv bookmarked /usr/local/bin/
```

#### Quick Install Script

```bash
# Download and run the install script
curl -fsSL https://raw.githubusercontent.com/vivek-dodia/bookmarked-cli/main/install.sh | bash
```

### Linux

#### Option 1: Download Pre-built Binary (Recommended)

```bash
# Download latest release (adjust architecture as needed)
curl -LO https://github.com/vivek-dodia/bookmarked-cli/releases/latest/download/bookmarked-linux-amd64

# Make executable
chmod +x bookmarked-linux-amd64

# Move to PATH
sudo mv bookmarked-linux-amd64 /usr/local/bin/bookmarked
```

#### Option 2: Build from Source

```bash
# Install Go (if not already installed)
# For Ubuntu/Debian:
sudo apt update
sudo apt install golang-go

# For Fedora/RHEL:
sudo dnf install golang

# For Arch:
sudo pacman -S go

# Clone the repository
git clone git@github.com:vivek-dodia/bookmarked-cli.git
cd bookmarked-cli

# Build the binary
go build -o bookmarked ./cmd/bookmarked

# Move to PATH
sudo mv bookmarked /usr/local/bin/
```

#### Quick Install Script

```bash
# Download and run the install script
curl -fsSL https://raw.githubusercontent.com/vivek-dodia/bookmarked-cli/main/install.sh | bash
```

## Quick Start

### 1. Create a GitHub Repository

Create a new **private** repository on GitHub to store your bookmarks:

```bash
# Example: https://github.com/username/my-bookmarks
```

**Important**: Use a **private repository** since bookmarks may contain sensitive URLs.

### 2. Generate GitHub Personal Access Token

1. Visit [GitHub Settings → Tokens](https://github.com/settings/tokens)
2. Click **"Generate new token (classic)"**
3. Give it a name: `Bookmarked Sync`
4. Select scope: **`repo`** (Full control of private repositories)
5. Click **"Generate token"**
6. **Copy the token** (starts with `ghp_...`) - you won't see it again!

### 3. Initialize Configuration

```bash
bookmarked init
```

This creates a config file at:
- **Windows**: `C:\Users\<username>\.bookmarked\config.yaml`
- **macOS**: `~/.bookmarked/config.yaml`
- **Linux**: `~/.bookmarked/config.yaml`

### 4. Configure Your Settings

Edit the config file and add your GitHub repository and token:

```yaml
github_repo: "your-username/my-bookmarks"
github_token: "ghp_your_token_here"
```

### 5. Test the Sync

```bash
bookmarked sync
```

If successful, you should see your bookmarks in the GitHub repository!

### 6. Install as Background Service

```bash
bookmarked install
```

**Platform-specific behavior:**
- **Windows**: Creates a scheduled task that runs on login
- **macOS**: Creates a launchd agent
- **Linux**: Creates a systemd user service

### 7. Start the Service

The service will start automatically on next login, or start it manually:

```bash
# Windows: Service starts automatically on login
# Or use Task Scheduler to start it manually

# macOS:
launchctl start com.bookmarked.sync

# Linux:
systemctl --user start bookmarked
```

**That's it!** Your bookmarks will now sync automatically to GitHub whenever you make changes.

## Configuration

The configuration file is located at `~/.bookmarked/config.yaml` (or `C:\Users\<username>\.bookmarked\config.yaml` on Windows).

### Configuration Options

```yaml
# GitHub repository (format: username/repo-name)
github_repo: "your-username/my-bookmarks"

# GitHub personal access token with repo permissions
github_token: "ghp_xxxxxxxxxxxxx"

# Branch to push to (default: main)
github_branch: "main"

# Debounce delay in milliseconds (default: 500)
# Prevents excessive commits during bulk bookmark operations
debounce_ms: 500

# Log file path (optional, logs to stdout if not set)
log_path: ""

# Commit message template (default: "Update bookmarks")
commit_message: "Update bookmarks"
```

### Chrome Bookmark Locations

The tool automatically detects Chrome bookmarks based on your OS:

| Platform | Bookmark File Location |
|----------|------------------------|
| **Windows** | `%LOCALAPPDATA%\Google\Chrome\User Data\Default\Bookmarks` |
| **macOS** | `~/Library/Application Support/Google/Chrome/Default/Bookmarks` |
| **Linux** | `~/.config/google-chrome/Default/Bookmarks` |

**Note**: Only the `Default` profile is synced. For multiple profiles, modify the code or run multiple instances.

## Usage

### Commands

```bash
# Initialize configuration file
bookmarked init

# Manually trigger a one-time sync
bookmarked sync

# Start the service in foreground (see live logs)
bookmarked start

# Install as a background service
bookmarked install

# Uninstall the background service
bookmarked uninstall

# Check service status
bookmarked status
```

### Managing the Service

#### Windows

```powershell
# Check status
bookmarked status

# View logs
type C:\Users\<username>\.bookmarked\bookmarked.log

# Manually start via Task Scheduler
# Open Task Scheduler → Find "Bookmarked" task → Right-click → Run
```

#### macOS

```bash
# Start service
launchctl start com.bookmarked.sync

# Stop service
launchctl stop com.bookmarked.sync

# Check status
bookmarked status

# View logs
tail -f ~/.bookmarked/bookmarked.log
```

#### Linux

```bash
# Start service
systemctl --user start bookmarked

# Stop service
systemctl --user stop bookmarked

# Restart service
systemctl --user restart bookmarked

# Check status
systemctl --user status bookmarked
# or
bookmarked status

# View logs
journalctl --user -u bookmarked -f
```

## How It Works

1. **File Watching**: Monitors Chrome's bookmarks file using filesystem events (`fsnotify`)
2. **Debouncing**: Waits 500ms after detecting changes to avoid excessive commits during bulk operations
3. **Formatting**: Reads Chrome's JSON bookmarks and formats with pretty-printing for readable diffs
4. **Git Operations**:
   - Pulls latest changes from GitHub (handles multi-device scenarios)
   - Copies and formats bookmarks to local repository
   - Commits changes with timestamp
   - Pushes to GitHub
5. **Background Service**: Runs continuously, watching for changes and syncing automatically

### Data Flow

```
Chrome Bookmarks File
         ↓
   File Watcher (fsnotify)
         ↓
   Debounce (500ms)
         ↓
   Read & Format JSON
         ↓
   Local Git Repo (~/.bookmarked/repo/)
         ↓
   Git Commit & Push
         ↓
   GitHub Repository
```

## Troubleshooting

### Bookmarks not syncing?

```bash
# Check if service is running
bookmarked status

# Try manual sync to see errors
bookmarked sync

# Check logs
# Windows:
type C:\Users\<username>\.bookmarked\bookmarked.log

# macOS/Linux:
cat ~/.bookmarked/bookmarked.log
```

### "Chrome bookmarks file not found"

- Ensure Chrome is installed and you've created at least one bookmark
- Verify the bookmark file exists at the expected location (see [Chrome Bookmark Locations](#chrome-bookmark-locations))
- For non-default Chrome profiles, you'll need to modify the code

### "Failed to push" error

- Verify your GitHub token is valid and has `repo` scope
- Ensure the repository exists and you have write access
- Check your internet connection
- Try regenerating the GitHub token

### "Repository is empty" error

The tool requires at least one commit in the repository. Initialize it:

```bash
# Clone your empty repo
git clone https://github.com/username/my-bookmarks.git
cd my-bookmarks

# Create initial commit
echo "# My Bookmarks" > README.md
git add README.md
git commit -m "Initial commit"
git push origin main
```

### Service not starting on boot

**Windows:**
- Open Task Scheduler and verify the "Bookmarked" task exists
- Check task properties and ensure trigger is set to "At log on"

**macOS:**
- Verify plist exists: `ls ~/Library/LaunchAgents/com.bookmarked.sync.plist`
- Load manually: `launchctl load ~/Library/LaunchAgents/com.bookmarked.sync.plist`

**Linux:**
- Check service file: `systemctl --user cat bookmarked`
- Enable if needed: `systemctl --user enable bookmarked`

## Development

### Project Structure

```
bookmarked-cli/
├── cmd/
│   └── bookmarked/
│       └── main.go              # CLI entry point with Cobra commands
├── internal/
│   ├── bookmarks/
│   │   └── bookmarks.go         # Chrome bookmark detection & formatting
│   ├── config/
│   │   └── config.go            # YAML configuration management
│   ├── watcher/
│   │   └── watcher.go           # File watching with debouncing
│   ├── sync/
│   │   └── sync.go              # Git operations (clone, commit, push)
│   └── service/
│       ├── service.go           # Main service logic
│       ├── install.go           # Platform dispatcher
│       ├── install_windows.go   # Windows Task Scheduler
│       ├── install_darwin.go    # macOS launchd
│       └── install_linux.go     # Linux systemd
├── .github/
│   └── workflows/
│       ├── ci.yml               # CI pipeline
│       └── release.yml          # Automated releases
├── go.mod
├── go.sum
├── README.md
├── LICENSE
├── Makefile
└── config.example.yaml
```

### Building from Source

```bash
# Clone the repository
git clone git@github.com:vivek-dodia/bookmarked-cli.git
cd bookmarked-cli

# Download dependencies
go mod download

# Build for current platform
go build -o bookmarked ./cmd/bookmarked

# Build for all platforms
make build-all

# Run tests
go test ./...
```

### Cross-Compilation

```bash
# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o bookmarked.exe ./cmd/bookmarked

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o bookmarked-mac-intel ./cmd/bookmarked

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o bookmarked-mac-arm ./cmd/bookmarked

# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o bookmarked-linux ./cmd/bookmarked
```

### Technologies Used

- **[Go](https://golang.org/)** - Programming language
- **[Cobra](https://github.com/spf13/cobra)** - CLI framework
- **[go-git](https://github.com/go-git/go-git)** - Pure Go Git implementation
- **[fsnotify](https://github.com/fsnotify/fsnotify)** - Cross-platform file system notifications
- **[yaml.v3](https://gopkg.in/yaml.v3)** - YAML parsing

## Security Considerations

- **Private Repository**: Always use a private GitHub repository for storing bookmarks
- **Token Security**:
  - Config file permissions are set to `0600` (user read/write only)
  - Never commit your config file to version control
  - Rotate tokens periodically
- **Minimal Permissions**: GitHub token only needs `repo` scope
- **Local Storage**: Bookmark data stored locally in `~/.bookmarked/repo/`

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [go-git](https://github.com/go-git/go-git) for Git operations
- Uses [fsnotify](https://github.com/fsnotify/fsnotify) for file watching
- CLI powered by [cobra](https://github.com/spf13/cobra)

## Support

If you encounter any issues or have questions:

1. Check the [Troubleshooting](#troubleshooting) section
2. Search [existing issues](https://github.com/vivek-dodia/bookmarked-cli/issues)
3. Create a [new issue](https://github.com/vivek-dodia/bookmarked-cli/issues/new) with:
   - OS and version
   - Go version (if building from source)
   - Error messages
   - Steps to reproduce

---

**Made for Chrome users who want their bookmarks backed up and version-controlled.**
