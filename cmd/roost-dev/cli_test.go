package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckHelpFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{"no args", []string{}, false},
		{"no help flag", []string{"foo", "bar"}, false},
		{"short help", []string{"-h"}, true},
		{"long help", []string{"--help"}, true},
		{"help word", []string{"help"}, true},
		{"help in middle", []string{"foo", "-h", "bar"}, true},
		{"similar but not help", []string{"-help", "helper"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkHelpFlag(tt.args, "test usage")
			if result != tt.expected {
				t.Errorf("checkHelpFlag(%v) = %v, want %v", tt.args, result, tt.expected)
			}
		})
	}
}

func TestGetDefaultConfigDir(t *testing.T) {
	dir := getDefaultConfigDir()

	homeDir, _ := os.UserHomeDir()
	expected := filepath.Join(homeDir, ".config", "roost-dev")

	if dir != expected {
		t.Errorf("getDefaultConfigDir() = %q, want %q", dir, expected)
	}
}

func TestGetConfigWithDefaults(t *testing.T) {
	cfg, configDir := getConfigWithDefaults()

	if cfg == nil {
		t.Fatal("getConfigWithDefaults() returned nil config")
	}

	if cfg.TLD == "" {
		t.Error("expected non-empty TLD")
	}

	if configDir == "" {
		t.Error("expected non-empty configDir")
	}
}

func TestIsRoot(t *testing.T) {
	// Just verify it doesn't panic and returns a boolean
	result := isRoot()
	// We can't easily test the actual value since it depends on how tests are run
	_ = result
}

func TestRequireNonRoot(t *testing.T) {
	err := requireNonRoot("test action")

	if isRoot() {
		if err == nil {
			t.Error("expected error when running as root")
		}
	} else {
		if err != nil {
			t.Errorf("unexpected error when not root: %v", err)
		}
	}
}
