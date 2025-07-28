# Agent Guidelines for clong Repository

This document outlines the essential commands and code style guidelines for agents operating within this repository.

## Build/Lint/Test Commands

- **Build**: `go build -o bin/clong`
- **Lint**: `golangci-lint run`
- **Run all tests**: `go test -race -cover ./...`
- **Run a single test**: `go test -run <TestName> <package_path>` (e.g., `go test -run TestMyFunction ./internal/clong/service`)

## Code Style Guidelines

- **Formatting**: Adhere to `go fmt` standards. Run `go fmt ./...` to format code.
- **Static Analysis**: Use `golangci-lint` for static analysis.
- **Imports**: Organize imports using `goimports` (which `go fmt` typically handles).
- **Naming Conventions**: Follow Go's idiomatic naming conventions (e.g., `CamelCase` for exported names, `camelCase` for unexported names).
- **Error Handling**: Handle errors explicitly; do not ignore them. Return errors as the last return value.
- **Comments**: Add comments for complex logic or exported functions/types.
