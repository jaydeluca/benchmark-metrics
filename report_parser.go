package main

import (
	"math"
	"strconv"
	"strings"
	"time"
)

type ReportMetrics struct {
	Date    time.Time
	Metrics map[string]map[string]float64
}

func ParseReport(report string) ReportMetrics {

	split := strings.Split(report, "----------------------------------------------------------\n")
	date := strings.Split(strings.Split(split[1], "Run at ")[1], "\n")[0]

	metricsSplit := strings.Split(split[2], "\n")

	var metrics = map[string]map[string]float64{}

	// initialize maps for each config
	entities := splitByMultipleSpaces(strings.Split(metricsSplit[0], ":")[1])
	for _, entity := range entities {
		metrics[entity] = map[string]float64{}
	}

	for _, line := range metricsSplit {
		// skip header line and time
		if strings.HasPrefix(line, "Agent") || strings.HasPrefix(line, "Run duration") || line == "" {
			continue
		}
		metricList := strings.Split(line, ":")
		metricName := strings.TrimSpace(metricList[0])
		for index, value := range entities {
			thisMetric, _ := strconv.ParseFloat(splitByMultipleSpaces(metricList[1])[index], 32)
			metrics[value][metricName] = math.Round(thisMetric*100) / 100
		}
	}

	return ReportMetrics{
		Date:    convertDateFormat(date),
		Metrics: metrics,
	}
}
