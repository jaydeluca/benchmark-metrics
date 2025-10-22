package internal

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
		// skip header line
		if strings.HasPrefix(line, "Agent") || line == "" {
			continue
		}

		// Bad data in some reports
		if strings.Contains(line, "8796093022208.00") {
			continue
		}

		metricList := strings.Split(line, ":")
		if len(metricList) < 2 {
			continue
		}

		metricName := strings.TrimSpace(metricList[0])

		for index, value := range entities {
			silo := splitByMultipleSpaces(metricList[1])
			if index >= len(silo) {
				continue
			}
			thisMetric, err := strconv.ParseFloat(silo[index], 32)
			if err == nil {
				metrics[value][metricName] = math.Round(thisMetric*100) / 100
			}
		}
	}

	return ReportMetrics{
		Date:    convertDateFormat(date),
		Metrics: metrics,
	}
}
