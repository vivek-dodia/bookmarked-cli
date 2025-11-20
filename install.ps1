# Quick installation script for bookmarked (Windows)

Write-Host "=== Bookmarked Installation ===" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Error: Go is not installed." -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/"
    exit 1
}

# Build the binary
Write-Host "Building bookmarked..." -ForegroundColor Yellow
go build -o bookmarked.exe .\cmd\bookmarked

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Build successful!" -ForegroundColor Green
} else {
    Write-Host "✗ Build failed!" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Binary created: bookmarked.exe" -ForegroundColor Green
Write-Host ""
Write-Host "To add to PATH, move bookmarked.exe to a directory in your PATH" -ForegroundColor Yellow
Write-Host "Or keep it in the current directory and use: .\bookmarked.exe" -ForegroundColor Yellow
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "1. Run: .\bookmarked.exe init"
Write-Host "2. Edit config file with your GitHub repo and token"
Write-Host "3. Run: .\bookmarked.exe sync (to test)"
Write-Host "4. Run: .\bookmarked.exe install (to set up automatic syncing)"
Write-Host ""
