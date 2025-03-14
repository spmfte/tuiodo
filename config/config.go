package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Default config file locations
const (
	DefaultConfigFileName = "tuiodo.yaml"
)

// Config represents the application configuration
type Config struct {
	General     GeneralConfig     `yaml:"general"`
	UI          UIConfig          `yaml:"ui"`
	Colors      ColorsConfig      `yaml:"colors"`
	Keybindings KeybindingsConfig `yaml:"keybindings"`
	Storage     StorageConfig     `yaml:"storage"`
}

// GeneralConfig contains general application settings
type GeneralConfig struct {
	DefaultCategory string `yaml:"default_category"`
	ShowStatusBar   bool   `yaml:"show_status_bar"`
	TasksPerPage    int    `yaml:"tasks_per_page"`
	ClearStatus     int    `yaml:"clear_status_after_seconds"` // Seconds before clearing status message (0 = never)
}

// UIConfig contains UI-related settings
type UIConfig struct {
	ShowHeader      bool   `yaml:"show_header"`
	HeaderFormat    string `yaml:"header_format"`
	ShowCategories  bool   `yaml:"show_categories"`
	ShowPriorities  bool   `yaml:"show_priorities"`
	ShowDueDates    bool   `yaml:"show_due_dates"`
	TaskSeparator   string `yaml:"task_separator"`
	EnableTabs      bool   `yaml:"enable_tabs"`
	EnableBorders   bool   `yaml:"enable_borders"`
	BorderStyle     string `yaml:"border_style"` // rounded, normal, double, thick, none
	DateFormat      string `yaml:"date_format"`
	CursorIndicator string `yaml:"cursor_indicator"`
	CheckboxDone    string `yaml:"checkbox_done"`
	CheckboxPending string `yaml:"checkbox_pending"`
}

// ColorsConfig contains color schemes
type ColorsConfig struct {
	ColorMode      string            `yaml:"color_mode"` // auto, true_color, 256, 16, none
	Theme          string            `yaml:"theme"`      // default, dark, light, custom
	Primary        string            `yaml:"primary"`
	Secondary      string            `yaml:"secondary"`
	Tertiary       string            `yaml:"tertiary"`
	Success        string            `yaml:"success"`
	Warning        string            `yaml:"warning"`
	Error          string            `yaml:"error"`
	Text           string            `yaml:"text"`
	TextDim        string            `yaml:"text_dim"`
	TextMuted      string            `yaml:"text_muted"`
	Highlight      string            `yaml:"highlight"`
	Border         string            `yaml:"border"`
	BorderFocus    string            `yaml:"border_focus"`
	Subtle         string            `yaml:"subtle"`
	PriorityHigh   string            `yaml:"priority_high"`
	PriorityMedium string            `yaml:"priority_medium"`
	PriorityLow    string            `yaml:"priority_low"`
	TaskDone       string            `yaml:"task_done"`
	TaskPending    string            `yaml:"task_pending"`
	CategoryColors map[string]string `yaml:"category_colors"` // Map category names to specific colors
}

// KeybindingsConfig contains keyboard shortcuts
type KeybindingsConfig struct {
	QuitKey           []string `yaml:"quit"`
	AddTaskKey        []string `yaml:"add_task"`
	EditTaskKey       []string `yaml:"edit_task"`
	DeleteTaskKey     []string `yaml:"delete_task"`
	ToggleTaskKey     []string `yaml:"toggle_task"`
	CyclePriorityKey  []string `yaml:"cycle_priority"`
	CycleCategoryKey  []string `yaml:"cycle_category"`
	CycleTabKey       []string `yaml:"cycle_tab"`
	NextPageKey       []string `yaml:"next_page"`
	PrevPageKey       []string `yaml:"prev_page"`
	MoveCursorUpKey   []string `yaml:"move_cursor_up"`
	MoveCursorDownKey []string `yaml:"move_cursor_down"`
	HelpKey           []string `yaml:"help"`
}

