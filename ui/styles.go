package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Colors defines the color palette for the application
type Colors struct {
	Primary     lipgloss.Color
	Secondary   lipgloss.Color
	Tertiary    lipgloss.Color
	Success     lipgloss.Color
	Warning     lipgloss.Color
	Error       lipgloss.Color
	Critical    lipgloss.Color
	Text        lipgloss.Color
	TextDim     lipgloss.Color
	TextMuted   lipgloss.Color
	Highlight   lipgloss.Color
	Border      lipgloss.Color
	BorderFocus lipgloss.Color
	Subtle      lipgloss.Color
	Background  lipgloss.Color

	// Category colors (dynamic map)
	CategoryColors map[string]lipgloss.Color
}

// Global variables to store the current styles and colors
var (
	currentColors Colors
	currentStyles map[string]lipgloss.Style

	// App info
	appVersion   string
	appCommit    string
	appBuildTime string
)

// SetAppInfo sets the application version information
func SetAppInfo(version, commit, buildTime string) {
	appVersion = version
	appCommit = commit
	appBuildTime = buildTime
}

// GetAppInfo returns the application version information
func GetAppInfo() (string, string, string) {
	return appVersion, appCommit, appBuildTime
}

// AppColors returns the color palette for the application with a modern aesthetic
func AppColors() Colors {
	return Colors{
		Primary:     lipgloss.Color("#7C3AED"), // Purple
		Secondary:   lipgloss.Color("#2563EB"), // Blue
		Tertiary:    lipgloss.Color("#10B981"), // Green
		Success:     lipgloss.Color("#10B981"), // Green
		Warning:     lipgloss.Color("#F59E0B"), // Amber
		Error:       lipgloss.Color("#EF4444"), // Red
		Critical:    lipgloss.Color("#991B1B"), // Dark Red
		Text:        lipgloss.Color("#F9FAFB"), // Nearly white
		TextDim:     lipgloss.Color("#E5E7EB"), // Light gray
		TextMuted:   lipgloss.Color("#9CA3AF"), // Medium gray
		Highlight:   lipgloss.Color("#C4B5FD"), // Light purple
		Border:      lipgloss.Color("#4B5563"), // Dark gray
		BorderFocus: lipgloss.Color("#8B5CF6"), // Medium purple
		Subtle:      lipgloss.Color("#374151"), // Very dark gray
		Background:  lipgloss.Color("#1F2937"), // Dark blue-gray

		// Initialize with some default category colors
		CategoryColors: map[string]lipgloss.Color{
			"ui":            lipgloss.Color("#8B5CF6"), // Purple for UI tasks
			"add-task":      lipgloss.Color("#EC4899"), // Pink for add-task
			"bug":           lipgloss.Color("#EF4444"), // Red for bugs
			"function":      lipgloss.Color("#10B981"), // Green for function
			"fix":           lipgloss.Color("#F59E0B"), // Amber for fixes
			"functionality": lipgloss.Color("#3B82F6"), // Blue for functionality
			"layout":        lipgloss.Color("#6366F1"), // Indigo for layout
			"docs":          lipgloss.Color("#2563EB"), // Blue for docs
			"storage":       lipgloss.Color("#14B8A6"), // Teal for storage
		},
	}
}

// init initializes the styles when the package is loaded
func init() {
	// Setup default colors and styles
	currentColors = AppColors()
	currentStyles = CreateStyles(currentColors)
}

// GetStyle returns a style by name
func GetStyle(name string) lipgloss.Style {
	if style, ok := currentStyles[name]; ok {
		return style
	}
	return lipgloss.NewStyle() // Return a default style if not found
}

// UpdateStyles updates the global styles with the provided configuration
func UpdateStyles(styles interface{}) {
	// Default colors
	colors := AppColors()

	// If we have a config.Styles object, use its colors
	if cfg, ok := styles.(struct{ CategoryColors map[string]string }); ok && cfg.CategoryColors != nil {
		// Update category colors based on config
		for category, colorStr := range cfg.CategoryColors {
			colors.CategoryColors[strings.ToLower(category)] = lipgloss.Color(colorStr)
		}
	}

	// Create styles from updated colors
	currentColors = colors
	currentStyles = CreateStyles(colors)

	// Generate dynamic category styles based on the color map
	baseStyle := lipgloss.NewStyle().Italic(true).Padding(0, 1)
	for category, color := range colors.CategoryColors {
		categoryKey := "category_" + strings.ToLower(category)
		currentStyles[categoryKey] = baseStyle.Copy().Foreground(color)
	}
}

