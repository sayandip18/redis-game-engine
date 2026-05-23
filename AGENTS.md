# Development guide for AI SWE Agent

## Project overview

This is a Golang app, where the server communicates with Redis running inside a docker container.

## Development Workflow

After modifying any Go file, ALWAYS run:

```bash
gofumpt -w .
golangci-lint run
go test ./...
```

Never leave the repository in a failing state.
