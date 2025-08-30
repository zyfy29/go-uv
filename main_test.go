package uv

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInstallAndUvx(t *testing.T) {
	// Setup: Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "go-uv-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set the install path to our temp directory
	Init(tempDir)

	// Test 1: Install uv
	client := &http.Client{Timeout: 30 * time.Second}
	err = Install(client)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Test 2: Verify uv and uvx binaries exist
	uvPath := filepath.Join(InstallPath, "uv")
	uvxPath := filepath.Join(InstallPath, "uvx")

	if _, err := os.Stat(uvPath); os.IsNotExist(err) {
		t.Errorf("uv binary not found at %s", uvPath)
	}
	if _, err := os.Stat(uvxPath); os.IsNotExist(err) {
		t.Errorf("uvx binary not found at %s", uvxPath)
	}

	// Test 3: Run uvx --help (lightest command)
	cmd := Uvx("--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("uvx --help failed: %v\nOutput: %s", err, string(output))
	}

	if len(output) == 0 {
		t.Error("uvx --help produced no output")
	}

	t.Logf("uvx --help output length: %d bytes", len(output))

	// Test 4: Run uvx --version (also very light)
	cmd = Uvx("--version")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("uvx --version failed: %v\nOutput: %s", err, string(output))
	}

	if len(output) == 0 {
		t.Error("uvx --version produced no output")
	}

	t.Logf("uvx --version output: %s", string(output))

	// Test 5: Test that Install is idempotent (doesn't reinstall if already installed)
	err = Install(client)
	if err != nil {
		t.Fatalf("Second Install call failed: %v", err)
	}

	// Test 6: Run a simple uvx command that requires a plugin (e.g., cowsay)
	cowsayCmd := Uvx("cowsay", "-t", "Hello, uv!")
	output, err = cowsayCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("uvx cowsay failed: %v\nOutput: %s", err, string(output))
	}
	t.Logf("uvx cowsay output:\n%s", string(output))

}
