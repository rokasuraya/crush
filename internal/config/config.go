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
		MaxTokens: 16384, // bumped up: 8192 was often cutting off longer refactors
		// default system prompt: nudge the model toward concise, direct responses
		SystemPrompt: "Be concise and direct. Avoid unnecessary preamble or filler phrases.",
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
		return Default(), fmt.Errorf("reading config file: %w", err)
	}

	cfg := Default()
	if err := json.Unmarshal(data, cfg); err != nil {
		return Default(), fmt.Errorf("parsing config file: %w", err)
	}

	return cfg, nil
}

// Save writes the config to disk, creating parent directories as needed.
func (c *Config) Save() error {
	path, err := Path()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing config: %w", err)
	}

	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
