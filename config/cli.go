package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// CLIFlags contains the parsed command-line flags
type CLIFlags struct {
	ConfigFile          string
	PrintConfig         bool
	CreateDefaultConfig bool
	StoragePath         string
	TasksPerPage        int
	ShowHelp            bool
	ShowVersion         bool
	Debug               bool
	NoMouse             bool
	NoColor             bool
	BackupDir           string
	MaxBackups          int
	NoAutoSave          bool
	NoBackup            bool
	Category            string
	Sort                string
	View                string
}

// ParseFlags parses command-line flags
func ParseFlags() CLIFlags {
	var flags CLIFlags

	// Define command-line flags
	flag.StringVar(&flags.ConfigFile, "config", "", "Path to config file")
	flag.StringVar(&flags.ConfigFile, "c", "", "Path to config file (shorthand)")

	flag.BoolVar(&flags.PrintConfig, "print-config", false, "Print current configuration and exit")
	flag.BoolVar(&flags.CreateDefaultConfig, "create-default-config", false, "Create default configuration file and exit")

	flag.StringVar(&flags.StoragePath, "storage", "", "Path to storage file (overrides config)")
	flag.StringVar(&flags.StoragePath, "s", "", "Path to storage file (shorthand)")

	flag.IntVar(&flags.TasksPerPage, "tasks-per-page", 0, "Number of tasks per page (overrides config)")
	flag.IntVar(&flags.TasksPerPage, "t", 0, "Number of tasks per page (shorthand)")

	flag.BoolVar(&flags.ShowHelp, "help", false, "Show help and exit")
	flag.BoolVar(&flags.ShowHelp, "h", false, "Show help and exit (shorthand)")

	flag.BoolVar(&flags.ShowVersion, "version", false, "Show version information")
	flag.BoolVar(&flags.ShowVersion, "v", false, "Show version information (shorthand)")

	flag.BoolVar(&flags.Debug, "debug", false, "Enable debug mode with detailed logging")
	flag.BoolVar(&flags.NoMouse, "no-mouse", false, "Disable mouse support")
	flag.BoolVar(&flags.NoColor, "no-color", false, "Disable color output")

	flag.StringVar(&flags.BackupDir, "backup-dir", "", "Set backup directory (overrides config)")
	flag.IntVar(&flags.MaxBackups, "max-backups", 5, "Set maximum number of backups (overrides config)")
	flag.BoolVar(&flags.NoAutoSave, "no-auto-save", false, "Disable auto-save feature")
	flag.BoolVar(&flags.NoBackup, "no-backup", false, "Disable backup on save")

	flag.StringVar(&flags.Category, "category", "", "Start with specific category filter")
	flag.StringVar(&flags.Sort, "sort", "", "Initial sort field (priority|created|category)")
	flag.StringVar(&flags.View, "view", "", "Initial view (all|pending|completed)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of tuiodo:\n")
		fmt.Fprintf(os.Stderr, "  tuiodo [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  tuiodo --config ~/.config/tuiodo/my-config.yaml\n")
		fmt.Fprintf(os.Stderr, "  tuiodo --create-default-config\n")
		fmt.Fprintf(os.Stderr, "  tuiodo --storage ~/my-tasks.md\n")
		fmt.Fprintf(os.Stderr, "  tuiodo --category Work\n")
		fmt.Fprintf(os.Stderr, "  tuiodo --sort priority\n")
		fmt.Fprintf(os.Stderr, "  tuiodo --view pending\n")
		fmt.Fprintf(os.Stderr, "  tuiodo --no-mouse --no-color\n")
	}

	// Parse flags
	flag.Parse()

	return flags
}

// HandleConfigFlags processes the CLI flags related to configuration
// Returns the loaded configuration and a boolean indicating if the program should exit
func HandleConfigFlags(flags CLIFlags) (Config, bool) {
	// Handle help flag
	if flags.ShowHelp {
		flag.Usage()
		return Config{}, true
	}

	// Handle create default config flag
	if flags.CreateDefaultConfig {
		configPath, err := GetConfigFilePath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error determining config path: %v\n", err)
			return Config{}, true
		}

		// Check if file already exists and confirm overwrite
		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("Config file already exists at %s\n", configPath)
			fmt.Print("Overwrite? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Aborting.")
				return Config{}, true
			}
		}

		if err := SaveDefaultConfig(configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating default config: %v\n", err)
			return Config{}, true
		}

		fmt.Printf("Created default config file at %s\n", configPath)
		return Config{}, true
	}

	// Load configuration
	config, err := LoadConfig(flags.ConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		return Config{}, true
	}

	// Handle print config flag
	if flags.PrintConfig {
		configData, err := yaml.Marshal(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error serializing config: %v\n", err)
			return Config{}, true
		}

		fmt.Println(string(configData))
		return Config{}, true
	}

	// Apply CLI overrides
	if flags.StoragePath != "" {
		config.Storage.FilePath = flags.StoragePath
	}

	if flags.TasksPerPage > 0 {
		config.General.TasksPerPage = flags.TasksPerPage
	}

	if flags.BackupDir != "" {
		config.Storage.BackupDirectory = flags.BackupDir
	}

	if flags.MaxBackups > 0 {
		config.Storage.MaxBackups = flags.MaxBackups
	}

	return config, false
}

// InitConfigIfNeeded creates a default config file if none exists
func InitConfigIfNeeded() error {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return fmt.Errorf("failed to determine config path: %w", err)
	}

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil {
		// Config exists, nothing to do
		return nil
	}

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Save default config
	if err := SaveDefaultConfig(configPath); err != nil {
		return fmt.Errorf("failed to save default config: %w", err)
	}

	return nil
}
