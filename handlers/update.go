package handlers

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spmfte/tuiodo/model"
	"github.com/spmfte/tuiodo/storage"
)

// Update processes messages and updates the model accordingly
func Update(msg tea.Msg, m model.Model) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return handleKeypress(msg, m)
	case tea.WindowSizeMsg:
		m.UpdateWindowSize(msg.Width, msg.Height)
		return m, nil
	}
	return m, nil
}

// HandleKeypress processes keyboard input
func handleKeypress(msg tea.KeyMsg, m model.Model) (model.Model, tea.Cmd) {
	// If help is visible, only respond to help toggle or quit
	if m.HelpVisible {
		switch msg.String() {
		case "?", "h", "escape", "esc", "q":
			m.HelpVisible = false
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		default:
			// Allow any key to dismiss help
			m.HelpVisible = false
			return m, nil
		}
	}

	// If delete confirmation is active, any key other than 'd' cancels it
	if m.DeleteConfirm && msg.String() != "d" {
		m.DeleteConfirm = false
		m.SetStatus("Deletion cancelled")
		return m, nil
	}

	// If in input mode, handle input-specific keys
	if m.InputMode {
		return handleInputMode(msg, m)
	}

	// If editing a task, handle edit-specific keys
	if m.EditingTask {
		return handleEditMode(msg, m)
	}

	// Normal mode key handling
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		m.MoveCursorUp()
	case "down", "j":
		m.MoveCursorDown()
	case "tab", "t":
		m.CycleTab()
	case "a": // Add new task
		m.InputMode = true
		m.Input = ""
	case "e": // Edit current task
		filteredTasks := m.GetVisibleTasks()
		if len(filteredTasks) > 0 && m.Cursor < len(filteredTasks) {
			task := filteredTasks[m.Cursor]
			m.EditingTask = true
			m.EditingTaskIdx = m.Cursor

			// Pre-populate input with existing task info
			if task.Category != "" {
				m.Input = task.Category + ": " + task.Description
			} else {
				m.Input = task.Description
			}
		}
	case "d": // Delete task
		filteredTasks := m.GetVisibleTasks()
		if len(filteredTasks) > 0 {
			if m.DeleteConfirm {
				// If already in confirmation mode, execute the delete
				m.DeleteCurrentTask()
				storage.SaveTasks(m.Tasks)
				m.SetStatus("Task deleted (press 'u' to undo)")
				m.DeleteConfirm = false

				// Recalculate pagination after deleting a task
				m.RecalculatePagination()
			} else {
				// First press just enters confirmation mode
				m.DeleteConfirm = true
				m.SetStatus("Press 'd' again to confirm deletion, or any other key to cancel")
			}
		}
	case "enter", " ", "space": // Toggle task completion - added explicit " " and "space" matches
		filteredTasks := m.GetVisibleTasks()
		if len(filteredTasks) > 0 && m.Cursor < len(filteredTasks) {
			// Get the actual task from the filtered list
			taskToToggle := filteredTasks[m.Cursor]

			// Find the task's index in the main task list by comparing relevant fields
			found := false
			for i, task := range m.Tasks {
				// Compare only the essential fields that identify a task
				// Added more robust comparison that doesn't rely as heavily on CreatedAt
				sameDescription := task.Description == taskToToggle.Description
				sameCategory := task.Category == taskToToggle.Category
				sameCreationTime := task.CreatedAt.Equal(taskToToggle.CreatedAt)

				// Use multiple matching criteria for more robust identification
				if sameDescription && sameCategory && sameCreationTime {
					// Toggle completion status
					m.Tasks[i].Done = !m.Tasks[i].Done
					storage.SaveTasks(m.Tasks)

					// Show status message
					if m.Tasks[i].Done {
						m.SetStatus("Task marked as complete")
					} else {
						m.SetStatus("Task marked as incomplete")
					}
					found = true
					break
				}
			}

			// If task wasn't found in main list, log an error status
			if !found {
				m.SetStatus("Error: Could not find task to toggle")
				// Log more details for debugging
				log.Printf("Failed to toggle task: %+v", taskToToggle)
			}
		}
	case "c": // Cycle through categories for filtering
		m.CycleCategory()
	case "p": // Cycle through priorities
		filteredTasks := m.GetVisibleTasks()
		if len(filteredTasks) > 0 && m.Cursor < len(filteredTasks) {
			m.CyclePriority()
			storage.SaveTasks(m.Tasks)

			// Show status based on new priority
			switch m.Tasks[m.Cursor].Priority {
			case model.PriorityHigh:
				m.SetStatus("Task priority set to HIGH")
			case model.PriorityMedium:
				m.SetStatus("Task priority set to MEDIUM")
			case model.PriorityLow:
				m.SetStatus("Task priority set to LOW")
			default:
				m.SetStatus("Task priority cleared")
			}
		}
	case "right", "l", "n": // Next page
		m.NextPage()
	case "left", "h", "b": // Previous page
		m.PrevPage()
	case "?", "F1": // Show help
		m.HelpVisible = true
	case "s": // Sort tasks
		m.SortTasks(model.SortByPriority)
		m.SetStatus("Sorted by priority")
	case "S": // Sort by creation date
		m.SortTasks(model.SortByCreatedAt)
		m.SetStatus("Sorted by creation date")
	case "C": // Sort by category
		m.SortTasks(model.SortByCategory)
		m.SetStatus("Sorted by category")
	case "x": // Expand/collapse task details
		filteredTasks := m.GetVisibleTasks()
		if len(filteredTasks) > 0 && m.Cursor < len(filteredTasks) {
			if m.TaskExpanded && m.ExpandedTaskIdx == m.Cursor {
				// Collapse if already expanded
				m.TaskExpanded = false
			} else {
				// Expand the task
				m.TaskExpanded = true
				m.ExpandedTaskIdx = m.Cursor
			}
		}
	case "u": // Undo last delete
		if m.LastDeleted != nil {
			if m.UndoDelete() {
				storage.SaveTasks(m.Tasks)
				m.SetStatus("Task restored")
			}
		}
	}

	return m, nil
}

