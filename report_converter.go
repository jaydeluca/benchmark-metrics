package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func ConvertReports(reports, timeframe []string, data map[string]string, exp metric.Exporter, ctx context.Context) {
	for _, report := range reports {
		dataPoints := map[string][]metricdata.DataPoint[float64]{}
		for _, timestamp := range timeframe {

			report := ParseReport(data[fmt.Sprintf("%v-%v", timestamp, report)])
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

		resourceMetrics := generateResourceMetrics(metrics)

		// export to collector
		fmt.Printf("Exporting metrics for %v\n", report)
		_ = exp.Export(ctx, resourceMetrics)

		// Update Dashboard based on metrics
		generateDashboard(report, metricNames)
		fmt.Printf("Generated dashboard for %v", report)
	}
}
