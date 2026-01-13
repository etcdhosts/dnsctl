// Package editor provides system editor integration.
package editor

import (
	"fmt"
	"os"
	"os/exec"
)

// getDefault returns the default system editor.
// It checks $EDITOR, $VISUAL, then falls back to "vi".
func getDefault() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	return "vi"
}

// EditResult represents the result of an edit operation.
type EditResult struct {
	Content  []byte
	Modified bool
}

// Edit opens content in the system editor and returns the edited result.
// Returns Modified=false if user made no changes.
func Edit(content string) (*EditResult, error) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "dnsctl-*.hosts")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	if _, err := tmpFile.WriteString(content); err != nil {
		_ = tmpFile.Close()
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}

	// Get file info before editing
	infoBefore, err := os.Stat(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat temp file: %w", err)
	}

	// Open editor
	editor := getDefault()
	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("editor failed: %w", err)
	}

	// Check if file was modified
	infoAfter, err := os.Stat(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat temp file: %w", err)
	}

	if infoAfter.ModTime().Equal(infoBefore.ModTime()) {
		return &EditResult{Modified: false}, nil
	}

	// Read edited content
	edited, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read temp file: %w", err)
	}

	return &EditResult{
		Content:  edited,
		Modified: true,
	}, nil
}
