package main

import (
	"fmt"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
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
	b.ResourceMetrics = *generateResourceMetrics(metrics)
	b.MetricNames = metricNames
}

func (b *BenchmarkReport) FetchReports(timeframe []string, commitCache, reportCache SingleFileCache, client ReportSource, repo string) {
	results := make(map[string]string)

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
			contents, _ = client.GetFileAtCommit(repo, "benchmark-overhead/results/release/summary.txt", commit)
			err := reportCache.AddToCache(timestamp, contents)
			if err != nil {
				fmt.Println("Error adding to cache")
			}
		} else {
			contents = cached
		}
		results[timestamp] = contents

	}
	b.ReportData = results
}
