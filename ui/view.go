package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spmfte/tuiodo/model"
)

// View renders the application
func View(m model.Model) string {
	// Get styles
	colors := AppColors()
	styles := CreateStyles(colors)

	// Use terminal dimensions for proper sizing
	maxWidth := m.Width
	if maxWidth <= 0 {
		maxWidth = 80 // Fallback if window size isn't available
	}

	// Set container width (no margins to avoid border issues)
	containerWidth := maxWidth - 2 // Just minimal padding
	if containerWidth < 20 {
		containerWidth = 20
	}

	// Update widths for components that need explicit widths
	styles["titleBar"] = styles["titleBar"].Width(containerWidth)
	styles["listContainer"] = styles["listContainer"].Width(containerWidth)
	styles["inputBox"] = styles["inputBox"].Width(containerWidth)
	styles["statusBar"] = styles["statusBar"].Width(containerWidth)

	// If showing help screen, render that instead
	if m.HelpVisible {
		return renderHelpScreen(styles, containerWidth, m.Height)
	}

	// Build the UI components
	var appContent []string

	// === HEADER SECTION ===
	appContent = append(appContent, renderHeader(m, styles, containerWidth))

	// === TABS SECTION ===
	appContent = append(appContent, renderTabs(m, styles, containerWidth))

	// === INPUT FORM (when in input mode) ===
	if m.InputMode || m.EditingTask {
		appContent = append(appContent, renderInputForm(m, styles, containerWidth))
	} else {
		// === TASKS SECTION (when not in input mode) ===
		appContent = append(appContent, renderTaskList(m, styles, containerWidth))
	}

	// === STATUS BAR ===
	appContent = append(appContent, renderStatusBar(m, styles, containerWidth))

	// Render the full UI (no vertical centering to avoid rendering issues)
	return strings.Join(appContent, "\n")
}

// renderHeader creates the application header
func renderHeader(m model.Model, styles map[string]lipgloss.Style, width int) string {
	var titleBar strings.Builder

	// App title with version badge
	title := styles["title"].Render("TUIODO")
	version := styles["versionBadge"].Render("v1.0")

	// Filter indicator if a category filter is active
	filterLabel := ""
	if m.CurrentFilter != "" {
		filterLabel = styles["filterIndicator"].Render("Category: " + m.CurrentFilter)
	}

	// Assemble the title bar with correct spacing
	emptySpace := width - lipgloss.Width(title) - lipgloss.Width(version) - lipgloss.Width(filterLabel) - 4
	if emptySpace < 0 {
		emptySpace = 0
	}

	titleBar.WriteString(title)
	titleBar.WriteString(" ")
	titleBar.WriteString(version)
	titleBar.WriteString(strings.Repeat(" ", emptySpace))
	titleBar.WriteString(filterLabel)

	return styles["titleBar"].Render(titleBar.String())
}

// renderTabs creates the tab navigation bar
func renderTabs(m model.Model, styles map[string]lipgloss.Style, width int) string {
	tabs := []struct {
		title string
		view  model.TabView
	}{
		{"All", model.TabAll},
		{"Pending", model.TabPending},
		{"Completed", model.TabCompleted},
	}

	var renderedTabs []string

	// Render each tab with appropriate active/inactive styling
	for _, tab := range tabs {
		var tabStyle lipgloss.Style
		if tab.view == m.CurrentView {
			tabStyle = styles["tabActive"]
		} else {
			tabStyle = styles["tabInactive"]
		}
		renderedTabs = append(renderedTabs, tabStyle.Render(tab.title))
	}

	return lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(1).
		Render(strings.Join(renderedTabs, " "))
}

// renderInputForm creates the input form for adding/editing tasks
func renderInputForm(m model.Model, styles map[string]lipgloss.Style, width int) string {
	var title string
	if m.EditingTask {
		title = "Edit Task"
	} else {
		title = "New Task"
	}

	prompt := styles["inputPrompt"].Render(title)
	hint := styles["inputHint"].Render(" (Format: Category: Task description)")
	cursor := styles["inputCursor"].Render("▋")

	inputForm := []string{
		prompt + hint,
		"",
		styles["input"].Render("→ " + m.Input + cursor),
		"",
		styles["inputHint"].Render("Press Enter to save, Esc to cancel"),
	}

	return styles["inputBox"].Render(strings.Join(inputForm, "\n"))
}

