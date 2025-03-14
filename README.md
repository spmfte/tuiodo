# tuiodo

A beautiful terminal-based TODO application with support for task categories and filtering.

## Features

- Interactive terminal UI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Elegant styling with [Lipgloss](https://github.com/charmbracelet/lipgloss)
- Markdown-based task persistence in a `TODO.md` file
- Support for task categories and filtering by category
- Keyboard-driven interface for efficient task management

## Usage

### Installation

```bash
go install github.com/spmfte/tuiodo@latest
```

Or clone and build:

```bash
git clone https://github.com/spmfte/tuiodo.git
cd tuiodo
go build
```

### Keyboard Controls

- `j` / `down arrow`: Move cursor down
- `k` / `up arrow`: Move cursor up
- `space`: Toggle task completion
- `a`: Add new task
- `d`: Delete selected task
- `c`: Cycle through category filters
- `q` or `Ctrl+C`: Quit

### Adding Tasks with Categories

When adding a new task, you can specify a category using the format:

```
Category: Task description
```

For example:
- "Work: Finish weekly report"
- "Personal: Call mom"

## File Format

Tasks are stored in a `TODO.md` file with the following format:

```markdown
## Category1

- [ ] Task description
- [x] Completed task

## Category2

- [ ] Another task
```

## License

MIT 