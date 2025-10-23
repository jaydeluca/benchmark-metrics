package internal

import (
	"fmt"
	"strings"
)

// PanelMetadata holds display information for dashboard panels
type PanelMetadata struct {
	Description string
	Unit        string
	Title       string
}

// getMetricMetadata returns display metadata for a given metric name
func getMetricMetadata(metricName string) PanelMetadata {
	metadataMap := map[string]PanelMetadata{
		"Startup time (ms)":       {Description: "Application startup time in milliseconds", Unit: "ms", Title: "Startup time"},
		"Run duration":            {Description: "Total benchmark run duration", Unit: "s", Title: "Run duration"},
		"Iter. mean (ms)":         {Description: "Mean iteration time in milliseconds", Unit: "ms", Title: "Iteration mean latency"},
		"Iter. p95 (ms)":          {Description: "95th percentile iteration time", Unit: "ms", Title: "Iteration p95 latency"},
		"Req. mean (ms)":          {Description: "Mean request processing time", Unit: "ms", Title: "Request mean latency"},
		"Req. p95 (ms)":           {Description: "95th percentile request processing time", Unit: "ms", Title: "Request p95 latency"},
		"Max heap used (MB)":      {Description: "Maximum heap memory used in megabytes", Unit: "mbytes", Title: "Max heap used"},
		"Min heap used (MB)":      {Description: "Minimum heap memory used", Unit: "mbytes", Title: "Min heap used"},
		"Total allocated MB":      {Description: "Total memory allocated over benchmark run", Unit: "mbytes", Title: "Total allocated memory"},
		"GC time (ms)":            {Description: "Total garbage collection time", Unit: "ms", Title: "GC time"},
		"GC pause time (ms)":      {Description: "Garbage collection pause time", Unit: "ms", Title: "GC pause time"},
		"Avg. CPU (user) %":       {Description: "Average CPU user utilization", Unit: "percent", Title: "Average CPU user"},
		"Max. CPU (user) %":       {Description: "Maximum CPU user utilization observed", Unit: "percent", Title: "Max CPU user"},
		"Peak threads":            {Description: "Peak number of threads during benchmark", Unit: "short", Title: "Peak threads"},
		"Thread switch rate":      {Description: "Thread context switching rate", Unit: "short", Title: "Thread switch rate"},
		"Net read avg (bps)":      {Description: "Average network read throughput in bytes per second", Unit: "Bps", Title: "Network read throughput"},
		"Net write avg (bps)":     {Description: "Average network write throughput in bytes per second", Unit: "Bps", Title: "Network write throughput"},
		"Avg. mch tot cpu %":      {Description: "Average machine total CPU utilization", Unit: "percent", Title: "Average machine total CPU"},
	}

	// Return metadata if found, otherwise return generic metadata
	if metadata, exists := metadataMap[metricName]; exists {
		return metadata
	}
	return PanelMetadata{
		Description: fmt.Sprintf("Metric: %s", metricName),
		Unit:        "short",
		Title:       metricName,
	}
}

// PanelLayout defines the order and width for each panel
type PanelLayout struct {
	MetricName string
	Width      int
}

// getPanelLayout returns the ordered layout configuration from the reference PDF
func getPanelLayout() []PanelLayout {
	return []PanelLayout{
		// Page 1 - Row 1: Two 12-wide panels
		{MetricName: "Startup time (ms)", Width: 12},
		{MetricName: "Run duration", Width: 12},
		// Page 1 - Row 2: Three 8-wide panels
		{MetricName: "Iter. mean (ms)", Width: 8},
		{MetricName: "Iter. p95 (ms)", Width: 8},
		{MetricName: "Req. mean (ms)", Width: 8},
		// Page 2 - Row 1: Two 12-wide panels
		{MetricName: "Req. p95 (ms)", Width: 12},
		{MetricName: "Max heap used (MB)", Width: 12},
		// Page 2 - Row 2: Three 8-wide panels
		{MetricName: "Min heap used (MB)", Width: 8},
		{MetricName: "Total allocated MB", Width: 8},
		{MetricName: "GC time (ms)", Width: 8},
		// Page 3 - Row 1: Two 12-wide panels
		{MetricName: "GC pause time (ms)", Width: 12},
		{MetricName: "Avg. CPU (user) %", Width: 12},
		// Page 3 - Row 2: Three 8-wide panels
		{MetricName: "Max. CPU (user) %", Width: 8},
		{MetricName: "Peak threads", Width: 8},
		{MetricName: "Thread switch rate", Width: 8},
		// Page 4 - Row 1: Two 12-wide panels
		{MetricName: "Net read avg (bps)", Width: 12},
		{MetricName: "Net write avg (bps)", Width: 12},
	}
}

