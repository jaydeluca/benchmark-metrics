package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"os"
	"strings"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	repo := "open-telemetry/opentelemetry-java-instrumentation"
	filePath := "benchmark-overhead/results/release/summary.txt"

	client := NewGitHubClient(token)

	timeframe, _ := generateTimeframeToToday("2023-04-01", 30)

	dataPoints := map[string][]metricdata.DataPoint[float64]{}

	for _, timestamp := range timeframe {
		commit, _ := client.GetMostRecentCommit(repo, timestamp, "gh-pages")
		contents, _ := client.GetFileAtCommit(repo, filePath, commit)
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

	panels := []string{}

	var metrics []metricdata.Metrics
	for metric, metricData := range dataPoints {
		metrics = append(metrics, *generateMetrics(metric, metricData))
		panels = append(panels, generatePanel(metric, metric))
	}

	dashboard := generateDashboard(strings.Join(panels, ","))
	fmt.Print(dashboard)

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
	fmt.Sprintf("Exporting metrics")
	_ = exp.Export(ctx, resourceMetrics)

}