// renderTaskList creates the list of tasks
func renderTaskList(m model.Model, styles map[string]lipgloss.Style, width int) string {
	// Get filtered and paginated tasks
	visibleTasks := m.GetVisibleTasks()

	// If no tasks, render empty message
	if len(visibleTasks) == 0 {
		emptyText := "No tasks yet. Press 'a' to add one."
		if m.CurrentFilter != "" {
			emptyText = "No tasks in category '" + m.CurrentFilter + "'"
		}
		return styles["emptyMessage"].Render(emptyText)
	}

	var taskList []string

	// Create header row
	taskList = append(taskList, styles["taskHeader"].Render(fmt.Sprintf("%-3s %-30s %s", "", "TASK", "CATEGORY")))

	// Calculate available width for task description
	maxDescWidth := width - 35 // Adjust based on other elements
	if maxDescWidth < 15 {
		maxDescWidth = 15
	}

	// Render each task
	for i, task := range visibleTasks {
		var taskRow strings.Builder

		// Cursor indicator with checkbox
		var checkboxStyle lipgloss.Style
		checkbox := "[ ]"
		if task.Done {
			checkbox = "[✓]"
			checkboxStyle = styles["checkboxDone"]
		} else {
			checkboxStyle = styles["checkboxPending"]
		}

		// Cursor indicator
		cursor := "  "
		if i == m.Cursor {
			cursor = styles["cursor"].Render("→ ")
		}

		taskRow.WriteString(cursor)
		taskRow.WriteString(checkboxStyle.Render(checkbox))
		taskRow.WriteString(" ")

		// Priority indicator if set
		if task.Priority != "" {
			var priorityStyle lipgloss.Style
			switch task.Priority {
			case model.PriorityHigh:
				priorityStyle = styles["priorityHigh"]
			case model.PriorityMedium:
				priorityStyle = styles["priorityMedium"]
			case model.PriorityLow:
				priorityStyle = styles["priorityLow"]
			}
			taskRow.WriteString(priorityStyle.Render(string(task.Priority)))
			taskRow.WriteString(" ")
		}

		// Task description
		taskStyle := styles["taskPending"]
		if task.Done {
			taskStyle = styles["taskDone"]
		}

		// Make sure description fits without overflowing
		description := task.Description
		if len(description) > maxDescWidth {
			description = description[:maxDescWidth-3] + "..."
		}

		taskRow.WriteString(taskStyle.Render(description))

		// Category badge if present
		if task.Category != "" {
			// Add some spacing
			padding := 2
			taskRow.WriteString(strings.Repeat(" ", padding))
			taskRow.WriteString(styles["category"].Render(task.Category))
		}

		// Add due date if present
		if task.DueDate != "" {
			taskRow.WriteString(" ")
			taskRow.WriteString(styles["dueDate"].Render(task.DueDate))
		}

		taskList = append(taskList, taskRow.String())

		// Only add separators if we have more than 1 task
		if len(visibleTasks) > 1 && i < len(visibleTasks)-1 {
			taskList = append(taskList, styles["taskSeparator"].Render(strings.Repeat("─", width-4)))
		}
	}

	// Add pagination info if needed
	if m.Pagination.TotalPages > 1 {
		pageInfo := fmt.Sprintf("Page %d of %d", m.Pagination.Page+1, m.Pagination.TotalPages)
		taskList = append(taskList, "")
		taskList = append(taskList, styles["pageInfo"].Render(pageInfo))
	}

	return styles["listContainer"].Render(strings.Join(taskList, "\n"))
}

// renderStatusBar creates the status bar at the bottom
func renderStatusBar(m model.Model, styles map[string]lipgloss.Style, width int) string {
	// Default help text
	statusText := "Press ? for help"

	// Show temporary status message if set
	if m.StatusMessage != "" {
		statusText = m.StatusMessage
	}

	return styles["statusBar"].Render(statusText)
}

// renderHelpScreen creates a comprehensive help screen
func renderHelpScreen(styles map[string]lipgloss.Style, width int, height int) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles["title"].GetForeground()).
		Padding(0, 1).
		MarginBottom(1)

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles["secondary"].GetForeground()).
		MarginTop(1).
		MarginBottom(1)

	contentStyle := lipgloss.NewStyle().
		Foreground(styles["taskPending"].GetForeground())

	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles["helpCommand"].GetForeground())

	helpContent := []string{
		titleStyle.Render("TUIODO KEYBOARD SHORTCUTS"),
		"",
		sectionStyle.Render("NAVIGATION"),
		fmt.Sprintf("%s : Move cursor up", keyStyle.Render("j/k, ↑/↓")),
		fmt.Sprintf("%s : Navigate between tabs", keyStyle.Render("tab, t")),
		fmt.Sprintf("%s : Next/previous page", keyStyle.Render("n/b, →/←")),
		"",
		sectionStyle.Render("TASK MANAGEMENT"),
		fmt.Sprintf("%s : Add new task", keyStyle.Render("a")),
		fmt.Sprintf("%s : Edit current task", keyStyle.Render("e")),
		fmt.Sprintf("%s : Delete current task", keyStyle.Render("d")),
		fmt.Sprintf("%s : Toggle task completion", keyStyle.Render("space, enter")),
		fmt.Sprintf("%s : Cycle priority (none/low/medium/high)", keyStyle.Render("p")),
		"",
		sectionStyle.Render("FILTERING"),
		fmt.Sprintf("%s : Cycle through categories", keyStyle.Render("c")),
		fmt.Sprintf("%s : Switch between views (All/Pending/Completed)", keyStyle.Render("tab, t")),
		"",
		sectionStyle.Render("OTHER"),
		fmt.Sprintf("%s : Show/hide this help", keyStyle.Render("?, h, F1")),
		fmt.Sprintf("%s : Quit application", keyStyle.Render("q, Ctrl+C")),
		"",
		contentStyle.Render("Press any key to close this help screen"),
	}

	// Don't try to be clever with vertical centering - just render the help content
	return styles["helpBox"].Copy().Width(width).Render(strings.Join(helpContent, "\n"))
}
