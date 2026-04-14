// Package config handles loading, saving, and managing crush configuration.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	// AppName is the name of the application used for config directory resolution.
	AppName = "crush"
	// ConfigFileName is the default config file name.
	ConfigFileName = "config.json"
)

// Config holds the application configuration.
type Config struct {
	// APIKey is the API key used to authenticate with the AI provider.
	APIKey string `json:"api_key,omitempty"`
	// Model is the AI model to use for completions.
	Model string `json:"model,omitempty"`
	// Theme is the UI color theme.
	Theme string `json:"theme,omitempty"`
	// MaxTokens is the maximum number of tokens to generate.
	MaxTokens int `json:"max_tokens,omitempty"`
	// SystemPrompt is an optional system prompt to prepend to conversations.
	SystemPrompt string `json:"system_prompt,omitempty"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Model:     "claude-opus-4-5",
		Theme:     "dark", // personal preference: dark theme
		MaxTokens: 32768, // bumped up again: 16384 still cuts off on large codebases
		// default system prompt: nudge the model toward concise, direct responses
		// also ask for Go-idiomatic code since that's mostly what I use this for
		// added Python to the list since I've been doing more data work lately
		// added a note about error handling since I kept getting lazy try/except blocks
		// added TypeScript since I've been doing more frontend work recently
		SystemPrompt: "Be concise and direct. Avoid unnecessary preamble or filler phrases. When writing Go code, follow idiomatic Go style and conventions. When writing Python, follow PEP 8 and prefer stdlib over third-party packages where reasonable. When writing TypeScript, prefer explicit types over 'any' and use modern ES features. Always handle errors explicitly; do not swallow or silently ignore them.",
	}
}

// Dir returns the platform-appropriate configuration directory for crush.
func Dir() (string, error) {
	dir, err := xdg.ConfigFile(AppName)
	if err != nil {
		// Fallback to ~/.crush
		home, herr := os.UserHomeDir()
		if herr != nil {
			return "", fmt.Errorf("could not determine config directory: %w", err)
		}
		return filepath.Join(home, "."+AppName), nil
	}
	return filepath.Dir(dir), nil
}

// Path returns the full path to the config file.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ConfigFileName), nil
}

// Load reads and parses the config file from disk. If the file does not exist,
// a default config is returned without error.
func Load() (*Config, error) {
	path, err := Path()
	if err != nil {
		return Default(), err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Default(), nil
		}
		return Default(), fmt.Errorf("reading config file:
