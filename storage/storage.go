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
var DefaultTodoFilePath string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		DefaultTodoFilePath = "TODO.md"
	} else {
		DefaultTodoFilePath = filepath.Join(homeDir, "TODO.md")
	}
}

// findGitRepository searches for a git repository starting from the current directory
// and walking up the directory tree
func findGitRepository(startPath string) (string, error) {
	currentPath := startPath

	for {
		// Check if current directory contains .git
		gitPath := filepath.Join(currentPath, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			// Found a git repository
			return currentPath, nil
		}

		// Move up to parent directory
		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			// Reached root directory, no git repository found
			return "", fmt.Errorf("no git repository found")
		}
		currentPath = parentPath
	}
}

// getGitRootTodoPath returns the path to TODO.md at the root of the git repository
func getGitRootTodoPath() (string, error) {
	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// Find git repository root
	gitRoot, err := findGitRepository(currentDir)
	if err != nil {
		return "", err
	}

	// Return path to TODO.md at git root
	return filepath.Join(gitRoot, "TODO.md"), nil
}

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
	priorityPattern  = regexp.MustCompile(`@priority:(high|medium|low|critical)`)
	createdAtPattern = regexp.MustCompile(`@created:(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)`)
	duePattern       = regexp.MustCompile(`@due:(\d{4}-\d{2}-\d{2})`)
	tagPattern       = regexp.MustCompile(`@tag:([^\s@]+)`)
	statusPattern    = regexp.MustCompile(`@status:([^\s@]+)`)
)

// Initialize sets up the storage with configurable settings
func Initialize(filePath string, backupDir string, maxBackupFiles int, enableAutoSave bool, enableBackup bool) {
	if filePath != "" {
		// If path is not absolute, make it absolute from current directory
		if !filepath.IsAbs(filePath) {
			if absPath, err := filepath.Abs(filePath); err == nil {
				filePath = absPath
			}
		}
		todoFilePath = filePath
	} else {
		// No explicit path provided, try to find git repository first
		if gitPath, err := getGitRootTodoPath(); err == nil {
			todoFilePath = gitPath
		} else {
			// Fall back to current working directory
			if currentDir, err := os.Getwd(); err == nil {
				todoFilePath = filepath.Join(currentDir, "TODO.md")
			} else {
				// Last resort: home directory
				todoFilePath = DefaultTodoFilePath
			}
		}
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
				case "critical":
					priority = model.PriorityCritical
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

			// Extract creation date if present
			createdAt := time.Now() // Default to now if not found
			createdAtMatch := createdAtPattern.FindStringSubmatch(description)
			if len(createdAtMatch) > 1 {
				if parsedTime, err := time.Parse(time.RFC3339, createdAtMatch[1]); err == nil {
					createdAt = parsedTime
				}
				// Remove the created tag from description
				description = strings.TrimSpace(createdAtPattern.ReplaceAllString(description, ""))
			}

			// Extract due date if present
			var dueDate string
			dueMatch := duePattern.FindStringSubmatch(description)
			if len(dueMatch) > 1 {
				dueDate = dueMatch[1]
				// Keep the due date in the metadata but remove it from visible description
				description = strings.TrimSpace(duePattern.ReplaceAllString(description, ""))
			}

			// Extract tags if present
			var tags []string
			tagMatches := tagPattern.FindAllStringSubmatch(description, -1)
			for _, match := range tagMatches {
				if len(match) > 1 {
					tags = append(tags, match[1])
				}
			}
			// Remove the tag markers from visible description
			description = strings.TrimSpace(tagPattern.ReplaceAllString(description, ""))

			// Extract status if present
			var status string
			statusMatch := statusPattern.FindStringSubmatch(description)
			if len(statusMatch) > 1 {
				status = statusMatch[1]
				// Remove the status tag from visible description
				description = strings.TrimSpace(statusPattern.ReplaceAllString(description, ""))
			}

			// Create task object with all metadata
			task := model.Task{
				Description: description,
				Done:        isDone,
				Category:    currentCategory,
				Priority:    priority,
				CreatedAt:   createdAt,
				Metadata:    make(map[string]string),
			}

			// Store additional metadata
			if dueDate != "" {
				task.Metadata["due"] = dueDate
			}
			if len(tags) > 0 {
				task.Metadata["tags"] = strings.Join(tags, ",")
			}
			if status != "" {
				task.Metadata["status"] = status
			}

			tasks = append(tasks, task)
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

			// Add creation date
			description = fmt.Sprintf("%s @created:%s", description, task.CreatedAt.UTC().Format(time.RFC3339))

			// Add additional metadata if present
			if dueDate, ok := task.Metadata["due"]; ok {
				description = fmt.Sprintf("%s @due:%s", description, dueDate)
			}
			if tags, ok := task.Metadata["tags"]; ok {
				for _, tag := range strings.Split(tags, ",") {
					description = fmt.Sprintf("%s @tag:%s", description, tag)
				}
			}
			if status, ok := task.Metadata["status"]; ok {
				description = fmt.Sprintf("%s @status:%s", description, status)
			}

			content.WriteString(fmt.Sprintf("- [%s] %s\n", checkmark, description))
		}
		content.WriteString("\n")
	}

	// Add an invisible tag at the bottom for the tuiodo app
	// This tag doesn't appear in the TUI but will be visible in raw markdown files
	content.WriteString("\n<!-- Optimized for [tuiodo](https://github.com/spmfte/tuiodo) -->\n")

	// Ensure the directory exists
	dir := filepath.Dir(todoFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
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
