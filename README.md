# go-uv

go-uv is a [uv](https://docs.astral.sh/uv/) manager for Go which allows you run python scripts, commands and modules by `uv` and `uvx` without any other installation. 

## Prerequisites

- Go 1.18+
- Windows, macOS or Linux 

## Installation

```bash
go get -u github.com/zyfy29/go-uv
```


## Usage

```go
import "github.com/zyfy29/go-uv"

installPath := "/path/to/install" // Windows: "C:\\path\\to\\install"
uv.Init(installPath) // will be OS-appropriate temp dir if left empty:
                     // Windows: %TEMP%\go-uv, macOS/Linux: /tmp/go-uv
defer os.RemoveAll(installPath) // cleanup if you only need it once

// uv.Install require an *http.Client, nil for http.DefaultClient
if err := uv.Install(nil); err != nil {
    // ...
}

var command *exec.Cmd

// build commands
// uv run echo "hello world"
command = uv.Uv("run", "echo", "hello world")
// uvx --from openai-whisper whisper --help
command = uv.Uvx("--from", "openai-whisper", "whisper", "--help")

// and if you want to pass a context:
command = uv.UvContext(ctx, "run", "echo", "hello world")
conmmand = uv.UvxContext(ctx, "--from", "openai-whisper", "whisper", "--help")

// run commands
if err := command.Run(); err != nil {
    // ...
}
```

## Description

go-uv actually did an unmanaged installation by running [uv official install script](https://astral.sh/uv/install.sh) on Unix-like systems or [PowerShell script](https://astral.sh/uv/install.ps1) on Windows. This is convenient for temporary usage, and only `uv` and `uvx` will be persist in your disk.

**Platform Support:**
- **Windows**: Uses PowerShell script (`install.ps1`) and `.exe` executables
- **macOS/Linux**: Uses shell script (`install.sh`) and standard executables

## Contributing

Issues and pull requests are welcomed.
