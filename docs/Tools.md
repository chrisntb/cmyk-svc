# Tools

This component requires the tools documented here.

## Go

See `https://go.dev/dl/`.

After downloading and installing `Go`, make sure the toolchain is configured for command line use.

Optionally add the Go SDK to your `PATH` if you manually installed instead of using a package manager:

```shell
# Use .bashrc or .zshrc depending on your preferred shell
cat <<'EOFF' | tee -a ${HOME}/.zshrc
# GO_SDK_HOME=...
# export PATH=${PATH}:${GO_SDK_HOME}/bin
EOF
```

Add the custom tools to your `PATH`:

```shell
# Use .bashrc or .zshrc depending on your preferred shell
cat <<'EOFF' | tee -a ${HOME}/.zshrc
export PATH=${PATH}:${HOME}/go/bin
EOF
```

## Linting

Check code quality using 'revive', see `https://github.com/mgechev/revive`:

```shell
go install github.com/mgechev/revive@v1.13.0
```

## API Specification

Generate OpenAPI Specification using 'swaggo', see `https://github.com/swaggo/swag`:

```shell
go install github.com/swaggo/swag/cmd/swag@v1.16.6
```

## API Tests

Run API tests using 'Bruno', see `https://www.usebruno.com/`.

First install Node via `nvm`:

```shell
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash

# Restart your shell
exec bash -l
# OR
exec zsh -l

nvm install --lts
```

Next install the Bruno CLI:

```shell
npm install -g @usebruno/cli
```

## VS Code Config

Recommended VS Code settings are available when using the workpace `cmyk-svc.code-workspace`.

### Extensions

> Check extension specific documentation since there are often post-install instructions.

- Claude Code for VS Code
  - `https://marketplace.visualstudio.com/items?itemName=anthropic.claude-code`

- Go
  - `https://marketplace.visualstudio.com/items?itemName=golang.Go`

Documentation:

- Markdownlint
  - `https://marketplace.visualstudio.com/items?itemName=DavidAnson.vscode-markdownlint`
- Markdown Preview Mermaid Support
  - `https://marketplace.visualstudio.com/items?itemName=bierner.markdown-mermaid`
- Mermaid Markdown Syntax Highlighting
  - `https://marketplace.visualstudio.com/items?itemName=bpruitt-goddard.mermaid-markdown-syntax-highlighting`

### Settings

Recommended additional `user settings`:

```json
{
  //...
  "telemetry.telemetryLevel": "off",
  //...
  "window.nativeTabs": true,
  "window.openFoldersInNewWindow": "on",
  //...
}
```
