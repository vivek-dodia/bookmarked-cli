#!/bin/bash
# Quick installation script for bookmarked

set -e

echo "=== Bookmarked Installation ==="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed."
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

# Build the binary
echo "Building bookmarked..."
go build -o bookmarked ./cmd/bookmarked

# Move to bin directory
if [ -w "/usr/local/bin" ]; then
    echo "Installing to /usr/local/bin..."
    mv bookmarked /usr/local/bin/
elif [ -d "$HOME/.local/bin" ]; then
    echo "Installing to $HOME/.local/bin..."
    mv bookmarked "$HOME/.local/bin/"
    echo "Make sure $HOME/.local/bin is in your PATH"
else
    echo "Installing to current directory..."
    echo "You can manually move 'bookmarked' to a directory in your PATH"
fi

echo ""
echo "✓ Installation complete!"
echo ""
echo "Next steps:"
echo "1. Run: bookmarked init"
echo "2. Edit config file with your GitHub repo and token"
echo "3. Run: bookmarked sync (to test)"
echo "4. Run: bookmarked install (to set up automatic syncing)"
echo ""
