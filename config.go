package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	LastOpenedFile      string            `json:"lastOpenedFile"`
	LastOpenedDirectory string            `json:"lastOpenedDirectory"`
	RecentFiles         []string          `json:"recentFiles"`
	WindowWidth         int               `json:"windowWidth"`
	WindowHeight        int               `json:"windowHeight"`
	Theme               string            `json:"theme"` // "light" or "dark"
	ShowHiddenFiles     bool              `json:"showHiddenFiles"`
	CustomSettings      map[string]string `json:"customSettings"`
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	return &Config{
		LastOpenedFile:      "",
		LastOpenedDirectory: "",
		RecentFiles:         []string{},
		WindowWidth:         1024,
		WindowHeight:        768,
		Theme:               "light",
		ShowHiddenFiles:     false,
		CustomSettings:      make(map[string]string),
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".markdowns")
	configPath := filepath.Join(configDir, "config.json")

	return configPath, nil
}

// LoadConfig loads the configuration from disk
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to disk
func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON with indentation for readability
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// App methods for config management

// GetConfig returns the current configuration
func (a *App) GetConfig() (*Config, error) {
	return LoadConfig()
}

// UpdateConfig updates the configuration
func (a *App) UpdateConfig(configJSON string) error {
	var config Config
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return SaveConfig(&config)
}

// UpdateConfigField updates a single field in the configuration
func (a *App) UpdateConfigField(field string, value string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	switch field {
	case "lastOpenedFile":
		config.LastOpenedFile = value
	case "lastOpenedDirectory":
		config.LastOpenedDirectory = value
	case "theme":
		config.Theme = value
	case "showHiddenFiles":
		config.ShowHiddenFiles = value == "true"
	default:
		// Store in custom settings if not a known field
		config.CustomSettings[field] = value
	}

	return SaveConfig(config)
}

// SetShowHiddenFiles sets the showHiddenFiles config option
func (a *App) SetShowHiddenFiles(show bool) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.ShowHiddenFiles = show
	return SaveConfig(config)
}

// GetShowHiddenFiles returns the showHiddenFiles config option
func (a *App) GetShowHiddenFiles() (bool, error) {
	config, err := LoadConfig()
	if err != nil {
		return false, err
	}

	return config.ShowHiddenFiles, nil
}

// AddRecentFile adds a file to the recent files list
func (a *App) AddRecentFile(filePath string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// Remove if already exists (to avoid duplicates)
	for i, f := range config.RecentFiles {
		if f == filePath {
			config.RecentFiles = append(config.RecentFiles[:i], config.RecentFiles[i+1:]...)
			break
		}
	}

	// Add to beginning of list
	config.RecentFiles = append([]string{filePath}, config.RecentFiles...)

	// Keep only last 10 files
	if len(config.RecentFiles) > 10 {
		config.RecentFiles = config.RecentFiles[:10]
	}

	return SaveConfig(config)
}

// GetRecentFiles returns the list of recent files
func (a *App) GetRecentFiles() ([]string, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return config.RecentFiles, nil
}
