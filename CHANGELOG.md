# Changelog

All notable changes to TUIODO will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.1] - 2025-03-15

### Added
- Enhanced task expansion UI with improved formatting
- Advanced metadata tag support (@due, @tag, @status)
- Color-coded progress bar in status line
- Invisible attribution tag at the bottom of todo files

### Changed
- Updated cursor navigation to circular/wrap-around mode
- Improved task completion toggle (spacebar) functionality
- Enhanced delete confirmation with two-step process
- Improved task restoration through undo capability

### Fixed
- Fixed UI formatting issues in expanded task view
- Fixed space bar key handling for task completion
- Fixed task comparison for the Metadata map field

## [1.1.0] - 2025-03-14

### Added
- Comprehensive command-line interface with extensive flags
- Version information display (`--version`, `-v`)
- Monochrome mode support (`--no-color`)
- Dynamic category colors support through configuration
- Default storage path now set to `~/TODO.md`
- Comprehensive configuration system with YAML support
- Column alignment improvements in task list view
- Category-specific color styling in the UI
- Improved help screen with detailed keyboard shortcuts
- Debug mode with detailed logging (`--debug`)

### Changed
- Enhanced CLI flag handling with short and long forms
- Improved configuration validation and error handling
- Updated task sorting to include priority, creation date, and category
- Enhanced storage handling with better home directory support
- Simplified task metadata handling
- Improved UI responsiveness and layout

### Fixed
- Column alignment issues in the task list view
- Metadata visibility in task descriptions
- Storage path handling and expansion
- Category color application in the UI
- Configuration file handling and validation

## [1.0.0] - 2025-03-14

### Added
- Initial release of TUIODO
- Basic task management functionality
- Priority levels (Critical, High, Medium, Low)
- Category organization
- Keyboard-driven interface
- Mouse support
- Task filtering and sorting
- Markdown-based storage
- Automatic backups
- Configuration system 