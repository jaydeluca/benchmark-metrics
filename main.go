package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"os"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	repo := "open-telemetry/opentelemetry-java-instrumentation"

	ctx := context.Background()
	exp, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		panic(err)
	}

	meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(exp)))
	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	otel.SetMeterProvider(meterProvider)

	reports := []string{
		"release",
		"snapshot-regression",
		"snapshot",
	}

	// Cache API calls to github to prevent repeated calls when testing
	commitCache := NewSingleFileCache("cache/commit-cache.json")
	reportCache := NewSingleFileCache("cache/report-cache.json")

	client := NewGitHubClient(token)

	timeframe, _ := generateTimeframeToToday("2022-02-14", 7)

	data := FetchReports(timeframe, *commitCache, *reportCache, client, repo, reports)
	ConvertReports(reports, timeframe, data, exp, ctx)
}
