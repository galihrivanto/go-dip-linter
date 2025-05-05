# go-dip-linter

`go-dip-linter` is a custom plugin for [golangci-lint](https://golangci-lint.run/) designed to enforce Dependency Inversion (D from SOLID) Principle

## Features

Detect violation of Dependency Inversion Principle by checking instatiation concrete type instead of using abstract (interface) 

## Installation

1. Install `golangci-lint` if you haven't already:

    ```bash
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    ```

2. Clone this repository:

    ```bash
    git clone https://github.com/galihrivanto/go-dip-linter.git
    cd go-dip-linter
    go build -o go-dip-linter .
    ```

3. Add the plugin to your `golangci-lint` configuration file (`.golangci.yml`):

    ```yaml
    version: "2"

linters:
  disable-all: true
  enable:
    - dip

  settings:
    custom:
      dip:
        type: module
        description: Detects violations of the Dependency Inversion Principle`.
        original-url: github.com/galihrivanto/go-dip-linter
        settings:
          paths:
            - name: /some/path
    ```

## Usage

Run `golangci-lint` with the custom plugin:

```bash
golangci-lint run
```