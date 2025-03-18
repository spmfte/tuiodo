package ui

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spmfte/tuiodo/model"
	"github.com/spmfte/tuiodo/storage"
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

	// Get app info
	version, commit, _ := GetAppInfo()

	// Format version string
	versionStr := "v" + version
	if commit != "" {
		versionStr += " (" + commit + ")"
	}

	// App title with version badge
	title := styles["title"].Render("TUIODO")
	versionBadge := styles["versionBadge"].Render(versionStr)

	// Filter indicator if a category filter is active
	filterLabel := ""
	if m.CurrentFilter != "" {
		filterLabel = styles["filterIndicator"].Render("Category: " + m.CurrentFilter)
	}

	// Assemble the title bar with correct spacing
	emptySpace := width - lipgloss.Width(title) - lipgloss.Width(versionBadge) - lipgloss.Width(filterLabel) - 4
	if emptySpace < 0 {
		emptySpace = 0
	}

	titleBar.WriteString(title)
	titleBar.WriteString(" ")
	titleBar.WriteString(versionBadge)
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

	// Create header row with improved dynamic layout
	// Calculate widths based on screen size with fixed proportions
	spacing := 2 // Fixed spacing between columns
	minTitleWidth := 12
	minCategoryWidth := 10
	minDateWidth := 10

	// Reserve space for cursor, checkbox, priority indicators, etc.
	reservedWidth := 12

	// Available width for main columns
	contentWidth := width - reservedWidth - 6 // Account for borders and padding

	// Set column widths as proportions of available space
	// with minimums to maintain readability
	dateWidth := minDateWidth
	categoryWidth := max(minCategoryWidth, contentWidth*20/100)                       // 20% for category
	taskWidth := max(minTitleWidth, contentWidth-categoryWidth-dateWidth-(spacing*2)) // Remaining space for task

	// Create the header with proper spacing - we don't need headerFormat anymore
	taskHeader := styles["taskHeader"].Copy().
		MarginLeft(3).
		Bold(true).
		Render(fmt.Sprintf("TASK%s%s%sCREATED",
			strings.Repeat(" ", taskWidth-4),
			"CATEGORY",
			strings.Repeat(" ", categoryWidth+spacing-8)))

	taskList = append(taskList, taskHeader)

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

		// Priority indicator if set and task is not completed
		prioritySpace := 0
		if task.Priority != "" && !task.Done {
			var priorityStyle lipgloss.Style
			switch task.Priority {
			case model.PriorityCritical:
				priorityStyle = styles["priorityCritical"]
			case model.PriorityHigh:
				priorityStyle = styles["priorityHigh"]
			case model.PriorityMedium:
				priorityStyle = styles["priorityMedium"]
			case model.PriorityLow:
				priorityStyle = styles["priorityLow"]
			}
			// Adjust space for the priority tag
			priorityText := string(task.Priority)
			priorityColored := priorityStyle.Render(priorityText)
			prioritySpace = lipgloss.Width(priorityColored) + 1
			taskRow.WriteString(priorityColored)
			taskRow.WriteString(" ")
		}

		// Calculate adjusted task width considering priority
		adjustedTaskWidth := taskWidth - prioritySpace

		// Task description
		taskStyle := styles["taskPending"]
		if task.Done {
			taskStyle = styles["taskDone"]
		}

		// Clean description and truncate if needed
		description := cleanMetadata(task.Description)

		// Calculate category length for potential truncation
		categoryStrLen := len(task.Category)
		extraCategoryLen := max(0, categoryStrLen-int(categoryWidth))

		// If category is longer than its allocated space, take space from task
		adjustedTaskWidth = max(minTitleWidth, adjustedTaskWidth-extraCategoryLen)

		if len(description) > adjustedTaskWidth {
			description = description[:adjustedTaskWidth-3] + "..."
		}

		// Pad the description to its adjusted width
		paddedDesc := fmt.Sprintf("%-*s", adjustedTaskWidth, description)
		taskRow.WriteString(taskStyle.Render(paddedDesc))

		// Category with appropriate width
		category := task.Category
		if len(category) > int(categoryWidth) {
			category = category[:categoryWidth-3] + "..."
		}

		categoryStr := fmt.Sprintf("%-*s", categoryWidth, category)
		if task.Category != "" {
			categoryStyle := getCategoryStyle(styles, task.Category)

			// For completed tasks, use a dimmed version of the category style
			if task.Done {
				categoryStyle = styles["taskDone"].Copy().
					Strikethrough(false).
					Italic(true).
					Padding(0, 1).
					MarginLeft(1)
			}

			taskRow.WriteString(categoryStyle.Render(categoryStr))
		} else {
			taskRow.WriteString(strings.Repeat(" ", int(categoryWidth)))
		}

		// Add spacing between category and date
		taskRow.WriteString(strings.Repeat(" ", spacing))

		// Creation date (local time, date only)
		createdDate := task.CreatedAt.Local().Format("2006-01-02")
		taskRow.WriteString(styles["date"].Render(createdDate))

		taskList = append(taskList, taskRow.String())

		// If this task is expanded, show its full details
		if m.TaskExpanded && i == m.ExpandedTaskIdx {
			// Create a cleaner expanded view without blocks or borders
			expandedDetails := []string{
				"",
				styles["secondary"].Copy().Bold(true).Render("Full Description:"),
				styles["taskPending"].Copy().Render("  " + cleanMetadata(task.Description)),
				"",
			}

			// Add metadata in a cleaner format
			metaInfo := extractMetadata(task.Description)

			// Display metadata in a cleaner two-column format
			infoLayout := [][]string{
				{styles["taskHeader"].Copy().Render("Created:"), styles["inputHint"].Render(task.CreatedAt.Local().Format("2006-01-02 15:04:05"))},
			}

			// Add priority if present and task is not completed
			if task.Priority != "" && !task.Done {
				var priorityText string
				switch task.Priority {
				case model.PriorityCritical:
					priorityText = styles["priorityCritical"].Render(string(task.Priority))
				case model.PriorityHigh:
					priorityText = styles["priorityHigh"].Render(string(task.Priority))
				case model.PriorityMedium:
					priorityText = styles["priorityMedium"].Render(string(task.Priority))
				case model.PriorityLow:
					priorityText = styles["priorityLow"].Render(string(task.Priority))
				}
				infoLayout = append(infoLayout, []string{styles["taskHeader"].Copy().Render("Priority:"), priorityText})
			}

			// Add category if present
			if task.Category != "" {
				categoryStyle := getCategoryStyle(styles, task.Category)

				// For completed tasks, use a dimmed version of the category style
				if task.Done {
					categoryStyle = styles["taskDone"].Copy().
						Strikethrough(false).
						Italic(true).
						Padding(0, 1).
						MarginLeft(1)
				}

				infoLayout = append(infoLayout, []string{styles["taskHeader"].Copy().Render("Category:"), categoryStyle.Render(task.Category)})
			}

			// Format the two-column layout with consistent spacing
			for _, row := range infoLayout {
				expandedDetails = append(expandedDetails, fmt.Sprintf("  %-12s %s", row[0], row[1]))
			}

			// Add metadata tags if present
			if len(metaInfo) > 0 {
				expandedDetails = append(expandedDetails, "")
				expandedDetails = append(expandedDetails, styles["secondary"].Copy().Bold(true).Render("Metadata Tags:"))

				for key, value := range metaInfo {
					expandedDetails = append(expandedDetails, fmt.Sprintf("  @%-10s %s", key+":", value))
				}
			}

			// Add a help hint at the bottom
			expandedDetails = append(expandedDetails, "")
			expandedDetails = append(expandedDetails, styles["inputHint"].Italic(true).Render("  Press 'x' to collapse"))

			// Render without additional styling that might cause formatting issues
			expandedView := strings.Join(expandedDetails, "\n")
			taskList = append(taskList, expandedView)
		}

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

