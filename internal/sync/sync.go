package sync

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/yourusername/bookmarked/internal/config"
)

type GitSync struct {
	cfg      *config.Config
	repoPath string
	repo     *git.Repository
}

// New creates a new GitSync instance
func New(cfg *config.Config) (*GitSync, error) {
	repoPath, err := config.GetRepoPath()
	if err != nil {
		return nil, err
	}

	return &GitSync{
		cfg:      cfg,
		repoPath: repoPath,
	}, nil
}

// Initialize clones the repository or opens it if it already exists
func (gs *GitSync) Initialize() error {
	// Check if repo already exists
	if _, err := os.Stat(gs.repoPath); err == nil {
		// Repo exists, open it
		repo, err := git.PlainOpen(gs.repoPath)
		if err != nil {
			return fmt.Errorf("failed to open repository: %w", err)
		}
		gs.repo = repo
		log.Println("Opened existing repository")
		return nil
	}

	// Clone the repository
	log.Printf("Cloning repository: %s", gs.cfg.GitHubRepo)

	repoURL := fmt.Sprintf("https://github.com/%s.git", gs.cfg.GitHubRepo)

	repo, err := git.PlainClone(gs.repoPath, false, &git.CloneOptions{
		URL: repoURL,
		Auth: &http.BasicAuth{
			Username: "git", // can be anything
			Password: gs.cfg.GitHubToken,
		},
		Progress: os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	gs.repo = repo
	log.Println("Repository cloned successfully")
	return nil
}

// CommitAndPush commits changes and pushes to remote
func (gs *GitSync) CommitAndPush(message string) error {
	if gs.repo == nil {
		return fmt.Errorf("repository not initialized")
	}

	// Get the worktree
	w, err := gs.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Check status to see if there are changes
	status, err := w.Status()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	if status.IsClean() {
		log.Println("No changes to commit")
		return nil
	}

	// Add all changes
	_, err = w.Add(".")
	if err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	// Commit
	commit, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Bookmarked",
			Email: "bookmarked@local",
			When:  time.Now(),
		},
	})

	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	log.Printf("Created commit: %s", commit.String())

	// Push to remote
	log.Println("Pushing to remote...")
	err = gs.repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "git",
			Password: gs.cfg.GitHubToken,
		},
	})

	if err != nil {
		// Check if error is "already up-to-date"
		if err == git.NoErrAlreadyUpToDate {
			log.Println("Already up to date")
			return nil
		}
		return fmt.Errorf("failed to push: %w", err)
	}

	log.Println("Pushed successfully")
	return nil
}

// Pull fetches and merges changes from remote
func (gs *GitSync) Pull() error {
	if gs.repo == nil {
		return fmt.Errorf("repository not initialized")
	}

	w, err := gs.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = w.Pull(&git.PullOptions{
		Auth: &http.BasicAuth{
			Username: "git",
			Password: gs.cfg.GitHubToken,
		},
	})

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Println("Already up to date")
			return nil
		}
		return fmt.Errorf("failed to pull: %w", err)
	}

	log.Println("Pulled successfully")
	return nil
}

// GetRepoPath returns the local repository path
func (gs *GitSync) GetRepoPath() string {
	return gs.repoPath
}