// StorageConfig contains storage-related settings
type StorageConfig struct {
	FilePath        string `yaml:"file_path"`
	BackupDirectory string `yaml:"backup_directory"`
	AutoSave        bool   `yaml:"auto_save"`
	BackupOnSave    bool   `yaml:"backup_on_save"`
	MaxBackups      int    `yaml:"max_backups"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		General: GeneralConfig{
			DefaultCategory: "Uncategorized",
			ShowStatusBar:   true,
			TasksPerPage:    10,
			ClearStatus:     3,
		},
		UI: UIConfig{
			ShowHeader:      true,
			HeaderFormat:    "TUIODO",
			ShowCategories:  true,
			ShowPriorities:  true,
			ShowDueDates:    true,
			TaskSeparator:   "─",
			EnableTabs:      true,
			EnableBorders:   true,
			BorderStyle:     "rounded",
			DateFormat:      "2006-01-02",
			CursorIndicator: "→ ",
			CheckboxDone:    "[✓]",
			CheckboxPending: "[ ]",
		},
		Colors: ColorsConfig{
			ColorMode:      "auto",
			Theme:          "default",
			Primary:        "#7C3AED", // Purple
			Secondary:      "#2563EB", // Blue
			Tertiary:       "#10B981", // Green
			Success:        "#10B981", // Green
			Warning:        "#F59E0B", // Amber
			Error:          "#EF4444", // Red
			Text:           "#F9FAFB", // Nearly white
			TextDim:        "#E5E7EB", // Light gray
			TextMuted:      "#9CA3AF", // Medium gray
			Highlight:      "#C4B5FD", // Light purple
			Border:         "#4B5563", // Dark gray
			BorderFocus:    "#8B5CF6", // Medium purple
			Subtle:         "#374151", // Very dark gray
			PriorityHigh:   "#EF4444", // Red
			PriorityMedium: "#F59E0B", // Amber
			PriorityLow:    "#10B981", // Green
			TaskDone:       "#9CA3AF", // Medium gray
			TaskPending:    "#F9FAFB", // Nearly white
			CategoryColors: map[string]string{
				"Work":     "#3B82F6", // Blue
				"Personal": "#EC4899", // Pink
				"Health":   "#10B981", // Green
				"Finance":  "#6366F1", // Indigo
			},
		},
		Keybindings: KeybindingsConfig{
			QuitKey:           []string{"q", "ctrl+c"},
			AddTaskKey:        []string{"a"},
			EditTaskKey:       []string{"e"},
			DeleteTaskKey:     []string{"d"},
			ToggleTaskKey:     []string{"enter", "space"},
			CyclePriorityKey:  []string{"p"},
			CycleCategoryKey:  []string{"c"},
			CycleTabKey:       []string{"tab", "t"},
			NextPageKey:       []string{"right", "l", "n"},
			PrevPageKey:       []string{"left", "h", "b"},
			MoveCursorUpKey:   []string{"up", "k"},
			MoveCursorDownKey: []string{"down", "j"},
			HelpKey:           []string{"?", "F1"},
		},
		Storage: StorageConfig{
			FilePath:        "TODO.md",
			BackupDirectory: "~/.config/tuiodo/backups",
			AutoSave:        true,
			BackupOnSave:    true,
			MaxBackups:      5,
		},
	}
}

// LoadConfig loads the configuration from the specified file
// If path is empty, it will try to load from default locations
func LoadConfig(path string) (Config, error) {
	config := DefaultConfig()

	// If no path specified, try default locations
	if path == "" {
		// Always prioritize $HOME/.config location
		homeDir, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(homeDir, ".config", "tuiodo", DefaultConfigFileName)
			if _, err := os.Stat(path); err == nil {
				// File exists, try to load it
				return loadFromFile(path, config)
			}
		}

		// Try user config directory as fallback
		userConfigDir, err := os.UserConfigDir()
		if err == nil && userConfigDir != filepath.Join(homeDir, ".config") {
			path = filepath.Join(userConfigDir, "tuiodo", DefaultConfigFileName)
			if _, err := os.Stat(path); err == nil {
				// File exists, try to load it
				return loadFromFile(path, config)
			}
		}

		// Try current directory as a last resort
		path = DefaultConfigFileName
		if _, err := os.Stat(path); err == nil {
			// File exists, try to load it
			return loadFromFile(path, config)
		}

		// No config file found, return default config
		return config, nil
	}

	// Try to load from the specified path
	return loadFromFile(path, config)
}

// loadFromFile loads config from specified file path
func loadFromFile(path string, defaultConfig Config) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return defaultConfig, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return defaultConfig, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Merge with defaults to ensure all fields are set
	return mergeWithDefaults(config, defaultConfig), nil
}

// mergeWithDefaults ensures all required fields have values by filling in
// missing values from the default config
func mergeWithDefaults(config, defaults Config) Config {
	// Merge General section
	if config.General.DefaultCategory == "" {
		config.General.DefaultCategory = defaults.General.DefaultCategory
	}
	if config.General.TasksPerPage == 0 {
		config.General.TasksPerPage = defaults.General.TasksPerPage
	}

	// Merge UI section
	if config.UI.BorderStyle == "" {
		config.UI.BorderStyle = defaults.UI.BorderStyle
	}
	if config.UI.DateFormat == "" {
		config.UI.DateFormat = defaults.UI.DateFormat
	}
	if config.UI.CursorIndicator == "" {
		config.UI.CursorIndicator = defaults.UI.CursorIndicator
	}
	if config.UI.CheckboxDone == "" {
		config.UI.CheckboxDone = defaults.UI.CheckboxDone
	}
	if config.UI.CheckboxPending == "" {
		config.UI.CheckboxPending = defaults.UI.CheckboxPending
	}
	if config.UI.TaskSeparator == "" {
		config.UI.TaskSeparator = defaults.UI.TaskSeparator
	}

	// Merge Storage section
	if config.Storage.FilePath == "" {
		config.Storage.FilePath = defaults.Storage.FilePath
	}
	if config.Storage.BackupDirectory == "" {
		config.Storage.BackupDirectory = defaults.Storage.BackupDirectory
	}
	if config.Storage.MaxBackups == 0 {
		config.Storage.MaxBackups = defaults.Storage.MaxBackups
	}

	// Ensure proper color settings
	if config.Colors.ColorMode == "" {
		config.Colors.ColorMode = defaults.Colors.ColorMode
	}
	if config.Colors.Theme == "" {
		config.Colors.Theme = defaults.Colors.Theme
	}
	if config.Colors.Primary == "" {
		config.Colors.Primary = defaults.Colors.Primary
	}
	if config.Colors.Secondary == "" {
		config.Colors.Secondary = defaults.Colors.Secondary
	}
	if config.Colors.Tertiary == "" {
		config.Colors.Tertiary = defaults.Colors.Tertiary
	}
	if config.Colors.Success == "" {
		config.Colors.Success = defaults.Colors.Success
	}
	if config.Colors.Warning == "" {
		config.Colors.Warning = defaults.Colors.Warning
	}
	if config.Colors.Error == "" {
		config.Colors.Error = defaults.Colors.Error
	}
	if config.Colors.Text == "" {
		config.Colors.Text = defaults.Colors.Text
	}

	// Initialize category colors if missing
	if config.Colors.CategoryColors == nil {
		config.Colors.CategoryColors = make(map[string]string)
	}

	// Copy default category colors if not set in config
	for category, color := range defaults.Colors.CategoryColors {
		if _, exists := config.Colors.CategoryColors[category]; !exists {
			config.Colors.CategoryColors[category] = color
		}
	}

	return config
}

// SaveConfig saves the configuration to the specified file
func SaveConfig(config Config, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// SaveDefaultConfig saves the default configuration to the specified path
func SaveDefaultConfig(path string) error {
	return SaveConfig(DefaultConfig(), path)
}

// GetConfigFilePath returns the path to the configuration file
// It will create the directory if it doesn't exist
func GetConfigFilePath() (string, error) {
	// Always use $HOME/.config as the primary location
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config")

	// Create tuiodo config directory
	tuiodoConfigDir := filepath.Join(configDir, "tuiodo")
	if err := os.MkdirAll(tuiodoConfigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(tuiodoConfigDir, DefaultConfigFileName), nil
}
