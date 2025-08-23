# T<sub><sub>(UI)</sub></sub>ODO - A Modern Terminal Task Manager

TUIODO is a feature-rich, terminal-based task management application built with Go. It provides an intuitive, keyboard-driven interface with mouse support, making it easy to manage your tasks directly from the terminal.

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-blue?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/UI-Bubble%20Tea-ff69b4?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCA1MTIgNTEyIj48cGF0aCBkPSJNMTk3LjEgNDUzLjVjLTEzLS41LTI3LjItMS4yLTM4LjgtMy4yQzk0LjUgNDM5IDU3LjQgNDA0IDQ4IDM2Ny4yYy0xNy0xMy4yLTI0LjctMTkuMS0zOC45LTMwLjRDMTcgMzM4LjkgNi44IDM0MyA0LjQgMzU1LjljLTMuMyAxNi42IDkgMzIuNCAyNS4zIDMyLjQgMS44IDAgMy42LS4yIDUuMy0uNSA1LjItLjkgMTAuMy0zLjQgMTQuMy02LjkgMi44LTIuNCAxLjcgMCAxLjcgMHMtOS41IDE2LTE0LjggMjQuNmMtMS41IDIuNS0zIDUtNC4zIDcuNy0zIDYuNS0zLjQgMTMuOC0xLjEgMjAuNiA0LjEgMTIuMS0uNSAyOS41IDE0LjEgMzYuMiAxMiA1LjQgMjYuOC0uNSAzMy44LTExLjIgNC0xLjcgOC0zLjIgMTEuOS01LjMgMTguMy05LjUgMzMuMy0yMy45IDUwLjctMzMuNyAxMi45LTcuMiAyNi40LTEzLjkgNDAuNS0xOC4yIDE4LjItNS41IDM3LjUtNy4yIDU2LjUtOC42IDE5LjIgMCAzOC40LS4xIDU3LjYgMHYtM2MtMTkuMi0uMi0zOC41LS4xLTU3LjctLjEtMTktLjEtMzguNy43LTU3LjUgMy45LTIwLjMgMy40LTQwIDkuNy01OC42IDE3LjQtMTYuNSA2LjktMzEuOCAxNi4zLTQ2LjIgMjcuMi0xLjIuOS0yLjQgMS44LTMuNiAyLjggMy4xLTUuNyA2LjYtMTQuMyA3LjktMTguNC0xLjMgMi44IC4yIDIuNCAyLjMuOSA4LjQtNi4yIDE3LjEtMTEuOSAyNi44LTE2LjQgOC44LTQgMTguMS03LjQgMjcuNi05LjQgMTkuMy00IDM5LjEtNi4yIDU4LjctNi4yIDQwLjIgMCA4Mi4xIDAgMTIyLjctLjIgMTAuNyAwIDIxLjMtLjIgMzItLjQgNDYuNy03LjggODYuNy01Mi4yIDkzLjUtOTkuMiAyLjctMTguOS4xLTM4LjItNy0zOC40LTEwLjYgNC43LTExLjEgMzcuOC0xMS4xIDM3LjhzLTM4IDc5LjkxLTEyNy4yIDcyLjQxYy04OS4yLTcuNS0xNDUtNDUuODEtMTYxLjMtOTUuNzEtMTYuMi00OS45IDE3LjEtMTAxLjM5IDI3LjMtMTI3LjE5IDEwLjItMjUuOCAxMy43LTYwLjQtMjMuMS03Mi4xMi0zNi43NS0xMS43Mi03Mi4xNSAyNi44Mi01Ni41NSA3OS4xMiAxNS42IDUyLjI5IDEwLjA5IDk1Ljc5IDEwLjA5IDk1Ljc5Uzk4LjUgMzIwLjggODMuMSAzNDEuMmMtMTUuNCAyMC40LTYuOSAyNS4xLTYuOSAyNS4xcy0xNCAyLjctMjEuNSA5LjdjLTcuNSA3LTExLjQgMTYuNy0xMS40IDI2LjggMCAyMCA4LjcgMzUuNSAzNS41IDQ3LjEgOS42IDQuMSAyMS4xIDMuNiAzMC41IDMuOGgxMTIuNHYtMy4xYy03LjktLjEtMTUuOS0uMi0yNC42LjhoLjF6Ii8+PC9zdmc+" alt="UI Framework">
</p>

