package uv

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-uv-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	Init(tempDir)

	t.Run("Install", func(t *testing.T) {
		client := &http.Client{Timeout: 30 * time.Second}
		err = Install(client)
		if err != nil {
			t.Fatalf("Install failed: %v", err)
		}
		uvPath := filepath.Join(installPath, "uv")
		uvxPath := filepath.Join(installPath, "uvx")

		if _, err := os.Stat(uvPath); os.IsNotExist(err) {
			t.Errorf("uv binary not found at %s", uvPath)
		}
		if _, err := os.Stat(uvxPath); os.IsNotExist(err) {
			t.Errorf("uvx binary not found at %s", uvxPath)
		}
	})

	t.Run("uv", func(t *testing.T) {
		echoCmd := Uv("run", "echo", "Hello, uv!")
		output, err := echoCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("uv run echo failed: %v\nOutput: %s", err, string(output))
		}
		t.Logf("uv run echo output:\n%s", string(output))
	})

	t.Run("uv with context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		sleepCmd := UvContext(ctx, "run", "sleep", "5")
		if err := sleepCmd.Run(); err != nil {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				t.Log("uv run sleep timed out as expected")
			} else {
				t.Errorf("unexpected error %v", err)
			}
		} else {
			t.Errorf("uv run sleep unexpectedly succeeded")
		}

	})

	t.Run("uvx", func(t *testing.T) {
		cowsayCmd := Uvx("cowsay", "-t", "Hello, uv!")
		output, err := cowsayCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("uvx cowsay failed: %v\nOutput: %s", err, string(output))
		}
		t.Logf("uvx cowsay output:\n%s", string(output))
	})
}
