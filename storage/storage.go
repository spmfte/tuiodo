package storage

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spmfte/tuiodo/model"
)

const TodoFilePath = "TODO.md"

// Regex patterns for parsing metadata
var (
	priorityPattern = regexp.MustCompile(`@priority:(high|medium|low)`)
	dueDatePattern  = regexp.MustCompile(`@due:(\d{4}-\d{2}-\d{2})`)
)

// LoadTasks loads tasks from the TODO.md file
func LoadTasks() []model.Task {
	content, err := os.ReadFile(TodoFilePath)
	if err != nil {
		// File doesn't exist or can't be read, return empty task list
		return []model.Task{}
	}

	lines := strings.Split(string(content), "\n")
	tasks := []model.Task{}

	var currentCategory string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check if this is a category header (starts with ##)
		if strings.HasPrefix(line, "## ") {
			currentCategory = strings.TrimPrefix(line, "## ")
			continue
		}

		// Parse task items (- [ ] or - [x])
		if strings.HasPrefix(line, "- [") && len(line) > 5 {
			isDone := line[3] == 'x' || line[3] == 'X'
			description := strings.TrimSpace(line[5:])

			// Extract priority if present
			var priority model.Priority
			priorityMatch := priorityPattern.FindStringSubmatch(description)
			if len(priorityMatch) > 1 {
				priorityStr := priorityMatch[1]
				switch priorityStr {
				case "high":
					priority = model.PriorityHigh
				case "medium":
					priority = model.PriorityMedium
				case "low":
					priority = model.PriorityLow
				default:
					priority = model.PriorityNone
				}
				// Remove the priority tag from description
				description = strings.TrimSpace(priorityPattern.ReplaceAllString(description, ""))
			}

			// Extract due date if present
			dueDate := ""
			dueDateMatch := dueDatePattern.FindStringSubmatch(description)
			if len(dueDateMatch) > 1 {
				dueDate = dueDateMatch[1]
				// Remove the due date tag from description
				description = strings.TrimSpace(dueDatePattern.ReplaceAllString(description, ""))
			}

			tasks = append(tasks, model.Task{
				Description: description,
				Done:        isDone,
				Category:    currentCategory,
				Priority:    priority,
				DueDate:     dueDate,
			})
		}
	}

	return tasks
}

// SaveTasks saves tasks to the TODO.md file
func SaveTasks(tasks []model.Task) error {
	var content strings.Builder

	// Group tasks by category
	categorizedTasks := make(map[string][]model.Task)
	for _, task := range tasks {
		category := task.Category
		if category == "" {
			category = "Uncategorized"
		}
		categorizedTasks[category] = append(categorizedTasks[category], task)
	}

	// Write tasks by category
	for category, categoryTasks := range categorizedTasks {
		content.WriteString(fmt.Sprintf("## %s\n\n", category))

		for _, task := range categoryTasks {
			checkmark := " "
			if task.Done {
				checkmark = "x"
			}

			description := task.Description

			// Add priority tag if set
			if task.Priority != "" {
				description = fmt.Sprintf("%s @priority:%s", description, task.Priority)
			}

			// Add due date tag if set
			if task.DueDate != "" {
				description = fmt.Sprintf("%s @due:%s", description, task.DueDate)
			}

			content.WriteString(fmt.Sprintf("- [%s] %s\n", checkmark, description))
		}
		content.WriteString("\n")
	}

	return os.WriteFile(TodoFilePath, []byte(content.String()), 0644)
}
