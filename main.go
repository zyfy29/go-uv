package uv

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const InstallerScriptUrl = "https://astral.sh/uv/install.sh"

// getDefaultInstallPath returns the OS-appropriate default installation path
func getDefaultInstallPath() string {
	if runtime.GOOS == "windows" {
		tempDir := os.Getenv("TEMP")
		if tempDir == "" {
			tempDir = os.Getenv("TMP")
		}
		if tempDir == "" {
			tempDir = "C:\\temp"
		}
		return filepath.Join(tempDir, "go-uv")
	}
	return "/tmp/go-uv"
}

// getInstallerScriptUrl returns the appropriate installer script URL for the current OS
func getInstallerScriptUrl() string {
	if runtime.GOOS == "windows" {
		return "https://astral.sh/uv/install.ps1"
	}
	return InstallerScriptUrl
}

// getExecutableName returns the executable name with appropriate extension for the current OS
func getExecutableName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

var (
	installPath = getDefaultInstallPath()
)

// Init sets the installation path for uv and uvx binaries.
// if path is an empty string, it defaults to OS-appropriate path:
// Windows: %TEMP%\go-uv, macOS/Linux: /tmp/go-uv.
func Init(path string) {
	if path != "" {
		installPath = path
	}
}

// Install if uv, uvx are already installed, print a message and skip installation
// else, download the appropriate install script for the OS and execute it:
// Windows: https://astral.sh/uv/install.ps1 (PowerShell)
// macOS/Linux: https://astral.sh/uv/install.sh (Shell)
func Install(client *http.Client) error {
	fmt.Println("Installing uv to", installPath)
	if client == nil {
		client = http.DefaultClient
	}

	if _, err := os.Stat(filepath.Join(installPath, getExecutableName("uvx"))); err == nil {
		fmt.Println("uvx are already installed at", installPath)
		return nil
	}

	if err := os.MkdirAll(installPath, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	resp, err := client.Get(getInstallerScriptUrl())
	if err != nil {
		return fmt.Errorf("failed to download install script: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download install script: status %d", resp.StatusCode)
	}

	scriptContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read install script: %w", err)
	}

	var cmd *exec.Cmd
	env := append(os.Environ(), fmt.Sprintf("UV_UNMANAGED_INSTALL=%s", installPath))
	
	if runtime.GOOS == "windows" {
		// On Windows, use PowerShell to execute the .ps1 script
		cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", string(scriptContent))
	} else {
		// On Unix-like systems, use sh to execute the shell script
		cmd = exec.Command("sh", "-c", string(scriptContent))
	}
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install uv: %w", err)
	}

	fmt.Println("Successfully installed uv to", installPath)
	return nil
}

// Uv builds an *exec.Cmd to run uv with the given options.
func Uv(opt ...string) *exec.Cmd {
	return exec.Command(GetUvPath(), opt...)
}

// UvContext builds an *exec.Cmd to run uv with the given options and context.
func UvContext(ctx context.Context, opt ...string) *exec.Cmd {
	return exec.CommandContext(ctx, GetUvPath(), opt...)
}

// Uvx builds an *exec.Cmd to run uvx with the given options.
func Uvx(opt ...string) *exec.Cmd {
	return exec.Command(GetUvxPath(), opt...)
}

// UvxContext builds an *exec.Cmd to run uvx with the given options and context.
func UvxContext(ctx context.Context, opt ...string) *exec.Cmd {
	return exec.CommandContext(ctx, GetUvxPath(), opt...)
}

// GetUvPath returns the installation path of uv binaries.
func GetUvPath() string {
	return filepath.Join(installPath, getExecutableName("uv"))
}

// GetUvxPath returns the installation path of uvx binaries.
func GetUvxPath() string {
	return filepath.Join(installPath, getExecutableName("uvx"))
}
