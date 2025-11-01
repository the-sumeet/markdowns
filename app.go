package main

import (
	"context"
	"fmt"
	"net/url"
	"os" // Added for os.UserHomeDir()
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context

	currentDir         string
	currentFile        string
	currentFileContent string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Load config to get last opened directory and file
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		// Fallback to home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting user home directory: %v\n", err)
			a.currentDir = "/" // Fallback to root
		} else {
			a.currentDir = homeDir
		}
	} else {
		// Use last opened directory if it exists
		if config.LastOpenedDirectory != "" {
			// Verify directory still exists
			if _, err := os.Stat(config.LastOpenedDirectory); err == nil {
				a.currentDir = config.LastOpenedDirectory
			} else {
				// Directory no longer exists, fall back to home
				homeDir, err := os.UserHomeDir()
				if err != nil {
					a.currentDir = "/"
				} else {
					a.currentDir = homeDir
				}
			}
		} else {
			// No last opened directory, use home
			homeDir, err := os.UserHomeDir()
			if err != nil {
				a.currentDir = "/"
			} else {
				a.currentDir = homeDir
			}
		}

		// Use last opened file if it exists
		if config.LastOpenedFile != "" {
			// Verify file still exists
			if _, err := os.Stat(config.LastOpenedFile); err == nil {
				a.currentFile = config.LastOpenedFile
			}
		}
	}

	// Set initial window title
	a.UpdateWindowTitleWithCurrentDir()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SetWindowTitle sets the window title
func (a *App) SetWindowTitle(title string) {
	runtime.WindowSetTitle(a.ctx, title)
}

// UpdateWindowTitleWithCurrentDir updates the window title to show the current directory
func (a *App) UpdateWindowTitleWithCurrentDir() {
	if a.currentDir != "" {
		runtime.WindowSetTitle(a.ctx, fmt.Sprintf("Markdowns - %s", a.currentDir))
	} else {
		runtime.WindowSetTitle(a.ctx, "Markdowns")
	}
}

// PickImageFile opens a file picker dialog for selecting image files and returns the URL-encoded relative path
func (a *App) PickImageFile() (string, error) {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Image",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Images (*.png;*.jpg;*.jpeg;*.gif;*.webp;*.svg)",
				Pattern:     "*.png;*.jpg;*.jpeg;*.gif;*.webp;*.svg",
			},
		},
	})

	if err != nil {
		return "", err
	}

	if selection == "" {
		return "", nil
	}

	// Get relative path from current file's directory
	var relativePath string
	if a.currentFile != "" {
		// Get the directory of the current file
		currentFileDir := filepath.Dir(a.currentFile)
		// Get relative path from current file's directory to the selected image
		relativePath, err = filepath.Rel(currentFileDir, selection)
		if err != nil {
			// If we can't get relative path, fall back to absolute
			relativePath = selection
		}
	} else {
		// If no current file, use path relative to current directory
		relativePath, err = filepath.Rel(a.currentDir, selection)
		if err != nil {
			relativePath = selection
		}
	}

	// URL encode the path to handle spaces and special characters
	encodedPath := url.PathEscape(relativePath)
	return encodedPath, nil
}

// ResolveImagePath resolves a relative image path to an absolute path based on the current file
func (a *App) ResolveImagePath(relativePath string) (string, error) {
	var baseDir string

	if a.currentFile != "" {
		// Use the directory of the current file as base
		baseDir = filepath.Dir(a.currentFile)
	} else {
		// Use current directory as base
		baseDir = a.currentDir
	}

	// Decode URL-encoded path
	decodedPath, err := url.PathUnescape(relativePath)
	if err != nil {
		return "", fmt.Errorf("failed to decode path: %w", err)
	}

	// If path is already absolute, return as is
	if filepath.IsAbs(decodedPath) {
		return decodedPath, nil
	}

	// Resolve relative path to absolute
	absolutePath := filepath.Join(baseDir, decodedPath)
	absolutePath = filepath.Clean(absolutePath)

	// Verify file exists
	if _, err := os.Stat(absolutePath); err != nil {
		return "", fmt.Errorf("image file not found: %w", err)
	}

	return absolutePath, nil
}
