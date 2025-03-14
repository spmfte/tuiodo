package config

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles represents all the styles used in the application
type Styles struct {
	// General
	Title     lipgloss.Style
	Subtitle  lipgloss.Style
	StatusBar lipgloss.Style
	Footer    lipgloss.Style
	Accent    lipgloss.Style
	Warning   lipgloss.Style
	Error     lipgloss.Style

	// Container styles
	Container       lipgloss.Style
	TabContainer    lipgloss.Style
	TaskContainer   lipgloss.Style
	HeaderContainer lipgloss.Style

	// Tab styles
	Tab       lipgloss.Style
	ActiveTab lipgloss.Style

	// Task styles
	Task         lipgloss.Style
	SelectedTask lipgloss.Style
	DoneTask     lipgloss.Style
	TaskText     lipgloss.Style
	Checkbox     lipgloss.Style
	CheckboxDone lipgloss.Style
	Separator    lipgloss.Style

	// Priority indicators
	PriorityHigh   lipgloss.Style
	PriorityMedium lipgloss.Style
	PriorityLow    lipgloss.Style

	// Input mode
	InputField  lipgloss.Style
	InputPrompt lipgloss.Style
}

// GetStyles returns all the styles based on the configuration
func GetStyles(cfg Config) Styles {
	var s Styles

	// Create border style
	var border lipgloss.Border
	switch cfg.UI.BorderStyle {
	case "normal":
		border = lipgloss.NormalBorder()
	case "double":
		border = lipgloss.DoubleBorder()
	case "thick":
		border = lipgloss.ThickBorder()
	case "none":
		border = lipgloss.HiddenBorder()
	default:
		border = lipgloss.RoundedBorder()
	}

	// Parse colors
	primaryColor, _ := ParseColor(cfg.Colors.Primary)
	secondaryColor, _ := ParseColor(cfg.Colors.Secondary)
	tertiaryColor, _ := ParseColor(cfg.Colors.Tertiary)
	successColor, _ := ParseColor(cfg.Colors.Success)
	warningColor, _ := ParseColor(cfg.Colors.Warning)
	errorColor, _ := ParseColor(cfg.Colors.Error)
	textColor, _ := ParseColor(cfg.Colors.Text)
	textDimColor, _ := ParseColor(cfg.Colors.TextDim)
	textMutedColor, _ := ParseColor(cfg.Colors.TextMuted)
	borderColor, _ := ParseColor(cfg.Colors.Border)
	borderFocusColor, _ := ParseColor(cfg.Colors.BorderFocus)
	subtleColor, _ := ParseColor(cfg.Colors.Subtle)
	priorityHighColor, _ := ParseColor(cfg.Colors.PriorityHigh)
	priorityMediumColor, _ := ParseColor(cfg.Colors.PriorityMedium)
	priorityLowColor, _ := ParseColor(cfg.Colors.PriorityLow)
	taskDoneColor, _ := ParseColor(cfg.Colors.TaskDone)

	// Initialize styles
	s.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor)

	s.Subtitle = lipgloss.NewStyle().
		Foreground(secondaryColor)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(textColor)

	s.Footer = lipgloss.NewStyle().
		Foreground(textMutedColor)

	// Add styles for previously unused colors
	s.Accent = lipgloss.NewStyle().
		Foreground(tertiaryColor)

	s.Warning = lipgloss.NewStyle().
		Foreground(warningColor)

	s.Error = lipgloss.NewStyle().
		Foreground(errorColor)

	// Container styles
	s.Container = lipgloss.NewStyle()
	if cfg.UI.EnableBorders {
		s.Container = s.Container.
			Border(border).
			BorderForeground(borderColor)
	}

	s.HeaderContainer = lipgloss.NewStyle().
		Foreground(textDimColor)
	if cfg.UI.EnableBorders {
		s.HeaderContainer = s.HeaderContainer.
			Border(border).
			BorderBottom(false).
			BorderForeground(borderColor)
	}

	s.TabContainer = lipgloss.NewStyle()
	if cfg.UI.EnableBorders {
		s.TabContainer = s.TabContainer.
			Border(border).
			BorderTop(false).
			BorderBottom(false).
			BorderForeground(borderColor)
	}

	s.TaskContainer = lipgloss.NewStyle().
		Foreground(textColor)
	if cfg.UI.EnableBorders {
		s.TaskContainer = s.TaskContainer.
			Border(border).
			BorderTop(false).
			BorderForeground(borderColor)
	}

	// Tab styles
	s.Tab = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(textColor)

	s.ActiveTab = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(primaryColor).
		Bold(true)

	// Task styles
	s.Task = lipgloss.NewStyle().
		Foreground(textColor)

	s.SelectedTask = lipgloss.NewStyle().
		Foreground(textColor).
		Bold(true)

	s.DoneTask = lipgloss.NewStyle().
		Foreground(taskDoneColor).
		Strikethrough(true)

	s.TaskText = lipgloss.NewStyle().
		Foreground(textColor)

	s.Checkbox = lipgloss.NewStyle().
		Foreground(primaryColor)

	s.CheckboxDone = lipgloss.NewStyle().
		Foreground(successColor)

	s.Separator = lipgloss.NewStyle().
		Foreground(subtleColor)

	// Priority indicators
	s.PriorityHigh = lipgloss.NewStyle().
		Foreground(priorityHighColor).
		Bold(true)

	s.PriorityMedium = lipgloss.NewStyle().
		Foreground(priorityMediumColor)

	s.PriorityLow = lipgloss.NewStyle().
		Foreground(priorityLowColor)

	// Input mode
	s.InputField = lipgloss.NewStyle().
		Foreground(textColor)
	if cfg.UI.EnableBorders {
		s.InputField = s.InputField.
			Border(border).
			BorderForeground(borderFocusColor)
	}

	s.InputPrompt = lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true)

	return s
}

// GetKeyBindings returns the mapped keybindings from config
func GetKeyBindings(cfg Config) KeyBindings {
	// This would convert the config's keybindings to whatever format
	// the application uses for key handling
	return KeyBindings{
		Quit:           cfg.Keybindings.QuitKey,
		AddTask:        cfg.Keybindings.AddTaskKey,
		EditTask:       cfg.Keybindings.EditTaskKey,
		DeleteTask:     cfg.Keybindings.DeleteTaskKey,
		ToggleTask:     cfg.Keybindings.ToggleTaskKey,
		CyclePriority:  cfg.Keybindings.CyclePriorityKey,
		CycleCategory:  cfg.Keybindings.CycleCategoryKey,
		CycleTab:       cfg.Keybindings.CycleTabKey,
		NextPage:       cfg.Keybindings.NextPageKey,
		PrevPage:       cfg.Keybindings.PrevPageKey,
		MoveCursorUp:   cfg.Keybindings.MoveCursorUpKey,
		MoveCursorDown: cfg.Keybindings.MoveCursorDownKey,
		Help:           cfg.Keybindings.HelpKey,
	}
}

// KeyBindings represents the key mappings used by the application
type KeyBindings struct {
	Quit           []string
	AddTask        []string
	EditTask       []string
	DeleteTask     []string
	ToggleTask     []string
	CyclePriority  []string
	CycleCategory  []string
	CycleTab       []string
	NextPage       []string
	PrevPage       []string
	MoveCursorUp   []string
	MoveCursorDown []string
	Help           []string
}
