package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sea-monkeys/bob/config"
)

func ValidatePaths(config config.Config) error {
	// Check if prompt file exists
	if _, err := os.Stat(config.PromptPath); err != nil {
		return fmt.Errorf("prompt file not found: %v", err)
	}

	// Check if settings directory exists
	if info, err := os.Stat(config.SettingsPath); err != nil {
		return fmt.Errorf("settings directory not found: %v", err)
	} else if !info.IsDir() {
		return fmt.Errorf("settings path must be a directory")
	}

	// Check if output directory exists
	outputDir := filepath.Dir(config.OutputPath)
	if info, err := os.Stat(outputDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %v", err)
			}
		} else {
			return fmt.Errorf("error checking output directory: %v", err)
		}
	} else if !info.IsDir() {
		return fmt.Errorf("output path parent must be a directory")
	}

	return nil
}
