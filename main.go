package uv

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const InstallerScriptUrl = "https://astral.sh/uv/install.sh"

var (
	InstallPath = "/tmp/go-uv"
)

// Init sets the installation path for uv and uvx binaries.
// if path is an empty string, it defaults to /tmp/go-uv.
func Init(path string) {
	if path != "" {
		InstallPath = path
	}
}

// Install if uv, uvx are already installed, print a message and skip installation
// else, download https://astral.sh/uv/install.sh using http.Client and execute it
func Install(client *http.Client) error {
	fmt.Println("Installing uv to", InstallPath)
	if client == nil {
		client = http.DefaultClient
	}

	if _, err := os.Stat(filepath.Join(InstallPath, "uvx")); err == nil {
		fmt.Println("uvx are already installed at", InstallPath)
		return nil
	}

	if err := os.MkdirAll(InstallPath, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	resp, err := client.Get(InstallerScriptUrl)
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

	cmd := exec.Command("sh", "-c", string(scriptContent))
	cmd.Env = append(os.Environ(), fmt.Sprintf("UV_UNMANAGED_INSTALL=%s", InstallPath))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install uv: %w", err)
	}

	fmt.Println("Successfully installed uv to", InstallPath)
	return nil
}

func Uvx(opt ...string) *exec.Cmd {
	return exec.Command(filepath.Join(InstallPath, "uvx"), opt...)
}

func UvxContext(ctx context.Context, opt ...string) *exec.Cmd {
	return exec.CommandContext(ctx, filepath.Join(InstallPath, "uvx"), opt...)
}
