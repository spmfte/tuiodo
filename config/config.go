package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
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
	Display     DisplayConfig     `yaml:"display"`
	Files       FilesConfig       `yaml:"files"`
	Sort        SortConfig        `yaml:"sort"`
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

// ColorsConfig contains color-related settings
type ColorsConfig struct {
	Theme          string            `yaml:"theme"`
	ColorMode      string            `yaml:"color_mode"`
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
	Background     string            `yaml:"background"`
	PriorityHigh   string            `yaml:"priority_high"`
	PriorityMedium string            `yaml:"priority_medium"`
	PriorityLow    string            `yaml:"priority_low"`
	TaskDone       string            `yaml:"task_done"`
	TaskPending    string            `yaml:"task_pending"`
	CategoryColors map[string]string `yaml:"category_colors"`
}

// KeybindingsConfig contains keybinding settings
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

// DisplayConfig contains display-related settings
type DisplayConfig struct {
	ShowDates bool `yaml:"show_dates"` // Whether to show creation dates for tasks
}

// FilesConfig contains file-related settings
type FilesConfig struct {
	GlobalTodoFile    string   `yaml:"global_todo_file"`    // Path to global todo file
	DirectoryTodoFile string   `yaml:"directory_todo_file"` // Name of directory-specific todo files
	ActiveFile        string   `yaml:"active_file"`         // Currently active todo file
	ExcludeDirs       []string `yaml:"exclude_dirs"`        // Directories to exclude from todo file search
}

// SortConfig contains sorting-related settings
type SortConfig struct {
	Field     string `yaml:"field"`     // Field to sort by (priority, date, category)
	Direction string `yaml:"direction"` // Sort direction (asc, desc)
}

// Colors represents the theme colors configuration
type Colors struct {
	Primary     string `yaml:"primary"`
	Secondary   string `yaml:"secondary"`
	Tertiary    string `yaml:"tertiary"`
	Success     string `yaml:"success"`
	Warning     string `yaml:"warning"`
	Error       string `yaml:"error"`
	Critical    string `yaml:"critical"`
	Text        string `yaml:"text"`
	TextDim     string `yaml:"text_dim"`
	TextMuted   string `yaml:"text_muted"`
	Highlight   string `yaml:"highlight"`
	Border      string `yaml:"border"`
	BorderFocus string `yaml:"border_focus"`
	Subtle      string `yaml:"subtle"`
	Background  string `yaml:"background"`
}

