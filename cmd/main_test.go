package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCLIBasicUsage(t *testing.T) {
	// Build the CLI binary for testing
	cmd := exec.Command("go", "build", "-o", "godsays-test", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("godsays-test")

	// Test basic usage
	cmd = exec.Command("./godsays-test")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("CLI execution failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("Expected output, got empty string")
	}

	// Check that we get some words
	words := strings.Fields(string(output))
	if len(words) == 0 {
		t.Error("Expected words in output, got none")
	}
}

func TestCLIWithAmount(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "godsays-test", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("godsays-test")

	cmd = exec.Command("./godsays-test", "-amount", "5")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("CLI execution failed: %v", err)
	}

	outputStr := strings.TrimSpace(string(output))
	if len(outputStr) == 0 {
		t.Error("Expected non-empty output")
	}
}

func TestCLIInvalidAmount(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "godsays-test", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("godsays-test")

	testCases := []string{"0", "-1", "1001", "abc"}
	for _, amount := range testCases {
		cmd = exec.Command("./godsays-test", "-amount", amount)
		err := cmd.Run()
		if err == nil {
			t.Errorf("Expected error for amount %s, got none", amount)
		}
	}
}

func TestCLIHelp(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "godsays-test", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("godsays-test")

	cmd = exec.Command("./godsays-test", "-help")
	output, _ := cmd.CombinedOutput() // Get both stdout and stderr, ignore error as help exits with 0

	outputStr := string(output)
	if !strings.Contains(outputStr, "Usage:") {
		t.Errorf("Expected help output to contain 'Usage:', got: %s", outputStr)
	}
}
