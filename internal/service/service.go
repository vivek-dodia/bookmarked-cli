package service

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vivek-dodia/bookmarked-cli/internal/bookmarks"
	"github.com/vivek-dodia/bookmarked-cli/internal/config"
	"github.com/vivek-dodia/bookmarked-cli/internal/sync"
	"github.com/vivek-dodia/bookmarked-cli/internal/watcher"
)

type Service struct {
	cfg          *config.Config
	gitSync      *sync.GitSync
	bookmarkPath string
	watcher      *watcher.Watcher
}

// New creates a new Service instance
func New(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,
	}
}

// Start runs the background sync service
func (s *Service) Start() error {
	// Set up logging
	if s.cfg.LogPath != "" {
		logFile, err := os.OpenFile(s.cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	log.Println("=== Bookmarked Service Starting ===")
	log.Printf("Time: %s", time.Now().Format(time.RFC3339))

	// Get bookmark path
	bookmarkPath, err := bookmarks.GetBookmarkPath()
	if err != nil {
		return fmt.Errorf("failed to get bookmark path: %w", err)
	}
	s.bookmarkPath = bookmarkPath
	log.Printf("Chrome bookmarks: %s", bookmarkPath)

	// Initialize git sync
	gitSync, err := sync.New(s.cfg)
	if err != nil {
		return fmt.Errorf("failed to create git sync: %w", err)
	}
	s.gitSync = gitSync

	// Initialize repository (clone or open)
	if err := s.gitSync.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Do initial sync
	log.Println("Performing initial sync...")
	if err := s.performSync(); err != nil {
		log.Printf("Initial sync failed: %v", err)
	}

	// Set up file watcher
	w, err := watcher.New(s.cfg.DebounceMs, func() {
		if err := s.performSync(); err != nil {
			log.Printf("Sync failed: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	s.watcher = w

	if err := w.Watch(bookmarkPath); err != nil {
		return fmt.Errorf("failed to watch bookmarks: %w", err)
	}

	log.Println("Service started successfully, watching for changes...")
	fmt.Println("✓ Bookmarked service is running")
	fmt.Println("Press Ctrl+C to stop")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
	s.watcher.Close()
	log.Println("Service stopped")

	return nil
}

// SyncOnce performs a one-time sync
func (s *Service) SyncOnce() error {
	log.Println("=== Manual Sync ===")

	// Get bookmark path
	bookmarkPath, err := bookmarks.GetBookmarkPath()
	if err != nil {
		return fmt.Errorf("failed to get bookmark path: %w", err)
	}
	s.bookmarkPath = bookmarkPath

	// Initialize git sync
	gitSync, err := sync.New(s.cfg)
	if err != nil {
		return fmt.Errorf("failed to create git sync: %w", err)
	}
	s.gitSync = gitSync

	// Initialize repository
	if err := s.gitSync.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Perform sync
	return s.performSync()
}

// performSync executes the sync operation
func (s *Service) performSync() error {
	startTime := time.Now()
	log.Println("--- Sync Starting ---")

	// Pull latest changes first
	if err := s.gitSync.Pull(); err != nil {
		log.Printf("Warning: Pull failed: %v", err)
	}

	// Copy and format bookmarks to repo
	if err := bookmarks.CopyToRepo(s.bookmarkPath, s.gitSync.GetRepoPath()); err != nil {
		return fmt.Errorf("failed to copy bookmarks: %w", err)
	}

	// Commit and push
	commitMsg := fmt.Sprintf("%s - %s", s.cfg.CommitMessage, time.Now().Format(time.RFC3339))
	if err := s.gitSync.CommitAndPush(commitMsg); err != nil {
		return fmt.Errorf("failed to commit and push: %w", err)
	}

	duration := time.Since(startTime)
	log.Printf("--- Sync Complete (took %v) ---", duration)
	return nil
}
