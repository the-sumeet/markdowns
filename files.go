package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	CurrentDir  string     `json:"currentDir"`
	CurrentFile string     `json:"currentFile"`
	FileInfo    *FileEntry `json:"fileInfo,omitempty"`
}

// ListFiles lists files and folders in the given path
func (a *App) ListFiles(path string) ([]FileEntry, error) {
	fmt.Println("listing", path)
	dirPath := path
	if dirPath == "" {
		dirPath = a.currentDir
	}

	entries, err := os.ReadDir(dirPath)
	fmt.Println(entries)
	fmt.Println("===")
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
func (a *App) OpenFile(path string) (*CurrentFilesState, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	if info.IsDir() {
		a.currentDir = path
	} else {
		a.currentFile = path

	}

	return &CurrentFilesState{
		CurrentDir:  a.currentDir,
		CurrentFile: a.currentFile,
		FileInfo: &FileEntry{
			Name:        info.Name(),
			Path:        path,
			IsDirectory: false,
			Size:        info.Size(),
			ModTime:     info.ModTime(),
		},
	}, nil
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

// GetCurrentState returns the current directory, file, and its content,
// along with the file list of the current directory.
func (a *App) GetCurrentFilesState() CurrentFilesState {
	return CurrentFilesState{
		CurrentDir:  a.currentDir,
		CurrentFile: a.currentFile,
	}
}
