package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileEntry represents a file or directory in the file system
type FileEntry struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	IsDirectory bool      `json:"isDirectory"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"modTime"`
}

// OpenFileResponse represents the state after opening a file or directory
type CurrentFilesState struct {
	CurrentDir  *FileEntry `json:"currentDir,omitempty"`
	CurrentFile *FileEntry `json:"currentFile,omitempty"`
	FileInfo    *FileEntry `json:"fileInfo,omitempty"`
	ContentHash string     `json:"contentHash,omitempty"`
}

func (a *App) ListFiles(path string) ([]FileEntry, error) {
	dirPath := path
	if dirPath == "" {
		dirPath = a.currentDir
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var fileEntries []FileEntry
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Printf("Error getting info for %s: %v\n", entry.Name(), err)
			continue
		}

		// Only include directories or markdown files
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext != ".md" && ext != ".markdown" {
				continue
			}
		}

		fileEntries = append(fileEntries, FileEntry{
			Name:        entry.Name(),
			Path:        filepath.Join(dirPath, entry.Name()),
			IsDirectory: entry.IsDir(),
			Size:        info.Size(),
			ModTime:     info.ModTime(),
		})
	}
	return fileEntries, nil
}

// OpenFile opens a file or changes the current directory
func (a *App) OpenFile(path string) (CurrentFilesState, error) {

	info, err := os.Stat(path)
	if err != nil {
		return CurrentFilesState{}, fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	if info.IsDir() {
		a.currentDir = path
		// Update window title when directory changes
		a.UpdateWindowTitleWithCurrentDir()
	} else {
		a.currentFile = path
	}

	return a.GetCurrentFilesState(), nil
}

// GoUp navigates one directory up from the current directory
func (a *App) GoUp() (*CurrentFilesState, error) {
	if a.currentDir == "" {
		return nil, fmt.Errorf("no current directory set")
	}

	// Get the parent directory
	parentDir := filepath.Dir(a.currentDir)

	// Check if we're already at the root
	if parentDir == a.currentDir {
		return nil, fmt.Errorf("already at root directory")
	}

	// Update current directory to parent
	a.currentDir = parentDir
	a.UpdateWindowTitleWithCurrentDir()

	// Get parent directory info
	dirInfo, err := os.Stat(parentDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory info for %s: %w", parentDir, err)
	}

	state := &CurrentFilesState{
		CurrentDir: &FileEntry{
			Name:        filepath.Base(parentDir),
			Path:        parentDir,
			IsDirectory: true,
			Size:        dirInfo.Size(),
			ModTime:     dirInfo.ModTime(),
		},
	}

	// Preserve CurrentFile if it still exists
	if a.currentFile != "" {
		fileInfo, err := os.Stat(a.currentFile)
		if err == nil {
			state.CurrentFile = &FileEntry{
				Name:        filepath.Base(a.currentFile),
				Path:        a.currentFile,
				IsDirectory: false,
				Size:        fileInfo.Size(),
				ModTime:     fileInfo.ModTime(),
			}
		}
	}

	return state, nil
}

// DeleteFile deletes a file or directory
func (a *App) DeleteFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	if info.IsDir() {
		// If deleting current directory, move to parent
		if a.currentDir == path {
			a.currentDir = filepath.Dir(path)
		}
	} else {
		// If deleting current file, clear current file state
		if a.currentFile == path {
			a.currentFile = ""
			a.currentFileContent = ""
		}
	}

	err = os.RemoveAll(path) // RemoveAll handles both files and directories recursively
	if err != nil {
		return fmt.Errorf("failed to delete %s: %w", path, err)
	}
	return nil
}