<div align="center">
  <img src="./tuiodo.demo.gif" alt="TUIODO Demo" width="800px"/>
</div>

<p align="center">
  <b>A beautiful, lightning-fast terminal task manager with extensive customization.</b>
</p>

## Features

### New UI Improvements
- **Dynamic Category Colors** - Customize colors for each task category
- **Priority Visibility** - Priority indicators automatically hidden for completed tasks
- **Intelligent Sorting** - Completed tasks always appear at the bottom regardless of sort order
- **Visual Consistency** - Categories for completed tasks appear with dimmed colors


### Core Features

- **Gorgeous Modern UI** with clean typography and visual hierarchy
- **Brilliant Colors** with monochrome mode support (`--no-color`)
- **Mouse Support** with hover effects and click interactions (where supported)
- **Keyboard-Driven** workflow with intuitive, vim-inspired shortcuts
- **Instant Performance** with optimized rendering and caching
- **Comprehensive CLI** with extensive configuration options

### Task Management

- **Priority Levels** (Critical/High/Medium/Low) with color coding
- **Dynamic Category Organization** with configurable colors
- **Smart Filtering** by status, category, and priority
- **Task Expansion** to view full details of any task
- **Circular Navigation** with wrap-around cursor movement
- **Delete Confirmation** with undo capability
- **Rich Metadata Support** using @tag notation
- **Markdown Storage** in simple, human-readable format (`~/TODO.md` by default)
- **Automatic Backups** with configurable options

### Metadata Tags Support

TUIODO supports several metadata tags for enhanced task tracking:

- **@priority** - Set task importance (critical/high/medium/low)
- **@created** - Automatically tracked creation timestamp
- **@due** - Set deadlines with YYYY-MM-DD format 
- **@tag** - Add custom tags to group related tasks
- **@status** - Track custom status values

### Content & Storage

- **Automatic Backups** with configurable options
- **Git Repository Detection** - Automatically uses TODO.md at git repository root when available
- **Multi-device Sync** via configurable storage paths (share tasks via Dropbox, etc.)
- **Import/Export** to standard formats (coming soon)

### Advanced Capabilities

- **Extensive Configuration** via YAML configuration files
- **Theme Support** with pre-built and custom themes
- **Custom Keybindings** to match your workflow
- **Command Palette** for quick access to all features (coming soon)
- **Plugin System** for extending functionality (coming soon)

## Screenshots

<div align="center">
  <table>
    <tr>
      <td>
        <img src="./screenshots/main_task_view.png" alt="Main View" width="350px"/>
        <p align="center"><i>Main Task View</i></p>
      </td>
      <td>
        <img src="./screenshots/task_editing.png" alt="Edit View" width="350px"/>
        <p align="center"><i>Task Editing</i></p>
      </td>
    </tr>
    <tr>
      <td>
        <img src="./screenshots/help_screen.png" alt="Help View" width="350px"/>
        <p align="center"><i>Help Screen</i></p>
      </td>
      <td>
        <img src="./screenshots/task_filtering.png" alt="Categories View" width="350px"/>
        <p align="center"><i>Category Filtering</i></p>
      </td>
    </tr>
  </table>
</div>

## Installation

### Prerequisites

- Go 1.24+ (or use the pre-built binaries)
- A terminal with true color support recommended (iTerm2, Alacritty, Kitty, etc.)

### Via Homebrew (macOS)

```bash
# Install from the official tap
brew tap spmfte/tuiodo
brew install tuiodo
```

### Via AUR (Arch Linux)

> NOTE: Work in progress.

```bash
# Using yay
yay -S tuiodo

# Using paru
paru -S tuiodo

# Manual installation from AUR
git clone https://aur.archlinux.org/tuiodo.git
cd tuiodo
makepkg -si
```

### Via Go Install

```bash
go install github.com/spmfte/tuiodo@latest
```

### From Source

```bash
# Clone the repository
git clone https://github.com/spmfte/tuiodo.git
cd tuiodo

# Build the project
go build

# Run it
./tuiodo
```

## Quick Start

### Basic Usage

Just run `tuiodo` to start the application. Use the keyboard shortcuts to navigate and manage your tasks.
It is recommended to use the following alias:
```bash
echo "alias todo='tuiodo'" >> "$HOME/.$(which "$SHELL" | awk -F'/' '{print $NF}')rc"
```

