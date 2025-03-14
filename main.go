package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const todoFilePath = "TODO.md"

// Task represents a single TODO item
type Task struct {
	Description string
	Done        bool
	Category    string
}

// Model represents the application state
type Model struct {
	Tasks         []Task
	Cursor        int
	SelectedTasks map[int]struct{}
	InputMode     bool
	Input         string
	Categories    map[string]struct{}
	CurrentFilter string
	Width         int
	Height        int
}

// Initial model
func initialModel() Model {
	tasks := loadTasksFromFile()
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
	}
}

// Load tasks from TODO.md file
func loadTasksFromFile() []Task {
	content, err := os.ReadFile(todoFilePath)
	if err != nil {
		// File doesn't exist or can't be read, return empty task list
		return []Task{}
	}

	lines := strings.Split(string(content), "\n")
	tasks := []Task{}

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

			tasks = append(tasks, Task{
				Description: description,
				Done:        isDone,
				Category:    currentCategory,
			})
		}
	}

	return tasks
}

// Save tasks to TODO.md file
func saveTasksToFile(tasks []Task) error {
	var content strings.Builder

	// Group tasks by category
	categorizedTasks := make(map[string][]Task)
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
			content.WriteString(fmt.Sprintf("- [%s] %s\n", checkmark, task.Description))
		}
		content.WriteString("\n")
	}

	return os.WriteFile(todoFilePath, []byte(content.String()), 0644)
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypress(msg)
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	}
	return m, nil
}

func (m Model) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.InputMode {
		return m.handleInputKeypress(msg)
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < len(m.Tasks)-1 {
			m.Cursor++
		}
	case "a": // Add new task
		m.InputMode = true
		m.Input = ""
	case "d": // Delete task
		if len(m.Tasks) > 0 {
			m.Tasks = append(m.Tasks[:m.Cursor], m.Tasks[m.Cursor+1:]...)
			if m.Cursor >= len(m.Tasks) && m.Cursor > 0 {
				m.Cursor--
			}
			saveTasksToFile(m.Tasks)
		}
	case "space": // Toggle task completion
		if len(m.Tasks) > 0 {
			m.Tasks[m.Cursor].Done = !m.Tasks[m.Cursor].Done
			saveTasksToFile(m.Tasks)
		}
	case "c": // Cycle through categories for filtering
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
	}

	return m, nil
}

func (m Model) handleInputKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if strings.TrimSpace(m.Input) != "" {
			// Extract category if included in format "Category: Task description"
			category := ""
			description := m.Input

			if parts := strings.SplitN(m.Input, ":", 2); len(parts) == 2 {
				category = strings.TrimSpace(parts[0])
				description = strings.TrimSpace(parts[1])

				if category != "" {
					m.Categories[category] = struct{}{}
				}
			}

			m.Tasks = append(m.Tasks, Task{
				Description: description,
				Category:    category,
			})
			saveTasksToFile(m.Tasks)
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
	default:
		if len(msg.String()) == 1 {
			m.Input += msg.String()
		}
	}

	return m, nil
}

func (m Model) View() string {
	var s strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	s.WriteString(titleStyle.Render(" TUIODO "))

	if m.CurrentFilter != "" {
		filterStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
		s.WriteString(" ")
		s.WriteString(filterStyle.Render("Filtered by: " + m.CurrentFilter))
	}

	s.WriteString("\n\n")

	// Show help
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	s.WriteString(helpStyle.Render("j/k: navigate • space: toggle • a: add • d: delete • c: filter • q: quit\n\n"))

	// Input mode
	if m.InputMode {
		inputPrompt := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Render("New task (optional Format: Category: Task) > ")
		s.WriteString(inputPrompt + m.Input + "▋\n\n")
		return s.String()
	}

	// Filter tasks based on current filter
	filteredTasks := m.Tasks
	if m.CurrentFilter != "" {
		filteredTasks = []Task{}
		for _, task := range m.Tasks {
			if task.Category == m.CurrentFilter {
				filteredTasks = append(filteredTasks, task)
			}
		}
	}

	// No tasks message
	if len(filteredTasks) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
		if m.CurrentFilter != "" {
			s.WriteString(emptyStyle.Render("No tasks in category '" + m.CurrentFilter + "'\n"))
		} else {
			s.WriteString(emptyStyle.Render("No tasks yet. Press 'a' to add one.\n"))
		}
		return s.String()
	}

	// List tasks
	for i, task := range filteredTasks {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		checkBox := "[ ]"
		if task.Done {
			checkBox = "[x]"
		}

		taskStyle := lipgloss.NewStyle()
		if task.Done {
			taskStyle = taskStyle.Strikethrough(true).Foreground(lipgloss.Color("#626262"))
		}

		cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
		checkBoxStyle := lipgloss.NewStyle()
		if task.Done {
			checkBoxStyle = checkBoxStyle.Foreground(lipgloss.Color("#00CC00"))
		}

		categoryText := ""
		if task.Category != "" {
			categoryStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
			categoryText = " " + categoryStyle.Render(task.Category)
		}

		s.WriteString(fmt.Sprintf("%s %s %s%s\n",
			cursorStyle.Render(cursor),
			checkBoxStyle.Render(checkBox),
			taskStyle.Render(task.Description),
			categoryText,
		))
	}

	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
