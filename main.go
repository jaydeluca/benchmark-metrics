package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v56/github"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"os"
)

func main() {
	repo := "opentelemetry-java-instrumentation"
	owner := "open-telemetry"

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

	// Cache API calls to github to prevent repeated calls when testing
	commitCache := NewSingleFileCache("cache/commit-cache.json")
	reportCache := NewSingleFileCache("cache/report-cache.json")
	client := &github.Client{}
	githubService := &GithubService{
		repo:         repo,
		owner:        owner,
		gitHubClient: client,
	}
	timeframe, _ := generateTimeframeToToday("2022-02-14", 7)

	benchmarkReport := BenchmarkReport{}
	benchmarkReport.FetchReports(ctx, timeframe, *commitCache, *reportCache, githubService)
	benchmarkReport.GenerateReport(timeframe)

	// export to collector
	fmt.Print("Exporting metrics")
	_ = exp.Export(ctx, &benchmarkReport.ResourceMetrics)

	// create grafana dashboard
	dashboard := generateDashboard("Benchmark Metrics", benchmarkReport.MetricNames)
	err = os.WriteFile("grafana/dashboards/instrumentation-benchmarks.json", []byte(dashboard), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Print("Generated dashboard")
}