#### What this command does:
1. `which "$SHELL"` retrieves the full path of the current shell (e.g., `/bin/zsh` or `/bin/bash`).
2. `awk -F'/' '{print $NF}'` extracts only the last part of the path, which is the shell name (e.g., `zsh` or `bash`).
3. `"$HOME/.$(which "$SHELL" | awk -F'/' '{print $NF}')rc"` constructs the appropriate shell configuration file path (`~/.zshrc`, `~/.bashrc`, etc.).
4. `echo "alias todo='tuiodo'" >> "$HOME/.$(which "$SHELL" | awk -F'/' '{print $NF}')rc"` appends the alias to the correct shell config file.
5. After running this command, reload your shell configuration with `source ~/.zshrc` (or `source ~/.bashrc`) to apply the alias.
6. You can now run `todo` instead of `tuiodo` to start the application.

### Creating Tasks

Press <kbd>a</kbd> to add a new task. You can organize tasks by category using the format:

"Category: Task description"

For example:

- "Work: Finish quarterly report"
- "Personal: Call mom about birthday plans"
- "Health: Schedule dentist appointment"

### Adding Priorities

Use the <kbd>p</kbd> key to cycle through priority levels for the selected task:

- `None → Low → Medium → High → Critical → None`

### Adding Metadata

Add metadata to tasks using @ notation:
- `@due:2023-12-31` - Sets a due date
- `@tag:important` - Adds a custom tag
- `@status:in-progress` - Sets a custom status

### Command-line Options

```bash
# Basic usage
tuiodo

# Use a specific config file
tuiodo --config ~/.config/tuiodo/my-config.yaml

# Use a different storage file
tuiodo --storage ~/projects/work-tasks.md

# Start with specific view and sorting
tuiodo --view pending --sort priority

# Terminal-friendly mode
tuiodo --no-mouse --no-color

# Configure backup behavior
tuiodo --backup-dir ~/backups --max-backups 10

# Start with specific category and view
tuiodo --category Work --view pending
```

## ⌨️ Keyboard Controls

| Action              | Keys                                   |
| ------------------- | -------------------------------------- |
| **Navigation**      |                                        |
| Move cursor down    | <kbd>j</kbd> <kbd>↓</kbd>              |
| Move cursor up      | <kbd>k</kbd> <kbd>↑</kbd>              |
| Next page           | <kbd>n</kbd> <kbd>→</kbd> <kbd>l</kbd> |
| Previous page       | <kbd>b</kbd> <kbd>←</kbd> <kbd>h</kbd> |
| Switch tabs         | <kbd>tab</kbd> <kbd>t</kbd>            |
| **Task Management** |                                        |
| Add task            | <kbd>a</kbd>                           |
| Edit task           | <kbd>e</kbd>                           |
| Delete task         | <kbd>d</kbd> (press twice to confirm)  |
| Undo delete         | <kbd>u</kbd>                           |
| Toggle completion   | <kbd>space</kbd> <kbd>enter</kbd>      |
| Expand task details | <kbd>x</kbd>                           |
| Cycle priority      | <kbd>p</kbd>                           |
| **Filtering**       |                                        |
| Cycle categories    | <kbd>c</kbd>                           |
| Sort by priority    | <kbd>s</kbd>                           |
| Sort by date        | <kbd>S</kbd>                           |
| Sort by category    | <kbd>C</kbd>                           |
| **Other**           |                                        |
| Show/hide help      | <kbd>?</kbd> <kbd>F1</kbd>             |
| Quit                | <kbd>q</kbd> <kbd>Ctrl+c</kbd>         |

## Configuration

TUIODO supports extensive configuration through a YAML file located at `~/.config/tuiodo/tuiodo.yaml`.

### Configuration Locations

Configuration files are automatically loaded from the following locations (in order):

1. Custom path specified with `--config /path/to/config.yaml`
2. `~/.config/tuiodo/tuiodo.yaml` (primary location)
3. User config directory as reported by OS
4. Current directory `./tuiodo.yaml`

### Creating a Default Config File

To generate a default configuration file:

```bash
tuiodo --create-default-config
```

### Configuration Sections

The configuration file is divided into these main sections:

#### 1. General Settings

