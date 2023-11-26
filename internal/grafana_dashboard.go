package internal

import (
	"fmt"
	"strings"
)

func generateDashboard(title string, metrics []string) string {
	var panels []string
	var currentX = 0
	var currentY = 0
	panelWidth := 8
	panelHeight := 8
	panelsPerRow := 3

	for _, metric := range metrics {
		panels = append(panels, generatePanel(metric, panelHeight, panelWidth, currentX, currentY))
		// Update the current position for the next panel
		currentX += panelWidth
		if currentX >= panelsPerRow*panelWidth {
			currentX = 0
			currentY += panelHeight
		}
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
  "editable": true,
  "fiscalYearStartMonth": 0,
  "graphTooltip": 0,
  "id": 2,
  "links": [],
  "liveNow": false,
  "panels": [
    %s
  ],
  "refresh": "",
  "schemaVersion": 38,
  "style": "dark",
  "tags": [],
  "templating": {
    "list": []
  },
  "time": {
    "from": "now/y",
    "to": "now/y"
  },
  "timepicker": {},
  "timezone": "",
  "title": "%v",
  "uid": "",
  "version": 2,
  "weekStart": ""
}`, panels, title)
}

func generatePanel(metricName string, panelHeight, panelWidth, currentX, currentY int) string {
	gridPos := fmt.Sprintf(`{
        "h": %d,
        "w": %d,
        "x": %d,
        "y": %d
    }`, panelHeight, panelWidth, currentX, currentY)

	return fmt.Sprintf(`{
      "datasource": {
        "type": "grafana-clickhouse-datasource",
        "uid": "P5C0FA5C61C0F8586"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
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
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": %s,
      "id": 1,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "grafana-clickhouse-datasource",
            "uid": "P5C0FA5C61C0F8586"
          },
          "format": 1,
          "meta": {
            "builderOptions": {
              "fields": [],
              "limit": 100,
              "mode": "list"
            }
          },
          "queryType": "sql",
          "rawSql": "SELECT\n    MetricName,\n    StartTimeUnix,\n    MAX(IF(Attributes['entity'] = 'none', Value, NULL)) AS none,\n    MAX(IF(Attributes['entity'] = 'snapshot', Value, NULL)) AS snapshot,\n    MAX(IF(Attributes['entity'] = 'latest', Value, NULL)) AS latest\nFROM otel.otel_metrics_sum\nWHERE MetricName = '%s'\nGROUP BY MetricName, StartTimeUnix\nORDER BY StartTimeUnix;",
          "refId": "A",
          "selectedFormat": 4
        }
      ],
      "title": "%s",
      "type": "timeseries"
    }`, gridPos, metricName, metricName)
}
