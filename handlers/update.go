package handlers

import (
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
		case "?", "h", "escape":
			m.HelpVisible = false
		case "q", "ctrl+c":
			return m, tea.Quit
		}
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
			m.DeleteCurrentTask()
			storage.SaveTasks(m.Tasks)
			m.SetStatus("Task deleted")

			// Recalculate pagination after deleting a task
			m.RecalculatePagination()
		}
	case "enter", "space": // Toggle task completion
		filteredTasks := m.GetVisibleTasks()
		if len(filteredTasks) > 0 && m.Cursor < len(filteredTasks) {
			m.ToggleCurrentTask()
			storage.SaveTasks(m.Tasks)

			// Show status message
			if filteredTasks[m.Cursor].Done {
				m.SetStatus("Task marked as complete")
			} else {
				m.SetStatus("Task marked as incomplete")
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
