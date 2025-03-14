package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spmfte/tuiodo/config"
	"github.com/spmfte/tuiodo/handlers"
	"github.com/spmfte/tuiodo/model"
	"github.com/spmfte/tuiodo/storage"
	"github.com/spmfte/tuiodo/ui"
)

// Version information
var (
	Version   = "1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// App is the main application model
type App struct {
	model model.Model
	cfg   config.Config
}

// Init initializes the application
func (a App) Init() tea.Cmd {
	return nil
}

// Update processes messages and updates the model
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.model, cmd = handlers.Update(msg, a.model)
	return a, cmd
}

// View renders the application UI
func (a App) View() string {
	return ui.View(a.model)
}

func main() {
	// Get version info if built with -ldflags
	if buildInfo, ok := debug.ReadBuildInfo(); ok && GitCommit == "unknown" {
		for _, setting := range buildInfo.Settings {
			if setting.Key == "vcs.revision" {
				GitCommit = setting.Value[:7] // short hash
			} else if setting.Key == "vcs.time" {
				BuildTime = setting.Value
			}
		}
	}

	// Parse command line flags
	flags := config.ParseFlags()

	// Handle config-related flags (like --help, --create-default-config)
	cfg, shouldExit := config.HandleConfigFlags(flags)
	if shouldExit {
		return
	}

	// Initialize default config file if it doesn't exist
	// This is silent and won't overwrite existing config
	if err := config.InitConfigIfNeeded(); err != nil {
		fmt.Printf("Warning: Could not initialize config: %v\n", err)
	}

	// Set up storage path from config
	var storagePath string
	if flags.StoragePath != "" {
		// CLI flag overrides config
		storagePath = flags.StoragePath
	} else {
		// Otherwise use config value
		storagePath = cfg.Storage.FilePath

		// Expand ~ in path if needed
		expandedPath, err := config.ExpandPath(storagePath)
		if err == nil {
			storagePath = expandedPath
		}
	}

	// Expand backup directory path if needed
	backupDir := cfg.Storage.BackupDirectory
	if backupDir != "" {
		expandedBackupDir, err := config.ExpandPath(backupDir)
		if err == nil {
			backupDir = expandedBackupDir
		}
	}

	// Initialize storage with full config
	storage.Initialize(
		storagePath,
		backupDir,
		cfg.Storage.MaxBackups,
		cfg.Storage.AutoSave,
		cfg.Storage.BackupOnSave,
	)

	// Load tasks from storage
	tasks := storage.LoadTasks()

	// Update UI styles based on configuration
	styles := config.GetStyles(cfg)
	ui.UpdateStyles(styles)

	// Get key bindings from config
	keyBindings := config.GetKeyBindings(cfg)

	// Set application info in the UI
	ui.SetAppInfo(Version, GitCommit, BuildTime)

	// Create initial model with configuration
	initialModel := model.NewModelWithConfig(
		tasks,
		cfg.General.TasksPerPage,
		cfg.General.DefaultCategory,
		keyBindings,
	)

	// Create application instance
	app := App{
		model: initialModel,
		cfg:   cfg,
	}

	// Configure tea program options
	options := []tea.ProgramOption{tea.WithAltScreen()}

	// Add mouse support if on a suitable platform
	if runtime.GOOS != "windows" {
		options = append(options, tea.WithMouseCellMotion())
	}

	// Create and run the program
	p := tea.NewProgram(app, options...)

	// Run the application
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