```yaml
general:
  default_category: "Uncategorized" # Default category for new tasks
  show_status_bar: true # Show status bar at bottom
  tasks_per_page: 10 # Number of tasks to show per page
  clear_status_after_seconds: 3 # Time before status messages disappear
```

#### 2. UI Settings

```yaml
ui:
  show_header: true # Show app header
  header_format: "TUIODO" # Header text
  show_categories: true # Show category labels
  show_priorities: true # Show priority indicators
  task_separator: "─" # Character used to separate tasks
  enable_tabs: true # Show tab bar
  enable_borders: true # Show container borders
  border_style: "rounded" # Border style (rounded, normal, double, thick, none)
  date_format: "2006-01-02" # Go date format for creation dates
```

#### 3. Colors Settings

```yaml
colors:
  theme: "default"
  primary: "#7C3AED"
  secondary: "#2563EB"
  tertiary: "#10B981"
  # ... other base colors ...
  
  # Custom category colors
  category_colors:
    ui: "#8B5CF6"          # Purple for UI tasks
    add-task: "#EC4899"    # Pink for add-task
    bug: "#EF4444"         # Red for bugs
    function: "#10B981"    # Green for function
    fix: "#F59E0B"         # Amber for fixes
    functionality: "#3B82F6" # Blue for functionality
    layout: "#6366F1"      # Indigo for layout
    docs: "#2563EB"        # Blue for docs
    storage: "#14B8A6"     # Teal for storage
    # You can add your own categories:
    my-category: "#9333EA" # Custom color for your category
```

#### 4. Key Bindings

```yaml
keybindings:
  quit: ["q", "ctrl+c"]
  add_task: ["a"]
  edit_task: ["e"]
  delete_task: ["d"]
  toggle_task: ["enter", "space"]
  cycle_priority: ["p"]
  cycle_category: ["c"]
  cycle_tab: ["tab", "t"]
  next_page: ["right", "l", "n"]
  prev_page: ["left", "h", "b"]
  move_cursor_up: ["up", "k"]
  move_cursor_down: ["down", "j"]
  help: ["?", "F1"]
```

#### 5. Storage Settings

```yaml
storage:
  file_path: "TODO.md" # Path to task storage file
  backup_directory: "~/.config/tuiodo/backups" # Backup directory
  auto_save: true # Save automatically on changes
  backup_on_save: true # Create backups when saving
  max_backups: 5 # Maximum number of backups to keep
```

## 📝 Storage Format

Tasks are stored in a simple Markdown format that's human-readable and version-control friendly:

```markdown
## Work

- [ ] Prepare presentation @priority:high @due:2023-06-15
- [x] Send weekly report @priority:medium

## Personal

- [ ] Buy groceries @priority:low
- [ ] Call mom @due:2023-05-10
```

### Format Details

- **Categories**: Denoted by `## Category Name`
- **Tasks**: Uses GitHub-style checkbox syntax
  - `- [ ]` for pending tasks
  - `- [x]` for completed tasks
- **Metadata**:
  - Priorities: `@priority:high`, `@priority:medium`, `@priority:low`
  - Due dates: `@due:YYYY-MM-DD`

## Advanced Usage

### Custom Task Storage Location

You can store your tasks anywhere:

```bash
# Store in a Dropbox folder for sync between devices
tuiodo --storage ~/Dropbox/tasks.md

# Project-specific task list
tuiodo --storage ~/projects/awesome-project/TODO.md
```

### Theme Customization

Create a custom theme by defining your own colors in the config:

```yaml
colors:
  # Use a light theme
  theme: "custom"
  primary: "#8B5CF6" # Purple
  secondary: "#3B82F6" # Blue
  text: "#1F2937" # Dark gray
  background: "#F9FAFB" # Light gray
  # ... other color settings
```

### Task Dependencies (Coming Soon)

Link tasks together with dependencies:

```markdown
- [ ] Design project architecture @priority:high
- [ ] Implement database layer @depends:1
- [ ] Write unit tests @depends:2
```

## Upgrading

To upgrade to the latest version:

```bash
go install github.com/spmfte/tuiodo@latest
```

## Contributing

Contributions are very welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgements

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - The amazing TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions for terminal applications
- [Charm](https://charm.sh/) - For their suite of beautiful terminal tools
