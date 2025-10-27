package main

import (
	"context"
	"fmt"
	"os" // Added for os.UserHomeDir()

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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		a.currentDir = "/" // Fallback to root
	} else {
		a.currentDir = homeDir
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
