package model

import (
	"regexp"
	"sort"
	"strings"
	"time"
)

// Priority represents a task priority level
type Priority string

const (
	PriorityNone     Priority = ""
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// Task represents a single TODO item
type Task struct {
	Description string
	Done        bool
	Category    string
	Priority    Priority
	CreatedAt   time.Time         // When the task was created
	Archived    bool              // Whether the task is archived
	Metadata    map[string]string // Additional metadata like due dates, tags, status
}

// TabView represents the current view/filter mode
type TabView string

const (
	TabAll       TabView = "all"
	TabToday     TabView = "today"
	TabPending   TabView = "pending"
	TabCompleted TabView = "completed"
	TabCategory  TabView = "category" // Filtered by specific category
	TabArchived  TabView = "archived" // Show archived tasks
)

// SortType represents different ways to sort tasks
type SortType string

const (
	SortByPriority  SortType = "priority"
	SortByCreatedAt SortType = "created"
	SortByCategory  SortType = "category"
)

// Model represents the application state
type Model struct {
	Tasks           []Task
	Cursor          int
	SelectedTasks   map[int]struct{}
	InputMode       bool
	Input           string
	InputCursor     int // Position of cursor within input field
	Categories      map[string]struct{}
	CurrentFilter   string  // Category filter
	CurrentView     TabView // Current tab view
	CurrentCategory string  // When in TabCategory
	Width           int
	Height          int
	Pagination      Pagination
	StatusMessage   string // Temporary status messages
	EditingTask     bool   // Whether currently editing a task
	EditingTaskIdx  int    // Index of task being edited
	HelpVisible     bool   // Whether help is visible
	TaskExpanded    bool   // Whether task details are expanded
	ExpandedTaskIdx int    // Index of task being expanded
	DeleteConfirm   bool   // Whether delete confirmation is active
	LastDeleted     *Task  // Last deleted task for undo
	LastDeletedIdx  int    // Index where the task was deleted
}

// Pagination tracks position in a paginated list
type Pagination struct {
	Page          int
	TotalPages    int
	ItemsPerPage  int
	CurrentOffset int
}

// priorityValue maps priorities to numeric values for sorting
var priorityValue = map[Priority]int{
	PriorityCritical: 5,
	PriorityHigh:     4,
	PriorityMedium:   3,
	PriorityLow:      2,
	PriorityNone:     1,
}

// NewModel creates a new model with initial state
func NewModel(tasks []Task) Model {
	categories := make(map[string]struct{})

	// Extract all categories from tasks
	for _, task := range tasks {
		if task.Category != "" {
			categories[task.Category] = struct{}{}
		}
	}

	return Model{
		Tasks:         tasks,
		SelectedTasks: make(map[int]struct{}),
		InputMode:     false,
		Categories:    categories,
		CurrentFilter: "", // Empty string means no filter
		CurrentView:   TabAll,
		Pagination: Pagination{
			Page:         0,
			ItemsPerPage: 10,
		},
		HelpVisible: false,
	}
}

// NewModelWithConfig creates a new model with configuration options
func NewModelWithConfig(tasks []Task, tasksPerPage int, defaultCategory string, keyBindings interface{}) Model {
	// Create a basic model first
	m := NewModel(tasks)

	// Apply configuration
	if tasksPerPage > 0 {
		m.Pagination.ItemsPerPage = tasksPerPage
	}

	// Set default filter to the default category if specified
	if defaultCategory != "" {
		// First check if this category exists
		if _, exists := m.Categories[defaultCategory]; exists {
			m.CurrentFilter = defaultCategory
			m.CurrentView = TabCategory
			m.CurrentCategory = defaultCategory
		}
	}

	// Key bindings would be used in the handlers package, but we'll store them
	// in the model for now (or they could be stored in a global variable)

	// Recalculate pagination based on the configured settings
	m.recalculatePagination()

	return m
}

// GetFilteredTasks returns tasks filtered by the current tab view and filters
func (m Model) GetFilteredTasks() []Task {
	var filteredTasks []Task

	// First apply tab view filter
	switch m.CurrentView {
	case TabAll:
		// Show all tasks including archived
		filteredTasks = m.Tasks
	case TabToday:
		// TODO: Implement date filtering when we add due dates
		filteredTasks = m.Tasks
	case TabPending:
		for _, task := range m.Tasks {
			if !task.Done {
				filteredTasks = append(filteredTasks, task)
			}
		}
	case TabCompleted:
		for _, task := range m.Tasks {
			if task.Done {
				filteredTasks = append(filteredTasks, task)
			}
		}
	case TabCategory:
		if m.CurrentCategory != "" {
			for _, task := range m.Tasks {
				if task.Category == m.CurrentCategory {
					filteredTasks = append(filteredTasks, task)
				}
			}
		} else {
			filteredTasks = m.Tasks
		}
	default:
		filteredTasks = m.Tasks
	}

	// Then apply category filter if set
	if m.CurrentFilter != "" && m.CurrentView != TabCategory {
		tasksWithCategory := []Task{}
		for _, task := range filteredTasks {
			if task.Category == m.CurrentFilter {
				tasksWithCategory = append(tasksWithCategory, task)
			}
		}
		filteredTasks = tasksWithCategory
	}

	return filteredTasks
}

// GetVisibleTasks returns only the tasks that should be displayed on the current page
func (m Model) GetVisibleTasks() []Task {
	allFiltered := m.GetFilteredTasks()

	// If pagination disabled or unnecessary, return all tasks
	if m.Pagination.ItemsPerPage <= 0 || len(allFiltered) <= m.Pagination.ItemsPerPage {
		return allFiltered
	}

	// Calculate start and end indices for the current page
	startIdx := m.Pagination.Page * m.Pagination.ItemsPerPage
	endIdx := startIdx + m.Pagination.ItemsPerPage

	// Make sure we don't go out of bounds
	if startIdx >= len(allFiltered) {
		startIdx = 0
		m.Pagination.Page = 0
	}
	if endIdx > len(allFiltered) {
		endIdx = len(allFiltered)
	}

	return allFiltered[startIdx:endIdx]
}

// CycleCategory advances to the next category in the list
func (m *Model) CycleCategory() {
	categories := make([]string, 0, len(m.Categories)+1)
	categories = append(categories, "") // Empty filter (show all)

	for category := range m.Categories {
		categories = append(categories, category)
	}

	// Find current filter position
	currentIndex := 0
	for i, cat := range categories {
		if cat == m.CurrentFilter {
			currentIndex = i
			break
		}
	}

	// Cycle to next category
	nextIndex := (currentIndex + 1) % len(categories)
	m.CurrentFilter = categories[nextIndex]

	// Reset cursor and recalculate pagination when changing category
	m.Cursor = 0
	m.recalculatePagination()
}

// CycleTab changes to the next tab view
func (m *Model) CycleTab() {
	tabs := []TabView{TabAll, TabPending, TabCompleted}

	// Find current tab position
	currentIndex := 0
	for i, tab := range tabs {
		if tab == m.CurrentView {
			currentIndex = i
			break
		}
	}

	// Cycle to next tab
	nextIndex := (currentIndex + 1) % len(tabs)
	m.CurrentView = tabs[nextIndex]

	// Reset cursor position and recalculate pagination
	m.Cursor = 0
	m.recalculatePagination()
}

// CyclePriority cycles through priority levels for the current task
func (m *Model) CyclePriority() {
	if len(m.Tasks) == 0 || m.Cursor >= len(m.Tasks) {
		return
	}

	priorities := []Priority{PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical}

	// Find current priority position
	currentIndex := 0
	for i, priority := range priorities {
		if priority == m.Tasks[m.Cursor].Priority {
			currentIndex = i
			break
		}
	}

	// Cycle to next priority
	nextIndex := (currentIndex + 1) % len(priorities)
	m.Tasks[m.Cursor].Priority = priorities[nextIndex]

	// Auto-sort by priority to update UI immediately
	m.SortTasks(SortByPriority)

	// Find the task in the new sorted position and update cursor
	taskToFind := m.Tasks[m.Cursor]
	for i, task := range m.Tasks {
		if task.Description == taskToFind.Description &&
			task.Category == taskToFind.Category &&
			task.CreatedAt == taskToFind.CreatedAt {
			m.Cursor = i
			break
		}
	}

	// Recalculate pagination after sorting
	m.recalculatePagination()
}

// AddTask adds a new task to the model
func (m *Model) AddTask(description, category string, priority Priority) {
	if category != "" {
		m.Categories[category] = struct{}{}
	}

	// Default to low priority if none specified
	if priority == PriorityNone {
		priority = PriorityLow
	}

	m.Tasks = append(m.Tasks, Task{
		Description: description,
		Category:    category,
		Priority:    priority,
		CreatedAt:   time.Now(),
	})

	// Auto-sort by priority to ensure proper positioning
	m.SortTasks(SortByPriority)
}

// UpdateTask updates an existing task
func (m *Model) UpdateTask(index int, description, category string, priority Priority) {
	if index < 0 || index >= len(m.Tasks) {
		return
	}

	if category != "" && category != m.Tasks[index].Category {
		m.Categories[category] = struct{}{}
	}

	m.Tasks[index].Description = description
	m.Tasks[index].Category = category
	m.Tasks[index].Priority = priority

	// Auto-sort by priority to ensure proper positioning
	m.SortTasks(SortByPriority)
}

// DeleteCurrentTask deletes the task at the current cursor position
func (m *Model) DeleteCurrentTask() {
	filteredTasks := m.GetVisibleTasks()
	if len(filteredTasks) == 0 || m.Cursor >= len(filteredTasks) {
		return
	}

	// Find the actual index in the Tasks slice
	taskToDelete := filteredTasks[m.Cursor]
	var taskIdx int
	found := false

	for i, task := range m.Tasks {
		if task.Description == taskToDelete.Description &&
			task.Category == taskToDelete.Category &&
			task.CreatedAt == taskToDelete.CreatedAt {
			taskIdx = i
			found = true
			break
		}
	}

	if !found {
		return // Task not found in the main list
	}

	// Store the deleted task for potential undo
	deletedTask := m.Tasks[taskIdx]
	m.LastDeleted = &deletedTask
	m.LastDeletedIdx = taskIdx

	// Remove the task
	m.Tasks = append(m.Tasks[:taskIdx], m.Tasks[taskIdx+1:]...)

	// Update cursor if needed
	if m.Cursor >= len(filteredTasks)-1 && m.Cursor > 0 {
		m.Cursor--
	}
}

// UndoDelete restores the last deleted task
func (m *Model) UndoDelete() bool {
	if m.LastDeleted == nil {
		return false
	}

	// Calculate insertion point (either the original position or end of list)
	insertIdx := m.LastDeletedIdx
	if insertIdx > len(m.Tasks) {
		insertIdx = len(m.Tasks)
	}

	// Insert task back into original position
	if insertIdx == len(m.Tasks) {
		m.Tasks = append(m.Tasks, *m.LastDeleted)
	} else {
		// Create space and insert the task
		m.Tasks = append(m.Tasks, Task{})                // Add empty task at end
		copy(m.Tasks[insertIdx+1:], m.Tasks[insertIdx:]) // Shift tasks right
		m.Tasks[insertIdx] = *m.LastDeleted              // Insert task
	}

	// Clean up
	m.LastDeleted = nil
	m.RecalculatePagination()
	return true
}

// ToggleCurrentTask toggles the completion status of the current task
func (m *Model) ToggleCurrentTask() {
	if len(m.Tasks) == 0 || m.Cursor >= len(m.Tasks) {
		return
	}

	m.Tasks[m.Cursor].Done = !m.Tasks[m.Cursor].Done
}

// ToggleHelp shows or hides the help screen
func (m *Model) ToggleHelp() {
	m.HelpVisible = !m.HelpVisible
}

// NextPage advances to the next page if possible
func (m *Model) NextPage() {
	if m.Pagination.Page < m.Pagination.TotalPages-1 {
		m.Pagination.Page++
		m.Cursor = 0
	}
}

// PrevPage goes to the previous page if possible
func (m *Model) PrevPage() {
	if m.Pagination.Page > 0 {
		m.Pagination.Page--
		m.Cursor = 0
	}
}

// MoveCursorUp moves the cursor up one position, or wraps to the bottom if at the top
func (m *Model) MoveCursorUp() {
	visibleTasks := m.GetVisibleTasks()
	if len(visibleTasks) == 0 {
		return
	}

	if m.Cursor > 0 {
		m.Cursor--
	} else {
		// If at the top, wrap to the bottom
		m.Cursor = len(visibleTasks) - 1
	}
}

// MoveCursorDown moves the cursor down one position, or wraps to the top if at the bottom
func (m *Model) MoveCursorDown() {
	visibleTasks := m.GetVisibleTasks()
	if len(visibleTasks) == 0 {
		return
	}

	if m.Cursor < len(visibleTasks)-1 {
		m.Cursor++
	} else {
		// If at the bottom, wrap to the top
		m.Cursor = 0
	}
}

// MoveInputCursorLeft moves the cursor left within the input field
func (m *Model) MoveInputCursorLeft() {
	if m.InputCursor > 0 {
		m.InputCursor--
	}
}

// MoveInputCursorRight moves the cursor right within the input field
func (m *Model) MoveInputCursorRight() {
	if m.InputCursor < len(m.Input) {
		m.InputCursor++
	}
}

// MoveInputCursorToStart moves the cursor to the beginning of the input field
func (m *Model) MoveInputCursorToStart() {
	m.InputCursor = 0
}

// MoveInputCursorToEnd moves the cursor to the end of the input field
func (m *Model) MoveInputCursorToEnd() {
	m.InputCursor = len(m.Input)
}

// InsertTextAtCursor inserts text at the current cursor position
func (m *Model) InsertTextAtCursor(text string) {
	if m.InputCursor == len(m.Input) {
		// At the end, just append
		m.Input += text
	} else {
		// Insert in the middle
		m.Input = m.Input[:m.InputCursor] + text + m.Input[m.InputCursor:]
	}
	m.InputCursor += len(text)
}

// DeleteTextAtCursor deletes one character at the current cursor position
func (m *Model) DeleteTextAtCursor() {
	if m.InputCursor < len(m.Input) {
		m.Input = m.Input[:m.InputCursor] + m.Input[m.InputCursor+1:]
	}
}

// DeleteTextBeforeCursor deletes one character before the current cursor position
func (m *Model) DeleteTextBeforeCursor() {
	if m.InputCursor > 0 {
		m.Input = m.Input[:m.InputCursor-1] + m.Input[m.InputCursor:]
		m.InputCursor--
	}
}

// ArchiveCurrentTask archives the task at the current cursor position
func (m *Model) ArchiveCurrentTask() {
	filteredTasks := m.GetVisibleTasks()
	if len(filteredTasks) == 0 || m.Cursor >= len(filteredTasks) {
		return
	}

	// Find the actual index in the Tasks slice
	taskToArchive := filteredTasks[m.Cursor]
	var taskIdx int
	found := false

	for i, task := range m.Tasks {
		if task.Description == taskToArchive.Description &&
			task.Category == taskToArchive.Category &&
			task.CreatedAt == taskToArchive.CreatedAt {
			taskIdx = i
			found = true
			break
		}
	}

	if !found {
		return // Task not found in the main list
	}

	// Archive the task
	m.Tasks[taskIdx].Archived = true

	// Recalculate pagination after archiving
	m.recalculatePagination()
}

// UnarchiveTask unarchives a task by its index
func (m *Model) UnarchiveTask(index int) {
	if index < 0 || index >= len(m.Tasks) {
		return
	}
	m.Tasks[index].Archived = false
	m.recalculatePagination()
}

// getTaskStatus returns a numeric status for sorting (lower = higher priority)
func getTaskStatus(task Task) int {
	if task.Archived {
		return 2 // Archived tasks (middle priority)
	} else if task.Done {
		return 3 // Completed tasks (lowest priority)
	}
	return 1 // Active tasks (highest priority)
}

// ParsePriorityFromText extracts priority from text and returns cleaned text and priority
func ParsePriorityFromText(text string) (string, Priority) {
	// Regex to match @priority:value
	priorityRegex := regexp.MustCompile(`@priority:(high|medium|low|critical)`)

	match := priorityRegex.FindStringSubmatch(text)
	if len(match) > 1 {
		priorityStr := match[1]
		var priority Priority
		switch priorityStr {
		case "critical":
			priority = PriorityCritical
		case "high":
			priority = PriorityHigh
		case "medium":
			priority = PriorityMedium
		case "low":
			priority = PriorityLow
		default:
			priority = PriorityLow
		}

		// Remove the priority tag from text
		cleanedText := strings.TrimSpace(priorityRegex.ReplaceAllString(text, ""))
		return cleanedText, priority
	}

	// No priority found, return original text with default priority
	return text, PriorityLow
}

// UpdateWindowSize updates the stored window dimensions
func (m *Model) UpdateWindowSize(width, height int) {
	m.Width = width
	m.Height = height

	// Adjust items per page based on available height
	// We need to account for:
	// - Title bar (1 line)
	// - Tab bar (3 lines with margins)
	// - Status bar (1 line)
	// - Container borders (2 lines)
	// - Task header (1 line)
	// - Potential pagination info (2 lines)
	// - Each task takes approximately 2 lines (with separator)
	availableHeight := height - 10

	// Ensure a reasonable minimum
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Divide by 2 because each task + separator takes ~2 lines
	m.Pagination.ItemsPerPage = availableHeight / 2

	// Ensure at least 3 tasks per page
	if m.Pagination.ItemsPerPage < 3 {
		m.Pagination.ItemsPerPage = 3
	}

	// Recalculate pagination for current view
	m.recalculatePagination()
}

// recalculatePagination updates the pagination based on current filtered tasks
func (m *Model) recalculatePagination() {
	totalTasks := len(m.GetFilteredTasks())
	if m.Pagination.ItemsPerPage > 0 {
		m.Pagination.TotalPages = (totalTasks + m.Pagination.ItemsPerPage - 1) / m.Pagination.ItemsPerPage

		// Make sure we're not on a non-existent page
		if m.Pagination.Page >= m.Pagination.TotalPages && m.Pagination.TotalPages > 0 {
			m.Pagination.Page = m.Pagination.TotalPages - 1
		}

		// If we have no pages, reset to page 0
		if m.Pagination.TotalPages == 0 {
			m.Pagination.Page = 0
		}
	}
}

// RecalculatePagination updates the pagination state based on current filter
func (m *Model) RecalculatePagination() {
	m.recalculatePagination()
}

// SetStatus sets a temporary status message
func (m *Model) SetStatus(message string) {
	m.StatusMessage = message
}

// SortTasks sorts the tasks based on the specified sort type
func (m *Model) SortTasks(sortType SortType) {
	switch sortType {
	case SortByPriority:
		sort.SliceStable(m.Tasks, func(i, j int) bool {
			// First sort by status: Active > Archived > Completed
			statusI := getTaskStatus(m.Tasks[i])
			statusJ := getTaskStatus(m.Tasks[j])
			if statusI != statusJ {
				return statusI < statusJ // Lower status number = higher priority
			}
			// Within same status group, sort by priority
			return priorityValue[m.Tasks[i].Priority] > priorityValue[m.Tasks[j].Priority]
		})
	case SortByCreatedAt:
		sort.SliceStable(m.Tasks, func(i, j int) bool {
			// First sort by status: Active > Archived > Completed
			statusI := getTaskStatus(m.Tasks[i])
			statusJ := getTaskStatus(m.Tasks[j])
			if statusI != statusJ {
				return statusI < statusJ // Lower status number = higher priority
			}
			// Within same status group, sort by creation date
			return m.Tasks[i].CreatedAt.After(m.Tasks[j].CreatedAt)
		})
	case SortByCategory:
		sort.SliceStable(m.Tasks, func(i, j int) bool {
			// First sort by status: Active > Archived > Completed
			statusI := getTaskStatus(m.Tasks[i])
			statusJ := getTaskStatus(m.Tasks[j])
			if statusI != statusJ {
				return statusI < statusJ // Lower status number = higher priority
			}
			// Within same status group, sort by category
			return m.Tasks[i].Category < m.Tasks[j].Category
		})
	}
}