// cleanMetadata removes metadata tags from task description
func cleanMetadata(description string) string {
	// Remove @created tag with timestamp
	createdPattern := regexp.MustCompile(`\s*@created:\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)
	description = createdPattern.ReplaceAllString(description, "")

	// Remove @priority tag
	priorityPattern := regexp.MustCompile(`\s*@priority:(high|medium|low|critical)`)
	description = priorityPattern.ReplaceAllString(description, "")

	// Remove any trailing/leading whitespace
	return strings.TrimSpace(description)
}

// extractMetadata extracts metadata tags from a task description
func extractMetadata(description string) map[string]string {
	metadata := make(map[string]string)

	// Find all @tag:value patterns
	metadataPattern := regexp.MustCompile(`@([a-zA-Z0-9_]+):([^\s@]+)`)
	matches := metadataPattern.FindAllStringSubmatch(description, -1)

	for _, match := range matches {
		if len(match) == 3 {
			metadata[match[1]] = match[2]
		}
	}

	return metadata
}

// renderStatusBar creates the status bar at the bottom
func renderStatusBar(m model.Model, styles map[string]lipgloss.Style, width int) string {
	var statusBar strings.Builder

	// Show temporary status message if set
	var statusText string
	if m.StatusMessage != "" {
		// Use different styling for different types of messages
		if strings.Contains(strings.ToLower(m.StatusMessage), "error") ||
			strings.Contains(strings.ToLower(m.StatusMessage), "failed") {
			statusText = styles["error"].Render(m.StatusMessage)
		} else if strings.Contains(strings.ToLower(m.StatusMessage), "complete") ||
			strings.Contains(strings.ToLower(m.StatusMessage), "added") {
			statusText = styles["success"].Render(m.StatusMessage)
		} else {
			statusText = styles["secondary"].Bold(true).Render(m.StatusMessage)
		}
	} else {
		// Default help text when no status message
		statusText = styles["inputHint"].Render("Press ? for help")
	}

	// Left side: status message or help text
	leftSide := statusText

	// Right side: storage info and task stats with progress bar
	totalTasks := len(m.Tasks)
	doneTasks := 0
	for _, t := range m.Tasks {
		if t.Done {
			doneTasks++
		}
	}

	// Format stats with progress percentage and mini progress bar
	var progressPct float64
	if totalTasks > 0 {
		progressPct = float64(doneTasks) / float64(totalTasks) * 100
	}

	// Create a color based on progress
	progressColor := styles["success"].GetForeground()
	if progressPct < 30 {
		progressColor = styles["warning"].GetForeground()
	} else if progressPct < 70 {
		progressColor = styles["primary"].GetForeground()
	}

	// Create a mini progress bar
	progressBarWidth := 10
	filledCells := int(progressPct / 100 * float64(progressBarWidth))
	progressBar := strings.Repeat("█", filledCells) +
		strings.Repeat("░", progressBarWidth-filledCells)

	progressStyle := lipgloss.NewStyle().Foreground(progressColor)
	percentText := fmt.Sprintf("%.0f%%", progressPct)

	// Add file info
	filePath := storage.GetStoragePath()
	fileName := filepath.Base(filePath)

	// Format right side
	rightSide := fmt.Sprintf(
		"%s %s %s • %s",
		fmt.Sprintf("%d/%d", doneTasks, totalTasks),
		progressStyle.Render(progressBar),
		progressStyle.Render(percentText),
		styles["inputHint"].Render(fileName),
	)

	// Determine spacing
	spacerWidth := width - lipgloss.Width(leftSide) - lipgloss.Width(rightSide) - 2
	if spacerWidth < 1 {
		spacerWidth = 1
	}

	// Compose the status bar
	statusBar.WriteString(leftSide)
	statusBar.WriteString(strings.Repeat(" ", spacerWidth))
	statusBar.WriteString(rightSide)

	// Use the existing statusBar style
	return styles["statusBar"].Render(statusBar.String())
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
		fmt.Sprintf("%s : Delete task (press twice to confirm)", keyStyle.Render("d")),
		fmt.Sprintf("%s : Undo last deletion", keyStyle.Render("u")),
		fmt.Sprintf("%s : Toggle task completion", keyStyle.Render("space, enter")),
		fmt.Sprintf("%s : Expand/collapse task details", keyStyle.Render("x")),
		fmt.Sprintf("%s : Cycle priority (none/low/medium/high/critical)", keyStyle.Render("p")),
		"",
		sectionStyle.Render("SORTING & FILTERING"),
		fmt.Sprintf("%s : Sort by priority", keyStyle.Render("s")),
		fmt.Sprintf("%s : Sort by creation date", keyStyle.Render("S")),
		fmt.Sprintf("%s : Sort by category", keyStyle.Render("C")),
		fmt.Sprintf("%s : Cycle through categories", keyStyle.Render("c")),
		fmt.Sprintf("%s : Switch between views (All/Pending/Completed)", keyStyle.Render("tab, t")),
		"",
		sectionStyle.Render("OTHER"),
		fmt.Sprintf("%s : Show/hide this help", keyStyle.Render("?, h, F1")),
		fmt.Sprintf("%s : Quit application", keyStyle.Render("q, Ctrl+C")),
		"",
		contentStyle.Render("Press any key to close this help screen"),
	}

	return styles["helpBox"].Copy().Width(width).Render(strings.Join(helpContent, "\n"))
}

// max helper function
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// getCategoryStyle returns the appropriate lipgloss style for a category
func getCategoryStyle(styles map[string]lipgloss.Style, category string) lipgloss.Style {
	// Try to find a specific style for this category
	categoryKey := "category_" + strings.ToLower(category)
	if style, ok := styles[categoryKey]; ok {
		return style
	}

	// Fall back to default category style
	return styles["category"]
}