func generateDashboard(title string, metrics []string) string {
	// Create a map of metrics for quick lookup
	metricMap := make(map[string]bool)
	for _, metric := range metrics {
		metricMap[metric] = true
	}

	// Use the predefined layout order
	panelLayout := getPanelLayout()

	var panels []string
	var currentX = 0
	var currentY = 0
	panelHeight := 8
	panelID := 1

	for _, layout := range panelLayout {
		// Only include panels for metrics that exist in the data
		if !metricMap[layout.MetricName] {
			continue
		}

		// Check if panel fits on current row (24 units wide total)
		if currentX+layout.Width > 24 {
			// Move to next row
			currentX = 0
			currentY += panelHeight
		}

		panels = append(panels, generatePanel(layout.MetricName, panelID, panelHeight, layout.Width, currentX, currentY))

		currentX += layout.Width
		panelID++
	}

	return generateDashboardJson(title, strings.Join(panels, ","))
}

func generateDashboardJson(title, panels string) string {
	return fmt.Sprintf(`{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": {
          "type": "grafana",
          "uid": "-- Grafana --"
        },
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "description": "Java Agent benchmark performance metrics including startup time, iteration performance, memory usage, CPU utilization, and network I/O",
  "editable": true,
  "fiscalYearStartMonth": 0,
  "graphTooltip": 0,
  "id": null,
  "links": [],
  "panels": [
    %s
  ],
  "preload": false,
  "schemaVersion": 42,
  "tags": ["benchmarks", "opentelemetry", "performance"],
  "templating": {
    "list": [
      {
        "current": {
          "text": "All",
          "value": [
            "$__all"
          ]
        },
        "datasource": {
          "type": "grafana-clickhouse-datasource",
          "uid": "clickhouse"
        },
        "includeAll": true,
        "label": "Entity",
        "multi": true,
        "name": "entity",
        "options": [],
        "query": "SELECT DISTINCT Attributes['entity'] as entity FROM default.otel_metrics_sum WHERE MetricName = 'benchmark_startup_time_ms' ORDER BY entity",
        "refresh": 1,
        "sort": 1,
        "type": "query"
      }
    ]
  },
  "time": {
    "from": "now-6h",
    "to": "now"
  },
  "timepicker": {},
  "timezone": "browser",
  "title": "%v",
  "uid": "benchmark-metrics",
  "version": 1
}`, panels, title)
}

func generatePanel(metricName string, panelID, panelHeight, panelWidth, currentX, currentY int) string {
	metadata := getMetricMetadata(metricName)

	gridPos := fmt.Sprintf(`{
        "h": %d,
        "w": %d,
        "x": %d,
        "y": %d
    }`, panelHeight, panelWidth, currentX, currentY)

	// SQL query for ClickHouse - using time series format with proper grouping and aggregation
	// Note: using single-line format to avoid JSON escaping issues
	sqlQuery := fmt.Sprintf(`SELECT toStartOfInterval(StartTimeUnix, toIntervalSecond($__interval_s)) AS t, coalesce(Attributes['entity'], 'none') AS metric, argMax(toFloat64(Value), StartTimeUnix) AS value FROM default.otel_metrics_sum WHERE MetricName = '%s' AND $__timeFilter(StartTimeUnix) GROUP BY t, metric ORDER BY t`, metricName)

	return fmt.Sprintf(`{
      "datasource": {
        "type": "grafana-clickhouse-datasource",
        "uid": "clickhouse"
      },
      "description": "%s",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "barWidthFactor": 0.6,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "showValues": false,
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": 0
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "%s"
        },
        "overrides": [
          {
            "matcher": {
              "id": "byRegexp",
              "options": "value .*"
            },
            "properties": [
              {
                "id": "displayName",
                "value": "${__field.labels.metric}"
              }
            ]
          }
        ]
      },
      "gridPos": %s,
      "id": %d,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "table",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "hideZeros": false,
          "mode": "single",
          "sort": "none"
        }
      },
      "pluginVersion": "12.3.0",
      "targets": [
        {
          "datasource": {
            "type": "grafana-clickhouse-datasource",
            "uid": "clickhouse"
          },
          "editorType": "sql",
          "format": 0,
          "queryType": "timeseries",
          "rawSql": "%s",
          "refId": "A"
        }
      ],
      "title": "%s",
      "type": "timeseries"
    }`, metadata.Description, metadata.Unit, gridPos, panelID, sqlQuery, metadata.Title)
}
