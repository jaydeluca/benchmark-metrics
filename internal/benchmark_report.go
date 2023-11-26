package internal

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"time"
)

type BenchmarkReport struct {
	MetricNames     []string
	ResourceMetrics metricdata.ResourceMetrics
	ReportData      map[string]string
}

func (b *BenchmarkReport) GenerateReport(timeframe []string) {
	dataPoints := map[string][]metricdata.DataPoint[float64]{}
	for _, timestamp := range timeframe {
		report := ParseReport(b.ReportData[timestamp])
		for entity, metrics := range report.Metrics {
			for metricName, metricValue := range metrics {
				if _, ok := dataPoints[metricName]; !ok {
					dataPoints[metricName] = []metricdata.DataPoint[float64]{}
				}
				dataPoints[metricName] = append(dataPoints[metricName], *GenerateDataPoint(entity, report.Date, metricValue))
			}
		}
	}

	var metricNames []string
	var metrics []metricdata.Metrics
	for metric, metricData := range dataPoints {
		metrics = append(metrics, *GenerateMetrics(metric, metricData))
		metricNames = append(metricNames, metric)
	}
	b.ResourceMetrics = *GenerateResourceMetrics(metrics)
	b.MetricNames = metricNames
}

func (b *BenchmarkReport) FetchReports(ctx context.Context, timeframe []string, commitCache, reportCache SingleFileCache, githubService RepoSource) {
	results := make(map[string]string)

	for _, timestamp := range timeframe {
		var commitSHA string
		cached, _ := commitCache.RetrieveValue(timestamp)
		if cached == "" {
			convertedTimestamp, _ := time.Parse(layout, timestamp)
			commitSHA = githubService.GetMostRecentCommitSHA(ctx, convertedTimestamp, "gh-pages")
			err := commitCache.AddToCache(timestamp, commitSHA)
			if err != nil {
				fmt.Println("Error adding to cache")
			}
		} else {
			commitSHA = cached
		}

		var contents FileContents
		cached, _ = reportCache.RetrieveValue(timestamp)
		if cached == "" {
			contents = githubService.GetFileContentsAtCommit(ctx, commitSHA)
			err := reportCache.AddToCache(timestamp, contents.raw)
			if err != nil {
				fmt.Println("Error adding to cache")
			}
		} else {
			contents = FileContents{raw: cached}
		}
		results[timestamp] = contents.ToString()

	}
	b.ReportData = results
}
