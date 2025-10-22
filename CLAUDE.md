# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This project fetches benchmark test reports from a GitHub repository (specifically OpenTelemetry Java Instrumentation's gh-pages branch), parses the performance data, and exports it as OpenTelemetry metrics to be visualized in Grafana. The primary use case is tracking historical benchmark performance over time.

## Build and Development Commands

### Running the application
```bash
# Run locally with OTLP collector endpoint
make run
# or manually:
export OTEL_EXPORTER_OTLP_METRICS_ENDPOINT=http://localhost:4317
go run main.go
```

### Testing
```bash
# Run all tests with coverage and linting
make test
# or manually:
go test ./... -cover -v
golangci-lint run
```

### Docker deployment
```bash
# Start all services (app, collector, clickhouse, grafana)
docker compose up -d
```

## Architecture

### Data Pipeline Flow (internal/pipeline.go)

The main execution flow in `Run()`:
1. **Fetch phase**: Query GitHub API for commits and file contents from the gh-pages branch
   - Uses `GithubService` to fetch commits at specific timestamps
   - Retrieves benchmark report file (`benchmark-overhead/results/release/summary.txt`)
   - Caches both commit SHAs and report contents to avoid rate limiting
2. **Parse phase**: Extract metrics from text-based reports
   - Parses report format with entity columns (e.g., "latest", "default")
   - Converts metrics to float64 values with rounding
3. **Convert phase**: Transform to OpenTelemetry metric format
   - Creates datapoints with timestamps and entity attributes
   - Builds ResourceMetrics structure
4. **Export phase**: Send to OTLP collector and generate Grafana dashboard JSON

### Key Components

**internal/benchmark_report.go** - Core orchestration
- `BenchmarkReport` struct holds metric names, OTLP data, and raw report text
- `FetchReports()`: Retrieves data from GitHub using caches
- `GenerateReport()`: Converts raw text to OTLP metrics

**internal/github_service.go** - GitHub API integration
- `RepoSource` interface allows mocking in tests
- `GithubService` wraps go-github client
- Methods: `GetMostRecentCommitSHA()`, `GetFileContentsAtCommit()`
- Hardcoded to fetch from `benchmark-overhead/results/release/summary.txt` path

**internal/report_parser.go** - Text parsing logic
- `ParseReport()`: Converts text format to structured `ReportMetrics`
- Handles multi-space delimited data with entity columns
- Filters out bad data (e.g., "8796093022208.00" anomalies)

**internal/metric_converter.go** - OTLP transformation
- Functions like `GenerateDataPoint()`, `GenerateMetrics()`, `GenerateResourceMetrics()`
- Converts parsed data to `metricdata.ResourceMetrics` structure

**internal/file_cache.go** - Local JSON file caching
- `SingleFileCache` stores key-value pairs in `cache/` directory
- Prevents repeated GitHub API calls during development/testing

**internal/grafana_dashboard.go** - Dashboard generation
- `generateDashboard()`: Creates Grafana JSON from metric names
- Outputs to `grafana/dashboards/instrumentation-benchmarks.json`

### CLI Structure (Cobra)

The application uses Cobra for command-line parsing:
- `cmd/root.go`: Defines flags for `--repo` and `--owner` (defaults to OpenTelemetry Java)
- `main.go`: Entry point that calls `cmd.Execute()` which invokes `internal.Run()`

### Configuration

- **GitHub authentication**: Set `GITHUB_TOKEN` env variable (optional but recommended to avoid rate limiting)
- **OTLP endpoint**: Set `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT` (default: http://localhost:4317)
- **Repository target**: Currently hardcoded in `internal/pipeline.go` to "open-telemetry/opentelemetry-java-instrumentation"
- **Timeframe**: Configured via `generateTimeframeToToday()` call in pipeline (starts from "2022-02-14", 2-day intervals)

### Testing

Tests use:
- `github.com/stretchr/testify` for assertions
- `github.com/migueleliasweb/go-github-mock` for mocking GitHub API calls
- Test files follow `*_test.go` convention alongside implementation files