// CreateStyles returns the styles for the application
func CreateStyles(colors Colors) map[string]lipgloss.Style {
	// Define some common dimensions
	roundedBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	styles := map[string]lipgloss.Style{
		// App container
		"app": lipgloss.NewStyle().
			Padding(1, 2),

		// Title bar at the top
		"titleBar": lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Primary).
			Padding(0, 1).
			Margin(0, 0, 1, 0).
			Border(lipgloss.Border{Bottom: "─"}, false, false, true, false).
			BorderForeground(colors.Border),

		// Main title
		"title": lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Text).
			Background(colors.Primary).
			Padding(0, 1).
			MarginRight(1),

		// Version badge
		"versionBadge": lipgloss.NewStyle().
			Foreground(colors.Highlight).
			Padding(0, 1).
			MarginRight(1),

		// Secondary header text
		"secondary": lipgloss.NewStyle().
			Foreground(colors.Secondary).
			Bold(true),

		// Filter indicator
		"filterIndicator": lipgloss.NewStyle().
			Foreground(colors.Secondary).
			Italic(true).
			Padding(0, 1),

		// Help menu bar
		"helpBar": lipgloss.NewStyle().
			Foreground(colors.TextMuted).
			Padding(0, 1).
			Margin(1, 0),

		// Help box for help screen
		"helpBox": lipgloss.NewStyle().
			Border(roundedBorder).
			BorderForeground(colors.BorderFocus).
			Padding(1, 2),

		// Help command
		"helpCommand": lipgloss.NewStyle().
			Foreground(colors.Primary).
			Bold(true),

		// Task list container
		"listContainer": lipgloss.NewStyle().
			Border(roundedBorder).
			BorderForeground(colors.Border).
			Padding(1, 1),

		// Task header
		"taskHeader": lipgloss.NewStyle().
			Foreground(colors.TextDim).
			Bold(true).
			MarginBottom(1),

		// Input area
		"inputBox": lipgloss.NewStyle().
			Border(roundedBorder).
			BorderForeground(colors.BorderFocus).
			Padding(1, 1),

		// Input prompt
		"inputPrompt": lipgloss.NewStyle().
			Foreground(colors.Primary).
			Bold(true),

		// Input hint
		"inputHint": lipgloss.NewStyle().
			Foreground(colors.TextMuted).
			Italic(true),

		// Input cursor
		"inputCursor": lipgloss.NewStyle().
			Foreground(colors.Secondary),

		// Input field
		"input": lipgloss.NewStyle().
			Foreground(colors.Text),

		// Empty message
		"emptyMessage": lipgloss.NewStyle().
			Foreground(colors.TextMuted).
			Italic(true).
			Align(lipgloss.Center).
			Padding(1, 0),

		// Section divider
		"divider": lipgloss.NewStyle().
			Foreground(colors.Border),

		// Cursor indicator
		"cursor": lipgloss.NewStyle().
			Foreground(colors.Highlight).
			Bold(true),

		// Checkbox styles
		"checkboxPending": lipgloss.NewStyle().
			Foreground(colors.Warning),

		"checkboxDone": lipgloss.NewStyle().
			Foreground(colors.Success),

		// Task text styles
		"taskPending": lipgloss.NewStyle().
			Foreground(colors.Text),

		"taskDone": lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(colors.TextMuted),

		// Task expanded view
		"taskExpanded": lipgloss.NewStyle().
			Foreground(colors.Text).
			Background(colors.Subtle).
			Padding(0, 2),

		// Due date
		"dueDate": lipgloss.NewStyle().
			Foreground(colors.TextMuted).
			Italic(true),

		// Pagination info
		"pageInfo": lipgloss.NewStyle().
			Foreground(colors.TextMuted).
			Align(lipgloss.Center).
			Italic(true),

		// Priority indicators
		"priorityHigh": lipgloss.NewStyle().
			Foreground(colors.Error).
			Padding(0, 1).
			Margin(0, 1, 0, 0).
			Bold(true),

		"priorityCritical": lipgloss.NewStyle().
			Foreground(colors.Critical).
			Background(colors.Error).
			Padding(0, 1).
			Margin(0, 1, 0, 0).
			Bold(true),

		"priorityMedium": lipgloss.NewStyle().
			Foreground(colors.Warning).
			Padding(0, 1).
			Margin(0, 1, 0, 0),

		"priorityLow": lipgloss.NewStyle().
			Foreground(colors.Success).
			Padding(0, 1).
			Margin(0, 1, 0, 0),

		// Default category badge
		"category": lipgloss.NewStyle().
			Foreground(colors.Secondary).
			Italic(true).
			Padding(0, 1).
			MarginLeft(1),

		// Status bar at the bottom
		"statusBar": lipgloss.NewStyle().
			Foreground(colors.TextMuted).
			Padding(0, 1).
			Border(lipgloss.Border{Top: "─"}, false, false, true, false).
			BorderForeground(colors.Border),

		// Task separator
		"taskSeparator": lipgloss.NewStyle().
			Foreground(colors.Subtle),

		// Tab inactive
		"tabInactive": lipgloss.NewStyle().
			Foreground(colors.TextMuted).
			Padding(0, 2).
			MarginRight(1),

		// Tab active
		"tabActive": lipgloss.NewStyle().
			Foreground(colors.Text).
			Background(colors.Primary).
			Bold(true).
			Padding(0, 2).
			MarginRight(1),

		// Date styles
		"date": lipgloss.NewStyle().
			Foreground(colors.TextMuted).
			Width(10),
	}

	return styles
}
