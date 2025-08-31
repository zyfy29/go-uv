# go-uv

go-uv is a [uv](https://docs.astral.sh/uv/) manager for Go which allows you run python scripts, commands and modules by `uv` and `uvx` without any other installation. 

## Prerequisites

- Go 1.18+
- macOS or Linux 

## Installation

```bash
go get -u github.com/zyfy29/go-uv
```


## Usage

```go
import "github.com/zyfy29/go-uv"

installPath := "/path/to/install"
uv.Init(installPath) // will be /tmp/go-uv if left empty
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
command = uv.UvxContext(ctx, "--from", "openai-whisper", "whisper", "--help")

// run commands
if err := command.Run(); err != nil {
    // ...
}
```

## Description

go-uv actually did an unmanaged installation by running [uv official install script](https://astral.sh/uv/install.sh). This is convenient for temporary usage, and only `uv` and `uvx` will be persist in your disk.

## Contributing

Issues and pull requests are welcomed.
