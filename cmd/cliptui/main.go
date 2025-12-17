package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/dvd/cliptui/internal/clipboard"
	"github.com/dvd/cliptui/internal/config"
	"github.com/dvd/cliptui/internal/storage"
	"github.com/dvd/cliptui/internal/tui"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "cliptui",
	Short: "A beautiful terminal-based clipboard history manager",
	Long: `clipTUI is a modern, fast, and elegant clipboard history manager for Linux.
It watches your system clipboard in the background, stores every item locally,
and lets you browse, search, preview, and restore previous clipboard entries.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action: show TUI
		showTUI()
	},
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start the clipboard monitoring daemon",
	Long:  "Runs the clipboard monitor in the foreground, watching for clipboard changes.",
	Run: func(cmd *cobra.Command, args []string) {
		runDaemon()
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the clipboard history TUI",
	Long:  "Opens the terminal UI to browse, search, and restore clipboard history.",
	Run: func(cmd *cobra.Command, args []string) {
		showTUI()
	},
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all clipboard history",
	Long:  "Deletes all stored clipboard items from the database.",
	Run: func(cmd *cobra.Command, args []string) {
		clearHistory()
	},
}

func init() {
	cfg = config.Default()

	rootCmd.AddCommand(daemonCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(clearCmd)

	rootCmd.PersistentFlags().StringVar(&cfg.DBPath, "db", cfg.DBPath, "Database path")
	rootCmd.PersistentFlags().IntVar(&cfg.MaxItems, "max-items", cfg.MaxItems, "Maximum items to store")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runDaemon() {
	store, err := storage.New(cfg.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	monitor := clipboard.NewMonitor(store, time.Duration(cfg.PollInterval)*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down daemon...")
		cancel()
	}()

	fmt.Println("Starting clipboard monitor daemon...")
	fmt.Printf("Database: %s\n", cfg.DBPath)
	fmt.Printf("Poll interval: %dms\n", cfg.PollInterval)
	fmt.Println("Press Ctrl+C to stop.")

	if err := monitor.Start(ctx); err != nil && err != context.Canceled {
		fmt.Fprintf(os.Stderr, "Daemon error: %v\n", err)
		os.Exit(1)
	}
}

func showTUI() {
	store, err := storage.New(cfg.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	model, err := tui.New(store)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create TUI: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}

func clearHistory() {
	store, err := storage.New(cfg.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	if err := store.Clear(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to clear history: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Clipboard history cleared successfully.")
}
