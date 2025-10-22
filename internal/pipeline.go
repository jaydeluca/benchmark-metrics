package internal

import (
	"context"
	"fmt"
	"github.com/google/go-github/v56/github"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"os"
	"time"
)

func Run() {
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
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	githubService := &GithubService{
		repo:         repo,
		owner:        owner,
		gitHubClient: client,
	}
	// Calculate start date: 15 months before today
	startDate := time.Now().AddDate(0, -15, 0)
	timeframe, _ := generateTimeframeToToday(startDate.Format("2006-01-02"), 2)

	benchmarkReport := BenchmarkReport{}
	benchmarkReport.FetchReports(ctx, timeframe, *commitCache, *reportCache, githubService)
	benchmarkReport.GenerateReport(timeframe)

	// export to collector
	fmt.Println("Exporting metrics")
	_ = exp.Export(ctx, &benchmarkReport.ResourceMetrics)

	// create grafana dashboard
	dashboard := generateDashboard("Benchmark Metrics", benchmarkReport.MetricNames)
	err = os.WriteFile("grafana/dashboards/instrumentation-benchmarks.json", []byte(dashboard), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Generated dashboard")
}
