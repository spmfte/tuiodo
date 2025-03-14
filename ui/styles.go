package ui

import (
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

	// Category colors
	CategoryWork     lipgloss.Color
	CategoryPersonal lipgloss.Color
	CategoryHealth   lipgloss.Color
	CategoryFinance  lipgloss.Color
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

		// Category colors
		CategoryWork:     lipgloss.Color("#3B82F6"), // Blue
		CategoryPersonal: lipgloss.Color("#EC4899"), // Pink
		CategoryHealth:   lipgloss.Color("#10B981"), // Green
		CategoryFinance:  lipgloss.Color("#6366F1"), // Indigo
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
	// This is a simplified implementation
	// In a real implementation, this would parse the provided styles config
	// and update currentStyles and currentColors accordingly

	// For now, we'll just ensure the default styles are loaded
	// if they haven't been already
	if currentStyles == nil {
		currentColors = AppColors()
		currentStyles = CreateStyles(currentColors)
	}

	// When integrating with the config system fully, this would:
	// 1. Extract colors from the config
	// 2. Create a new Colors struct
	// 3. Generate new styles based on these colors
	// 4. Update currentColors and currentStyles
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

	return map[string]lipgloss.Style{
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

		// Category badge
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

		// Category-specific styles
		"category_work": lipgloss.NewStyle().
			Foreground(colors.CategoryWork).
			Italic(true).
			Padding(0, 1),

		"category_personal": lipgloss.NewStyle().
			Foreground(colors.CategoryPersonal).
			Italic(true).
			Padding(0, 1),

		"category_health": lipgloss.NewStyle().
			Foreground(colors.CategoryHealth).
			Italic(true).
			Padding(0, 1),

		"category_finance": lipgloss.NewStyle().
			Foreground(colors.CategoryFinance).
			Italic(true).
			Padding(0, 1),
	}
}