// HandleInputMode processes keyboard input in input mode
func handleInputMode(msg tea.KeyMsg, m model.Model) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if strings.TrimSpace(m.Input) != "" {
			// Extract category if included in format "Category: Task description"
			category := ""
			description := m.Input

			if parts := strings.SplitN(m.Input, ":", 2); len(parts) == 2 {
				category = strings.TrimSpace(parts[0])
				description = strings.TrimSpace(parts[1])
			}

			m.AddTask(description, category)
			storage.SaveTasks(m.Tasks)
			m.SetStatus("Task added")

			// Recalculate pagination after adding task
			m.RecalculatePagination()
		}
		m.InputMode = false
		m.Input = ""
	case "esc":
		m.InputMode = false
		m.Input = ""
	case "backspace":
		if len(m.Input) > 0 {
			m.Input = m.Input[:len(m.Input)-1]
		}
	case "tab":
		// Try to auto-complete with a category if typed part of one
		if !strings.Contains(m.Input, ":") {
			partialCategory := strings.TrimSpace(m.Input)
			if partialCategory != "" {
				for category := range m.Categories {
					if strings.HasPrefix(strings.ToLower(category), strings.ToLower(partialCategory)) {
						m.Input = category + ": "
						break
					}
				}
			}
		}
	default:
		// Only add to input if it's a printable character
		if len(msg.String()) == 1 {
			m.Input += msg.String()
		}
	}

	return m, nil
}

// HandleEditMode processes keyboard input in edit mode
func handleEditMode(msg tea.KeyMsg, m model.Model) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if strings.TrimSpace(m.Input) != "" {
			// Extract category if included in format "Category: Task description"
			category := ""
			description := m.Input

			if parts := strings.SplitN(m.Input, ":", 2); len(parts) == 2 {
				category = strings.TrimSpace(parts[0])
				description = strings.TrimSpace(parts[1])
			}

			// Get the task's current priority
			priority := model.PriorityNone
			if m.EditingTaskIdx < len(m.Tasks) {
				priority = m.Tasks[m.EditingTaskIdx].Priority
			}

			m.UpdateTask(m.EditingTaskIdx, description, category, priority)
			storage.SaveTasks(m.Tasks)
			m.SetStatus("Task updated")
		}
		m.EditingTask = false
		m.Input = ""
	case "esc":
		m.EditingTask = false
		m.Input = ""
	case "backspace":
		if len(m.Input) > 0 {
			m.Input = m.Input[:len(m.Input)-1]
		}
	case "tab":
		// Try to auto-complete with a category if typed part of one
		if !strings.Contains(m.Input, ":") {
			partialCategory := strings.TrimSpace(m.Input)
			if partialCategory != "" {
				for category := range m.Categories {
					if strings.HasPrefix(strings.ToLower(category), strings.ToLower(partialCategory)) {
						m.Input = category + ": "
						break
					}
				}
			}
		}
	default:
		// Only add to input if it's a printable character
		if len(msg.String()) == 1 {
			m.Input += msg.String()
		}
	}

	return m, nil
}
