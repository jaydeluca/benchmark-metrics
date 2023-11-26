run:
	export OTEL_EXPORTER_OTLP_METRICS_ENDPOINT=http://localhost:4317; \
	go run main.go

test:
	go test ./... -cover -v
	golangci-lint run
