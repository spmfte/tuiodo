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
	Version   = "1.1.3"
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

// printVersion prints version information
func printVersion() {
	fmt.Printf("TUIODO v%s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

// printUsage prints the help message
func printUsage() {
	fmt.Printf(`TUIODO - A Modern Terminal Task Manager v%s

Usage:
  tuiodo [options]

Options:
  -h, --help                    Show this help message
  -v, --version                 Show version information
  -c, --config <path>          Path to config file
  --create-default-config       Create default configuration file and exit
  --print-config               Print current configuration and exit
  -s, --storage <path>         Path to storage file (overrides config)
  -t, --tasks-per-page <num>   Number of tasks per page (overrides config)
  --debug                      Enable debug mode with detailed logging
  --no-mouse                   Disable mouse support
  --no-color                   Disable color output
  --backup-dir <path>          Set backup directory (overrides config)
  --max-backups <num>          Set maximum number of backups (overrides config)
  --no-auto-save              Disable auto-save feature
  --no-backup                  Disable backup on save
  --category <name>            Start with specific category filter
  --sort <field>              Initial sort field (priority|created|category)
  --view <type>               Initial view (all|pending|completed)

Examples:
  tuiodo                                    # Start with default settings
  tuiodo --config ~/.config/tuiodo.yaml     # Use custom config file
  tuiodo --storage ~/tasks.md               # Use specific storage file
  tuiodo --category Work                    # Start with Work category filter
  tuiodo --sort priority                    # Sort tasks by priority
  tuiodo --view pending                     # Show only pending tasks
  tuiodo --no-mouse --no-color             # Terminal-friendly mode

For more information and documentation:
  https://github.com/spmfte/tuiodo
`, Version)
}

func main() {
	// Parse command line flags
	flags := config.ParseFlags()

	// Handle version flag
	if flags.ShowVersion {
		printVersion()
		return
	}

	// Handle help flag
	if flags.ShowHelp {
		printUsage()
		return
	}

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

	// Validate flags
	if err := config.ValidateFlags(flags); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Handle config-related flags
	cfg, shouldExit := config.HandleConfigFlags(flags)
	if shouldExit {
		return
	}

	// Initialize default config file if it doesn't exist
	if err := config.InitConfigIfNeeded(); err != nil {
		fmt.Printf("Warning: Could not initialize config: %v\n", err)
	}

	// Set up storage path from config or flags
	var storagePath string
	if flags.StoragePath != "" {
		storagePath = flags.StoragePath
	} else {
		// No explicit storage path provided, check for git repository first
		storagePath = "" // Empty string will trigger git detection in storage.Initialize
	}

	// Handle backup directory configuration
	backupDir := cfg.Storage.BackupDirectory
	if flags.BackupDir != "" {
		backupDir = flags.BackupDir
	}
	if backupDir != "" {
		if expanded, err := config.ExpandPath(backupDir); err == nil {
			backupDir = expanded
		}
	}

	// Initialize storage with full configuration
	storage.Initialize(
		storagePath,
		backupDir,
		flags.MaxBackups,
		!flags.NoAutoSave,
		!flags.NoBackup,
	)

	// Load tasks from storage
	tasks := storage.LoadTasks()

	// Update UI styles based on configuration
	styles := config.GetStyles(cfg)
	if flags.NoColor {
		styles = config.GetMonochromeStyles()
	}

	// Pass the styles with category colors to update the UI
	ui.UpdateStyles(struct{ CategoryColors map[string]string }{
		CategoryColors: styles.CategoryColors,
	})

	// Get key bindings from config
	keyBindings := config.GetKeyBindings(cfg)

	// Set application info in the UI
	ui.SetAppInfo(Version, GitCommit, BuildTime)

	// Create initial model with configuration
	initialModel := model.NewModelWithConfig(
		tasks,
		flags.TasksPerPage,
		flags.Category,
		keyBindings,
	)

	// Apply initial view and sort if specified
	if flags.View != "" {
		switch flags.View {
		case "all":
			initialModel.CurrentView = model.TabAll
		case "pending":
			initialModel.CurrentView = model.TabPending
		case "completed":
			initialModel.CurrentView = model.TabCompleted
		}
	}

	// Set sort option
	if flags.Sort != "" {
		initialModel.SortTasks(model.SortType(flags.Sort))
	} else {
		// Default sort by priority
		initialModel.SortTasks(model.SortByPriority)
	}

	// Configure tea program options
	options := []tea.ProgramOption{}

	// Add alt screen by default
	options = append(options, tea.WithAltScreen())

	// Add mouse support if enabled and on suitable platform
	if !flags.NoMouse && runtime.GOOS != "windows" {
		options = append(options, tea.WithMouseCellMotion())
	}

	// Create application instance
	app := App{
		model: initialModel,
		cfg:   cfg,
	}

	// Create and run the program
	p := tea.NewProgram(app, options...)

	// Run the application
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
