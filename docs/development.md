# Development Guide

## Requirements

- Go 1.25 or newer.
- PostgreSQL 17 locally if you want to run the API against a database without Docker.
- Docker is optional later for Redis, MinIO, or containerized infrastructure.

## Common Commands

```sh
make run
make test
make vet
make fmt
make build
```

If you are not using Docker, create PostgreSQL locally with pgAdmin 4 using:

`docs/pgadmin4-macos.md`

## Environment

Copy `.env.local.example` to `.env` for local development. The current code reads environment variables directly; it does not parse `.env` files by itself. Use your shell, a task runner, or a process manager to load `.env`.

See `docs/environments.md` for the difference between `local`, `development`, `staging`, and `production`.

For pgAdmin 4 setup on macOS without Docker, see `docs/pgadmin4-macos.md`.

## Learning Order

1. Read `cmd/api/main.go`, then `internal/app/app.go`.
2. Trace configuration through `internal/config`.
3. Trace HTTP server setup through `internal/platform/httpserver`.
4. Implement one module vertically before implementing every module horizontally.

A useful first vertical slice would be registration or login. Do not start by writing every repository interface for every feature. That creates abstractions before you understand the real shape of the code.
