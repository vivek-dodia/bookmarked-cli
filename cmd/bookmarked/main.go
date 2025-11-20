package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/bookmarked/internal/config"
	"github.com/yourusername/bookmarked/internal/service"
)

var (
	version = "dev"
	cfgFile string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "bookmarked",
	Short: "Sync Chrome bookmarks to GitHub automatically",
	Long:  `A minimal background service that watches your Chrome bookmarks and syncs them to a GitHub repository.`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.InitConfig()
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Manually trigger a sync",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		svc := service.New(cfg)
		return svc.SyncOnce()
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the background sync service",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		svc := service.New(cfg)
		return svc.Start()
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install as a background service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return service.Install()
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the background service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return service.Uninstall()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check service status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return service.Status()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(statusCmd)
}