// GetFileContent reads and returns the content of the specified file
func (a *App) GetFileContent(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("path %s is a directory, not a file", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return string(content), nil
}

// GetFileContentPreview reads and returns the first 4 lines of the specified file
// Each line is truncated to max 128 characters with "..." appended if longer
func (a *App) GetFileContentPreview(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("path %s is a directory, not a file", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	lineCount := 0
	maxChars := 128

	for scanner.Scan() && lineCount < 4 {
		line := scanner.Text()
		// Truncate line if longer than maxChars
		if len(line) > maxChars {
			line = line[:maxChars] + "..."
		}
		lines = append(lines, line)
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return strings.Join(lines, "\n"), nil
}

// GetContentHash calculates and returns the SHA-256 hash of the given content
func (a *App) GetContentHash(content string) string {
	// Use SHA-256 for fast hashing with low collision chance
	hash := sha256.Sum256([]byte(content))
	fmt.Println("Content to hash:", content)
	fmt.Println("Calculated content hash:", hex.EncodeToString(hash[:]))
	return hex.EncodeToString(hash[:])
}

// GetCurrentState returns the current directory, file, and its content,
// along with the file list of the current directory.
func (a *App) GetCurrentFilesState() CurrentFilesState {
	state := CurrentFilesState{}

	// Populate CurrentDir
	if a.currentDir != "" {
		dirInfo, err := os.Stat(a.currentDir)
		if err == nil {
			state.CurrentDir = &FileEntry{
				Name:        filepath.Base(a.currentDir),
				Path:        a.currentDir,
				IsDirectory: true,
				Size:        dirInfo.Size(),
				ModTime:     dirInfo.ModTime(),
			}
		}
	}

	// Populate CurrentFile
	if a.currentFile != "" {
		fileInfo, err := os.Stat(a.currentFile)
		if err == nil {
			state.CurrentFile = &FileEntry{
				Name:        filepath.Base(a.currentFile),
				Path:        a.currentFile,
				IsDirectory: false,
				Size:        fileInfo.Size(),
				ModTime:     fileInfo.ModTime(),
			}

			// Calculate content hash
			content, err := os.ReadFile(a.currentFile)
			if err == nil {
				state.ContentHash = a.GetContentHash(string(content))
			}
		}
	}

	return state
}

// CreateFile creates a new file with the given name in the current directory
func (a *App) CreateFile(name string) error {
	if name == "" {
		return fmt.Errorf("file name cannot be empty")
	}

	filePath := filepath.Join(a.currentDir, name)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file %s already exists", name)
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", name, err)
	}
	defer file.Close()

	return nil
}

// CreateDir creates a new directory with the given name in the current directory
func (a *App) CreateDir(name string) error {
	if name == "" {
		return fmt.Errorf("directory name cannot be empty")
	}

	dirPath := filepath.Join(a.currentDir, name)

	// Check if directory already exists
	if _, err := os.Stat(dirPath); err == nil {
		return fmt.Errorf("directory %s already exists", name)
	}

	// Create the directory
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", name, err)
	}

	return nil
}

// RenameFile renames a file or directory
func (a *App) RenameFile(oldPath string, newName string) error {
	if newName == "" {
		return fmt.Errorf("new name cannot be empty")
	}

	// Check if old path exists
	info, err := os.Stat(oldPath)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", oldPath, err)
	}

	// Build the new path (same directory, new name)
	dir := filepath.Dir(oldPath)
	newPath := filepath.Join(dir, newName)

	// Check if new path already exists
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("a file or directory named %s already exists", newName)
	}

	// Rename the file/directory
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("failed to rename %s to %s: %w", oldPath, newName, err)
	}

	// Update current file/dir state if necessary
	if info.IsDir() && a.currentDir == oldPath {
		a.currentDir = newPath
	} else if !info.IsDir() && a.currentFile == oldPath {
		a.currentFile = newPath
	}

	return nil
}

// SaveFile saves content to the specified file
func (a *App) SaveFile(path string, content string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	// Don't allow saving to directories
	if info.IsDir() {
		return fmt.Errorf("path %s is a directory, not a file", path)
	}

	// Write content to file
	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to save file %s: %w", path, err)
	}

	// Update current file content cache
	if a.currentFile == path {
		a.currentFileContent = content
	}

	return nil
}
