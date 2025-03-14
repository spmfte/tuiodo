# tuiodo

A beautiful, intuitive terminal-based TODO application with modern UI design and powerful task management features.

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  <img src="https://img.shields.io/badge/UI-Bubble%20Tea-ff69b4" alt="UI Framework">
</p>

## Features

- **Beautiful Modern UI** - Clean, intuitive interface with clear visual hierarchy
- **Task Management** - Add, edit, complete, and delete tasks with ease
- **Task Prioritization** - Set and visualize task priorities (low/medium/high)
- **Category Support** - Organize tasks by custom categories
- **Multiple Views** - Filter tasks by status (all/pending/completed)
- **Keyboard-driven** - Fast and efficient workflow with intuitive shortcuts
- **Persistence** - Markdown-based storage in a simple `TODO.md` file
- **Pagination** - Handle large task lists with smart page navigation
- **Context-aware Help** - Comprehensive help screen (press `?`)

## Screenshots

[Screenshots would be here]

## Installation

```bash
# Via Go install
go install github.com/spmfte/tuiodo@latest

# Or clone and build
git clone https://github.com/spmfte/tuiodo.git
cd tuiodo
go build
```

## Keyboard Controls

### Navigation
- `j` / `↓`: Move cursor down
- `k` / `↑`: Move cursor up
- `tab` / `t`: Switch between tabs (All/Pending/Completed)
- `n` / `→` / `l`: Next page
- `b` / `←` / `h`: Previous page

### Task Management
- `a`: Add new task
- `e`: Edit selected task
- `d`: Delete selected task
- `space` / `enter`: Toggle task completion
- `p`: Cycle through priorities (none/low/medium/high)

### Filtering
- `c`: Cycle through category filters
- `tab` / `t`: Switch between views (All/Pending/Completed)

### Other
- `?` / `F1`: Show/hide help
- `q` / `Ctrl+C`: Quit

## Task Format

When adding tasks, you can specify categories using the format:

```
Category: Task description
```

For example:
- "Work: Finish quarterly report"
- "Personal: Schedule dentist appointment"

## File Format

Tasks are stored in a `TODO.md` file with metadata for priorities and due dates:

```markdown
## Work

- [ ] Prepare presentation @priority:high
- [x] Send weekly update email @priority:medium

## Personal

- [ ] Buy groceries @priority:low
- [ ] Call mom @due:2023-05-10
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT 