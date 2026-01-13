package editor

import (
	"os"
	"testing"
)

func TestGetDefault(t *testing.T) {
	// Save original env vars
	origEditor := os.Getenv("EDITOR")
	origVisual := os.Getenv("VISUAL")
	defer func() {
		_ = os.Setenv("EDITOR", origEditor)
		_ = os.Setenv("VISUAL", origVisual)
	}()

	tests := []struct {
		name     string
		editor   string
		visual   string
		expected string
	}{
		{
			name:     "EDITOR set",
			editor:   "vim",
			visual:   "",
			expected: "vim",
		},
		{
			name:     "VISUAL set, EDITOR empty",
			editor:   "",
			visual:   "code",
			expected: "code",
		},
		{
			name:     "both set, EDITOR takes precedence",
			editor:   "nvim",
			visual:   "code",
			expected: "nvim",
		},
		{
			name:     "neither set, fallback to vi",
			editor:   "",
			visual:   "",
			expected: "vi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("EDITOR", tt.editor)
			_ = os.Setenv("VISUAL", tt.visual)

			result := getDefault()
			if result != tt.expected {
				t.Errorf("getDefault() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestEdit_NoModification(t *testing.T) {
	// Skip if running in CI without terminal
	if os.Getenv("CI") != "" {
		t.Skip("skipping editor test in CI")
	}

	// Use 'true' command as editor (exits immediately without modifying)
	origEditor := os.Getenv("EDITOR")
	defer func() { _ = os.Setenv("EDITOR", origEditor) }()
	_ = os.Setenv("EDITOR", "true")

	content := "test content\n"
	result, err := Edit(content)
	if err != nil {
		t.Fatalf("Edit() error = %v", err)
	}

	if result.Modified {
		t.Errorf("Edit() Modified = true, want false (editor didn't modify file)")
	}
}

func TestEdit_WithModification(t *testing.T) {
	// Skip if running in CI without terminal
	if os.Getenv("CI") != "" {
		t.Skip("skipping editor test in CI")
	}

	// Create a script that modifies the file
	script, err := os.CreateTemp("", "editor-*.sh")
	if err != nil {
		t.Fatalf("failed to create temp script: %v", err)
	}
	scriptPath := script.Name()
	defer func() { _ = os.Remove(scriptPath) }()

	_, _ = script.WriteString("#!/bin/sh\necho modified > \"$1\"\n")
	_ = script.Close()
	_ = os.Chmod(scriptPath, 0755)

	origEditor := os.Getenv("EDITOR")
	defer func() { _ = os.Setenv("EDITOR", origEditor) }()
	_ = os.Setenv("EDITOR", scriptPath)

	content := "original content\n"
	result, err := Edit(content)
	if err != nil {
		t.Fatalf("Edit() error = %v", err)
	}

	if !result.Modified {
		t.Errorf("Edit() Modified = false, want true")
	}

	if string(result.Content) != "modified\n" {
		t.Errorf("Edit() Content = %q, want %q", string(result.Content), "modified\n")
	}
}

func TestEdit_EditorFailure(t *testing.T) {
	// Skip if running in CI
	if os.Getenv("CI") != "" {
		t.Skip("skipping editor test in CI")
	}

	origEditor := os.Getenv("EDITOR")
	defer func() { _ = os.Setenv("EDITOR", origEditor) }()
	_ = os.Setenv("EDITOR", "false") // 'false' command always exits with error

	_, err := Edit("content")
	if err == nil {
		t.Error("Edit() expected error when editor fails")
	}
}
