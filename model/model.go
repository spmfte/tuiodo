package model

// Priority represents a task priority level
type Priority string

const (
	PriorityNone   Priority = ""
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Task represents a single TODO item
type Task struct {
	Description string
	Done        bool
	Category    string
	Priority    Priority
	DueDate     string // Format: "YYYY-MM-DD" or empty string
}

// TabView represents the current view/filter mode
type TabView string

const (
	TabAll       TabView = "all"
	TabToday     TabView = "today"
	TabPending   TabView = "pending"
	TabCompleted TabView = "completed"
	TabCategory  TabView = "category" // Filtered by specific category
)

// Model represents the application state
type Model struct {
	Tasks           []Task
	Cursor          int
	SelectedTasks   map[int]struct{}
	InputMode       bool
	Input           string
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
}

// Pagination tracks position in a paginated list
type Pagination struct {
	Page          int
	TotalPages    int
	ItemsPerPage  int
	CurrentOffset int
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

	priorities := []Priority{PriorityNone, PriorityLow, PriorityMedium, PriorityHigh}

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
}

// AddTask adds a new task to the model
func (m *Model) AddTask(description, category string) {
	if category != "" {
		m.Categories[category] = struct{}{}
	}

	m.Tasks = append(m.Tasks, Task{
		Description: description,
		Category:    category,
	})
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
}

// DeleteCurrentTask removes the task at the current cursor position
func (m *Model) DeleteCurrentTask() {
	if len(m.Tasks) == 0 || m.Cursor >= len(m.Tasks) {
		return
	}

	m.Tasks = append(m.Tasks[:m.Cursor], m.Tasks[m.Cursor+1:]...)
	if m.Cursor >= len(m.Tasks) && m.Cursor > 0 {
		m.Cursor--
	}
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

// MoveCursorUp moves the cursor up one position if possible
func (m *Model) MoveCursorUp() {
	if m.Cursor > 0 {
		m.Cursor--
	} else if m.Pagination.Page > 0 {
		// If at the top of the page, go to previous page
		m.PrevPage()
		visibleTasks := m.GetVisibleTasks()
		if len(visibleTasks) > 0 {
			m.Cursor = len(visibleTasks) - 1
		}
	}
}

// MoveCursorDown moves the cursor down one position if possible
func (m *Model) MoveCursorDown() {
	visibleTasks := m.GetVisibleTasks()
	if len(visibleTasks) > 0 && m.Cursor < len(visibleTasks)-1 {
		m.Cursor++
	} else if m.Pagination.Page < m.Pagination.TotalPages-1 {
		// If at the bottom of the page, go to next page
		m.NextPage()
		m.Cursor = 0
	}
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