// DefaultColors returns the default color theme
func DefaultColors() Colors {
	return Colors{
		Primary:     "#7C3AED", // Purple
		Secondary:   "#2563EB", // Blue
		Tertiary:    "#10B981", // Green
		Success:     "#10B981", // Green
		Warning:     "#F59E0B", // Amber
		Error:       "#EF4444", // Red
		Critical:    "#991B1B", // Dark Red
		Text:        "#F9FAFB", // Nearly white
		TextDim:     "#E5E7EB", // Light gray
		TextMuted:   "#9CA3AF", // Medium gray
		Highlight:   "#C4B5FD", // Light purple
		Border:      "#4B5563", // Dark gray
		BorderFocus: "#8B5CF6", // Medium purple
		Subtle:      "#374151", // Very dark gray
		Background:  "#1F2937", // Dark blue-gray
	}
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "tuiodo")

	return Config{
		General: GeneralConfig{
			DefaultCategory: "Uncategorized",
			ShowStatusBar:   true,
			TasksPerPage:    10,
			ClearStatus:     5,
		},
		Storage: StorageConfig{
			FilePath:        filepath.Join(homeDir, "TODO.md"),
			BackupDirectory: filepath.Join(configDir, "backups"),
			AutoSave:        true,
			BackupOnSave:    true,
			MaxBackups:      5,
		},
		UI: UIConfig{
			ShowHeader:      true,
			HeaderFormat:    "TUIODO",
			ShowCategories:  true,
			ShowPriorities:  true,
			EnableTabs:      true,
			EnableBorders:   true,
			BorderStyle:     "rounded",
			DateFormat:      "2006-01-02",
			CursorIndicator: "→ ",
			CheckboxDone:    "[✓]",
			CheckboxPending: "[ ]",
			TaskSeparator:   "─",
		},
		Colors: ColorsConfig{
			Theme:          "default",
			ColorMode:      "auto",
			Primary:        "#7C3AED",
			Secondary:      "#2563EB",
			Tertiary:       "#10B981",
			Success:        "#10B981",
			Warning:        "#F59E0B",
			Error:          "#EF4444",
			Text:           "#F9FAFB",
			TextDim:        "#E5E7EB",
			TextMuted:      "#9CA3AF",
			Highlight:      "#C4B5FD",
			Border:         "#4B5563",
			BorderFocus:    "#8B5CF6",
			Subtle:         "#374151",
			Background:     "#1F2937",
			PriorityHigh:   "#DC2626",
			PriorityMedium: "#F59E0B",
			PriorityLow:    "#10B981",
			TaskDone:       "#6B7280",
			TaskPending:    "#F9FAFB",
			CategoryColors: map[string]string{
				"Work":     "#3B82F6",
				"Personal": "#EC4899",
				"Health":   "#10B981",
				"Finance":  "#6366F1",
			},
		},
		Display: DisplayConfig{
			ShowDates: true,
		},
		Files: FilesConfig{
			GlobalTodoFile:    filepath.Join(homeDir, "TODO.md"),
			DirectoryTodoFile: "TODO.md",
			ExcludeDirs:       []string{".git", "node_modules"},
		},
		Sort: SortConfig{
			Field:     "priority",
			Direction: "desc",
		},
		Keybindings: KeybindingsConfig{
			QuitKey:           []string{"q", "ctrl+c"},
			AddTaskKey:        []string{"a"},
			EditTaskKey:       []string{"e"},
			DeleteTaskKey:     []string{"d"},
			ToggleTaskKey:     []string{"space", "enter"},
			CyclePriorityKey:  []string{"p"},
			CycleCategoryKey:  []string{"c"},
			CycleTabKey:       []string{"tab", "t"},
			NextPageKey:       []string{"n", "right"},
			PrevPageKey:       []string{"b", "left"},
			MoveCursorUpKey:   []string{"k", "up"},
			MoveCursorDownKey: []string{"j", "down"},
			HelpKey:           []string{"?", "h", "f1"},
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
	if !config.General.ShowStatusBar {
		config.General.ShowStatusBar = defaults.General.ShowStatusBar
	}
	if config.General.ClearStatus == 0 {
		config.General.ClearStatus = defaults.General.ClearStatus
	}

	// Merge UI section
	if !config.UI.ShowHeader {
		config.UI.ShowHeader = defaults.UI.ShowHeader
	}
	if config.UI.HeaderFormat == "" {
		config.UI.HeaderFormat = defaults.UI.HeaderFormat
	}
	if !config.UI.ShowCategories {
		config.UI.ShowCategories = defaults.UI.ShowCategories
	}
	if !config.UI.ShowPriorities {
		config.UI.ShowPriorities = defaults.UI.ShowPriorities
	}
	if !config.UI.EnableTabs {
		config.UI.EnableTabs = defaults.UI.EnableTabs
	}
	if !config.UI.EnableBorders {
		config.UI.EnableBorders = defaults.UI.EnableBorders
	}
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

	// Merge Colors section
	if config.Colors.Theme == "" {
		config.Colors.Theme = defaults.Colors.Theme
	}
	if config.Colors.ColorMode == "" {
		config.Colors.ColorMode = defaults.Colors.ColorMode
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
	if config.Colors.TextDim == "" {
		config.Colors.TextDim = defaults.Colors.TextDim
	}
	if config.Colors.TextMuted == "" {
		config.Colors.TextMuted = defaults.Colors.TextMuted
	}
	if config.Colors.Highlight == "" {
		config.Colors.Highlight = defaults.Colors.Highlight
	}
	if config.Colors.Border == "" {
		config.Colors.Border = defaults.Colors.Border
	}
	if config.Colors.BorderFocus == "" {
		config.Colors.BorderFocus = defaults.Colors.BorderFocus
	}
	if config.Colors.Subtle == "" {
		config.Colors.Subtle = defaults.Colors.Subtle
	}
	if config.Colors.Background == "" {
		config.Colors.Background = defaults.Colors.Background
	}
	if config.Colors.PriorityHigh == "" {
		config.Colors.PriorityHigh = defaults.Colors.PriorityHigh
	}
	if config.Colors.PriorityMedium == "" {
		config.Colors.PriorityMedium = defaults.Colors.PriorityMedium
	}
	if config.Colors.PriorityLow == "" {
		config.Colors.PriorityLow = defaults.Colors.PriorityLow
	}
	if config.Colors.TaskDone == "" {
		config.Colors.TaskDone = defaults.Colors.TaskDone
	}
	if config.Colors.TaskPending == "" {
		config.Colors.TaskPending = defaults.Colors.TaskPending
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

	// Merge Display section
	if !config.Display.ShowDates {
		config.Display.ShowDates = defaults.Display.ShowDates
	}

	// Merge Files section
	if config.Files.GlobalTodoFile == "" {
		config.Files.GlobalTodoFile = defaults.Files.GlobalTodoFile
	}
	if config.Files.DirectoryTodoFile == "" {
		config.Files.DirectoryTodoFile = defaults.Files.DirectoryTodoFile
	}
	if len(config.Files.ExcludeDirs) == 0 {
		config.Files.ExcludeDirs = defaults.Files.ExcludeDirs
	}

	// Merge Sort section
	if config.Sort.Field == "" {
		config.Sort.Field = defaults.Sort.Field
	}
	if config.Sort.Direction == "" {
		config.Sort.Direction = defaults.Sort.Direction
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

// GetCategoryStyle returns the style for a category
func GetCategoryStyle(cfg Config, category string) lipgloss.Style {
	style := lipgloss.NewStyle().
		Italic(true).
		Padding(0, 1)

	if color, ok := cfg.Colors.CategoryColors[category]; ok {
		style = style.Foreground(lipgloss.Color(color))
	} else {
		style = style.Foreground(lipgloss.Color(cfg.Colors.Secondary))
	}

	return style
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "tuiodo.yaml"
	}
	return filepath.Join(homeDir, ".config", "tuiodo", "config.yaml")
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig() error {
	path := GetConfigPath()
	cfg := DefaultConfig()
	return SaveConfig(cfg, path)
}
