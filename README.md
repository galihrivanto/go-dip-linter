# go-dip-linter

`go-dip-linter` is a custom plugin for [golangci-lint](https://golangci-lint.run/) designed to enforce Dependency Inversion (D from SOLID) Principle

## Features

Detect violation of Dependency Inversion Principle by checking instatiation concrete type instead of using abstract (interface) 

## Installation

1. Install `golangci-lint` if you haven't already:

    ```bash
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    ```

2. See how to enable [custom linter](https://golangci-lint.run/plugins/module-plugins/)
sample `.custom-gcl.yml`
    ```yaml
    version: v1.64.8
    name: custom-golangci-lint
    plugins:
    - module: 'github.com/galihrivanto/go-dip-linter'
        path: /home/galih/codes/personal/go-dip-linter
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
              name_pattern:
                - service
    ```

## Usage

Run `golangci-lint` with the custom plugin:

```bash
golangci-lint run
```