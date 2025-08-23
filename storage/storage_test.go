package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindGitRepository(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "tuiodo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a nested directory structure
	nestedDir := filepath.Join(tempDir, "project", "src", "main")
	err = os.MkdirAll(nestedDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested dirs: %v", err)
	}

	// Create .git directory at the project level (not at root)
	gitDir := filepath.Join(tempDir, "project", ".git")
	err = os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	// Test finding git repository from nested directory
	foundGit, err := findGitRepository(nestedDir)
	if err != nil {
		t.Fatalf("Expected to find git repository, got error: %v", err)
	}

	expectedGit := filepath.Join(tempDir, "project")
	if foundGit != expectedGit {
		t.Errorf("Expected git root %s, got %s", expectedGit, foundGit)
	}

	// Test from a directory that's not in a git repository
	nonGitDir := filepath.Join(tempDir, "not-a-repo")
	err = os.Mkdir(nonGitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create non-git dir: %v", err)
	}

	_, err = findGitRepository(nonGitDir)
	if err == nil {
		t.Error("Expected error when not in git repository, got none")
	}
}

func TestGetGitRootTodoPath(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "tuiodo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a project directory with .git
	projectDir := filepath.Join(tempDir, "my-project")
	err = os.Mkdir(projectDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	gitDir := filepath.Join(projectDir, ".git")
	err = os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	// Create a TODO.md file at the project root
	todoPath := filepath.Join(projectDir, "TODO.md")
	todoContent := []byte("- [ ] Test task")
	err = os.WriteFile(todoPath, todoContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create TODO.md: %v", err)
	}

	// Change to the project directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	err = os.Chdir(projectDir)
	if err != nil {
		t.Fatalf("Failed to change to project directory: %v", err)
	}

	// Test getting the git root TODO path
	foundTodoPath, err := getGitRootTodoPath()
	if err != nil {
		t.Fatalf("Expected to find TODO.md in git repo, got error: %v", err)
	}

	// Resolve symlinks to handle macOS path differences
	resolvedTodoPath, err := filepath.EvalSymlinks(todoPath)
	if err != nil {
		resolvedTodoPath = todoPath
	}
	resolvedFoundPath, err := filepath.EvalSymlinks(foundTodoPath)
	if err != nil {
		resolvedFoundPath = foundTodoPath
	}

	if resolvedFoundPath != resolvedTodoPath {
		t.Errorf("Expected TODO path %s, got %s", resolvedTodoPath, resolvedFoundPath)
	}
}

func TestInitializeWithGitDetection(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "tuiodo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a project directory with .git
	projectDir := filepath.Join(tempDir, "my-project")
	err = os.Mkdir(projectDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	gitDir := filepath.Join(projectDir, ".git")
	err = os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	// Create a TODO.md file at the project root
	todoPath := filepath.Join(projectDir, "TODO.md")
	todoContent := []byte("- [ ] Test task")
	err = os.WriteFile(todoPath, todoContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create TODO.md: %v", err)
	}

	// Change to the project directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	err = os.Chdir(projectDir)
	if err != nil {
		t.Fatalf("Failed to change to project directory: %v", err)
	}

	// Test initialization with git detection
	Initialize("", "", 5, true, true)

	// Check if the storage path is set to the git repository TODO.md
	storagePath := GetStoragePath()

	// Resolve symlinks to handle macOS path differences
	resolvedTodoPath, err := filepath.EvalSymlinks(todoPath)
	if err != nil {
		resolvedTodoPath = todoPath
	}
	resolvedStoragePath, err := filepath.EvalSymlinks(storagePath)
	if err != nil {
		resolvedStoragePath = storagePath
	}

	if resolvedStoragePath != resolvedTodoPath {
		t.Errorf("Expected storage path to be git repo TODO.md %s, got %s", resolvedTodoPath, resolvedStoragePath)
	}

	// Test with explicit file path (should override git detection)
	explicitPath := filepath.Join(tempDir, "explicit-todo.md")
	Initialize(explicitPath, "", 5, true, true)

	storagePath = GetStoragePath()

	// Resolve symlinks for explicit path as well
	resolvedExplicitPath, err := filepath.EvalSymlinks(explicitPath)
	if err != nil {
		resolvedExplicitPath = explicitPath
	}
	resolvedStoragePath, err = filepath.EvalSymlinks(storagePath)
	if err != nil {
		resolvedStoragePath = storagePath
	}

	if resolvedStoragePath != resolvedExplicitPath {
		t.Errorf("Expected storage path to be explicit path %s, got %s", resolvedExplicitPath, resolvedStoragePath)
	}
}

func TestDebugGitDetection(t *testing.T) {
	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	t.Logf("Current working directory: %s", currentDir)

	// Check if we're in a git repository
	gitRoot, err := findGitRepository(currentDir)
	if err != nil {
		t.Logf("Not in git repository: %v", err)
	} else {
		t.Logf("Found git repository at: %s", gitRoot)

		// Check if TODO.md exists at git root
		todoPath := filepath.Join(gitRoot, "TODO.md")
		if _, err := os.Stat(todoPath); err == nil {
			t.Logf("TODO.md exists at git root: %s", todoPath)
		} else {
			t.Logf("TODO.md does not exist at git root: %s", todoPath)
		}
	}

	// Test the full path resolution
	if gitTodoPath, err := getGitRootTodoPath(); err == nil {
		t.Logf("Git root TODO path: %s", gitTodoPath)
	} else {
		t.Logf("Could not get git root TODO path: %v", err)
	}
}
