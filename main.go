package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"os"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	repo := "open-telemetry/opentelemetry-java-instrumentation"
	filePath := "benchmark-overhead/results/release/summary.txt"

	// Cache API calls to github to prevent repeated calls when testing
	commitCache := NewSingleFileCache("cache/commit-cache.json")
	reportCache := NewSingleFileCache("cache/report-cache.json")

	client := NewGitHubClient(token)

	timeframe, _ := generateTimeframeToToday("2022-02-14", 7)

	dataPoints := map[string][]metricdata.DataPoint[float64]{}

	for _, timestamp := range timeframe {
		var commit string
		cached, _ := commitCache.RetrieveValue(timestamp)
		if cached == "" {
			commit, _ = client.GetMostRecentCommit(repo, timestamp, "gh-pages")
			err := commitCache.AddToCache(timestamp, commit)
			if err != nil {
				fmt.Println("Error adding to cache")
			}
		} else {
			commit = cached
		}

		var contents string
		cached, _ = reportCache.RetrieveValue(timestamp)
		if cached == "" {
			contents, _ = client.GetFileAtCommit(repo, filePath, commit)
			err := reportCache.AddToCache(timestamp, contents)
			if err != nil {
				fmt.Println("Error adding to cache")
			}
		} else {
			contents = cached
		}

		report := ParseReport(contents)
		for entity, metrics := range report.Metrics {

			for metricName, metricValue := range metrics {
				if _, ok := dataPoints[metricName]; !ok {
					// If the key doesn't exist, initialize a new slice
					dataPoints[metricName] = []metricdata.DataPoint[float64]{}
				}
				dataPoints[metricName] = append(dataPoints[metricName], *generateDataPoint(entity, report.Date, metricValue))

			}
		}
	}

	var metricNames []string

	var metrics []metricdata.Metrics
	for metric, metricData := range dataPoints {
		metrics = append(metrics, *generateMetrics(metric, metricData))
		metricNames = append(metricNames, metric)
	}

	resourceMetrics := generateResourceMetrics(metrics)

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

	// export to collector
	fmt.Println("Exporting metrics")
	_ = exp.Export(ctx, resourceMetrics)

	// Update Dashboard based on metrics
	generateDashboard(metricNames)
}
