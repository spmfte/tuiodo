package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spmfte/tuiodo/handlers"
	"github.com/spmfte/tuiodo/model"
	"github.com/spmfte/tuiodo/storage"
	"github.com/spmfte/tuiodo/ui"
)

// App is the main application model
type App struct {
	model model.Model
}

// Init initializes the application
func (a App) Init() tea.Cmd {
	return nil
}

// Update processes messages and updates the model
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.model, cmd = handlers.Update(msg, a.model)
	return a, cmd
}

// View renders the application UI
func (a App) View() string {
	return ui.View(a.model)
}

func main() {
	// Load tasks from storage
	tasks := storage.LoadTasks()

	// Create initial model
	initialModel := model.NewModel(tasks)

	// Create application instance
	app := App{model: initialModel}

	// Create and run the program with alt screen enabled for full-terminal UI
	p := tea.NewProgram(app, tea.WithAltScreen())

	// Run the application
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
