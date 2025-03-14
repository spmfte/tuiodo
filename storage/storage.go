package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/spmfte/tuiodo/model"
)

// DefaultTodoFilePath is the default path for the TODO file
const DefaultTodoFilePath = "TODO.md"

// Storage configuration
var (
	todoFilePath    = DefaultTodoFilePath
	backupDirectory = ""
	maxBackups      = 5
	autoSave        = true
	backupOnSave    = true
	storageWriteMu  sync.Mutex // mutex to prevent concurrent writes
)

// Regex patterns for parsing metadata - compile only once for better performance
var (
	priorityPattern = regexp.MustCompile(`@priority:(high|medium|low)`)
	dueDatePattern  = regexp.MustCompile(`@due:(\d{4}-\d{2}-\d{2})`)
)

// Initialize sets up the storage with configurable settings
func Initialize(filePath string, backupDir string, maxBackupFiles int, enableAutoSave bool, enableBackup bool) {
	if filePath != "" {
		todoFilePath = filePath
	} else {
		todoFilePath = DefaultTodoFilePath
	}

	backupDirectory = backupDir

	if maxBackupFiles > 0 {
		maxBackups = maxBackupFiles
	}

	autoSave = enableAutoSave
	backupOnSave = enableBackup

	// Create the backup directory if it doesn't exist and backups are enabled
	if backupOnSave && backupDirectory != "" {
		os.MkdirAll(backupDirectory, 0755)
	}
}

// GetStoragePath returns the current storage file path
func GetStoragePath() string {
	return todoFilePath
}

// IsAutoSaveEnabled returns whether auto-save is enabled
func IsAutoSaveEnabled() bool {
	return autoSave
}

// LoadTasks loads tasks from the configured TODO file
func LoadTasks() []model.Task {
	content, err := os.ReadFile(todoFilePath)
	if err != nil {
		// File doesn't exist or can't be read, return empty task list
		return make([]model.Task, 0)
	}

	lines := strings.Split(string(content), "\n")
	tasks := make([]model.Task, 0, len(lines)/3) // Preallocate with estimated capacity

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

// SaveTasks saves tasks to the configured TODO file
func SaveTasks(tasks []model.Task) error {
	storageWriteMu.Lock()
	defer storageWriteMu.Unlock()

	// Create backup if configured and the file exists
	if backupOnSave && backupDirectory != "" {
		if _, err := os.Stat(todoFilePath); err == nil {
			createBackup()
		}
	}

	// Build the content with a single StringBuilder for better performance
	var content strings.Builder
	// Preallocate some capacity to reduce reallocations
	content.Grow(len(tasks) * 50)

	// Group tasks by category (using map to avoid n² operations)
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

	return os.WriteFile(todoFilePath, []byte(content.String()), 0644)
}

// createBackup creates a backup of the current todo file
func createBackup() error {
	if backupDirectory == "" {
		return nil
	}

	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(backupDirectory, 0755); err != nil {
		return err
	}

	// Read existing file
	data, err := os.ReadFile(todoFilePath)
	if err != nil {
		return err
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupFile := filepath.Join(backupDirectory, fmt.Sprintf("TODO-%s.md", timestamp))

	// Write backup file
	if err := os.WriteFile(backupFile, data, 0644); err != nil {
		return err
	}

	// Clean up old backups
	if err := cleanupOldBackups(); err != nil {
		return err
	}

	return nil
}

// cleanupOldBackups removes old backups exceeding the maximum count
func cleanupOldBackups() error {
	if backupDirectory == "" || maxBackups <= 0 {
		return nil
	}

	// List backup files
	entries, err := os.ReadDir(backupDirectory)
	if err != nil {
		return err
	}

	// Filter for backup files and sort by modification time
	backupFiles := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "TODO-") && strings.HasSuffix(entry.Name(), ".md") {
			backupFiles = append(backupFiles, filepath.Join(backupDirectory, entry.Name()))
		}
	}

	// If we have more backups than the limit, delete the oldest ones
	if len(backupFiles) > maxBackups {
		// Sort backup files by modification time (oldest first)
		// We use a simple bubble sort here as the number of files is typically small
		for i := 0; i < len(backupFiles)-1; i++ {
			for j := 0; j < len(backupFiles)-i-1; j++ {
				info1, _ := os.Stat(backupFiles[j])
				info2, _ := os.Stat(backupFiles[j+1])
				if info1.ModTime().After(info2.ModTime()) {
					backupFiles[j], backupFiles[j+1] = backupFiles[j+1], backupFiles[j]
				}
			}
		}

		// Delete the oldest files
		for i := 0; i < len(backupFiles)-maxBackups; i++ {
			os.Remove(backupFiles[i])
		}
	}

	return nil
}
