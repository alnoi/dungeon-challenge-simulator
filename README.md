# Dungeon Event Processor

Dungeon Event Processor is a Go 1.22 backend-style solution for a dungeon challenge event stream.

The project is built around an event-driven state machine: it reads dungeon configuration, parses incoming events, updates player state, emits output events, and builds a final report.

## Features

- Parses dungeon configuration from JSON.
- Parses sequential incoming events in `[HH:MM:SS] playerID eventID [extra]` format.
- Processes registered players through a state machine.
- Handles disqualification, impossible moves, death, healing, damage, floor clearing, boss timing, dungeon close time, and final report generation.
- Formats event logs and final report in the required task output format.
- Includes unit tests and an end-to-end scenario test.

## Architecture

- `cmd/app` - application entrypoint and minimal processing pipeline.
- `internal/config` - CLI flags, environment variables, and default runtime settings.
- `internal/domain` - core domain structs, enum-like types, and errors.
- `internal/parser` - config and event parsers.
- `internal/engine` - event-driven state machine, handlers, and validators.
- `internal/output` - event and report formatters.
- `internal/timeutil` - parsing and formatting `HH:MM:SS` timestamps.
- `internal/*_test.go` - unit tests for individual packages.
- `tests/e2e_test.go` - end-to-end test for the full pipeline.

## Requirements

- Go 1.22 or newer.

## Running

By default, the app expects `config.json` and `events` in the project root:

```bash
go run ./cmd/app
```

You can pass paths with flags:

```bash
go run ./cmd/app -config /path/to/config.json -events /path/to/events
```

Or with environment variables:

```bash
DUNGEON_CONFIG_PATH=/path/to/config.json DUNGEON_EVENTS_PATH=/path/to/events go run ./cmd/app
```

Priority order:

```text
flags > environment variables > defaults
```

## Sample Data

Create basic local `config.json` and `events` files:

```bash
make sample-data
```

Then run the app:

```bash
make run
```

`make sample-data` overwrites local `config.json` and `events`.

## Make Targets

```bash
make fmt          # format Go code with gofmt
make test         # run all tests
make sample-data  # create sample config.json and events
make run          # run the app with CONFIG and EVENTS paths
```

You can override paths for `make run`:

```bash
make run CONFIG=/path/to/config.json EVENTS=/path/to/events
```

## Tests

Run all tests:

```bash
go test ./...
```

The test suite includes:

- unit tests for config loading, parsers, formatters, and time utilities;
- behavior tests for the engine package;
- an e2e test that runs a multi-player dungeon scenario through config parsing, event parsing, engine processing, and output formatting.

## Project Structure

```text
.
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   ├── domain/
│   │   ├── event.go
│   │   ├── player.go
│   │   ├── dungeon.go
│   │   ├── report.go
│   │   └── errors.go
│   ├── parser/
│   │   ├── config_parser.go
│   │   ├── config_parser_test.go
│   │   ├── event_parser.go
│   │   └── event_parser_test.go
│   ├── engine/
│   │   ├── engine.go
│   │   ├── handlers.go
│   │   ├── validator.go
│   │   └── engine_test.go
│   ├── output/
│   │   ├── event_formatter.go
│   │   ├── event_formatter_test.go
│   │   ├── report_formatter.go
│   │   └── report_formatter_test.go
│   └── timeutil/
│       ├── time.go
│       └── time_test.go
├── e2e/
│   └── dungeon_flow_test.go
├── go.mod
├── go.sum
├── README.md
└── Makefile
```
