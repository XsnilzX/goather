# AGENTS.md

This repository is a small Go CLI/weather widget. The guidance below is tuned
to how the codebase is already written.

## Repo overview

- `main.go`: CLI entrypoint, orchestrates location lookup + weather fetch.
- `location.go`: IP geolocation providers + option plumbing.
- `weather.go`: Open-Meteo API client.
- `weathercache.go`: cache file read/write helpers.
- `weatherformat.go`: icons, localized descriptions, tooltip formatting.
- `makefile`: build helpers.

## Build, lint, test

- Go version: `go 1.25.1` (from `go.mod`).
- Build (release): `make release` or `make`.
- Build (debug): `make debug`.
- Clean: `make clean` (removes binary + cache file).

### Tests

- No tests currently exist in the repo.
- Add new tests in `*_test.go` files and use standard `go test`.

### Run all tests (when tests exist)

- `go test ./...`

### Run a single package test file (when tests exist)

- `go test ./... -run TestName`
- `go test ./path/to/pkg -run TestName`

### Run a single subtest (when tests exist)

- `go test ./... -run TestName/SubtestName`

### Lint/format

- Formatting: `gofmt -w .` (standard Go format).
- Vet: `go vet ./...` (not in Makefile but standard).
- No other linters are configured (no golangci-lint config found).

## Cursor/Copilot rules

- No `.cursor/rules`, `.cursorrules`, or `.github/copilot-instructions.md`
  were found. Use the guidance in this file.

## Code style guidelines

### General Go conventions

- Follow standard Go style and run `gofmt` on changes.
- Keep packages small and focused; everything is in `package main`.
- Prefer short, clear names and avoid abbreviations that are unclear.
- Use `context.Context` for network calls and timeouts.

### Imports

- Use Go import grouping as `gofmt` organizes it.
- Avoid unused imports; the project does not use blank imports.
- Standard library imports only (currently no third‑party deps).

### Formatting

- Use tabs and `gofmt` formatting (do not hand-align with spaces).
- Keep line lengths reasonable; prefer multi-line `fmt.Sprintf` for long
  format strings (see `main.go`).
- Avoid extra whitespace at end of lines.

### Types and structs

- Use explicit types for JSON structures with tags.
- Keep exported types and fields in Go style (PascalCase for exported).
- Internal fields can be unexported, but note JSON encoding only includes
  exported fields; in this repo, cache structs currently use unexported fields.
- Prefer concrete types over `interface{}` unless needed for JSON decoding.

### Naming conventions

- Prefer Go’s camelCase or PascalCase in new code.
- Existing code has some snake_case helpers (e.g. `init_cache`). Keep naming
  consistent inside a file unless you are refactoring the whole file.
- Use verbs for functions (`FetchWeather`, `GetLocation`, `load_cache`).
- Use descriptive option names like `WithOverallTimeout`.

### Error handling

- Return errors upward; avoid `panic` except for truly unrecoverable cases.
- When returning errors, wrap with context where useful (e.g. `fmt.Errorf`).
- For user-visible errors, follow the pattern in `main.go` of emitting a JSON
  response with an error tooltip.

### Logging/output

- The main output is JSON written to stdout (`WaybarOut`).
- Avoid noisy logging in production path; use logs only for debugging.
- Preserve the output format expected by Waybar (text + tooltip keys).

### Concurrency

- Location providers are executed concurrently; if adding providers, keep them
  non-blocking and honor context timeouts.
- Use buffered channels when fanning out; see `location.go` for pattern.
- Avoid data races when reading shared state; keep per-request state local.

### HTTP usage

- Use `http.NewRequestWithContext` and per-request timeouts.
- Check response status codes explicitly; return errors on non-200.
- Always close response bodies with `defer resp.Body.Close()`.
- Set a recognizable `User-Agent` header.

### JSON handling

- Use `encoding/json` with struct tags for API responses.
- Keep JSON field tags in snake_case to match API conventions.
- For writing JSON to stdout, use `json.NewEncoder(os.Stdout).Encode(...)`.

### Cache behavior

- Cache file is hardcoded as `/tmp/weather_cache.json` (see `weathercache.go`).
- `init_cache` creates an empty file if missing; `update_cache` overwrites.
- If you change cache schema, update read/write functions together.

### Localization / UI strings

- Weather descriptions are in German and icons are emoji-based.
- Keep tooltip markup consistent (`<span size='xx-large'>`, `<big>` tags).
- Use `formatHourlyForecast` to append hourly blocks; keep output stable.

### Configuration

- Runtime settings are currently constants or options; no env parsing yet.
- If adding new flags/env vars, keep defaults and document them in README.

### Testing expectations

- Favor table-driven tests for providers, formatting, and caching.
- Mock HTTP with `httptest` for API calls.
- Keep tests deterministic; avoid real network calls in unit tests.

### File/dir additions

- Place new Go files at repo root unless you introduce subpackages.
- If adding a subpackage, update imports and keep `package main` entrypoint.

## Quick development loop

- Edit code and run `go test ./...` (when tests exist).
- Run `make debug` to build a local binary with symbols.
- Run the binary to inspect JSON output.

## Notes for agentic changes

- Keep edits focused; preserve existing behavior unless explicitly asked.
- Avoid adding third‑party dependencies without a clear need.
- When in doubt about style, follow the conventions already in the file you
  are touching.
