package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	report = `----------------------------------------------------------
 Run at Fri Sep 01 05:16:59 UTC 2023
 release : compares no agent, latest stable, and latest snapshot agents
 5 users, 5000 iterations
----------------------------------------------------------
Agent               :              none           latest         snapshot
Run duration        :          00:01:54         00:02:11         00:02:11
Avg. CPU (user) %   :        0.48642316        0.5045967       0.51070267
Max. CPU (user) %   :             0.575       0.61650485        0.5920398
Avg. mch tot cpu %  :         0.9849136       0.98020655       0.98153454
Startup time (ms)   :             18360            13867            13802
Total allocated MB  :          32164.29         40026.94         40649.89
Min heap used (MB)  :             92.64           107.42           111.92
Max heap used (MB)  :            406.64           428.91           442.25
Thread switch rate  :         33919.312         34970.51         40582.11
GC time (ms)        :              1711             1484             2285
GC pause time (ms)  :              1727             1505             2285
Req. mean (ms)      :              8.35             9.65             9.67
Req. p95 (ms)       :             25.14            28.70            29.00
Iter. mean (ms)     :            111.87           129.19           129.52
Iter. p95 (ms)      :            186.74           207.58           208.03
Net read avg (bps)  :        7021778.00       6310764.00       6311767.00
Net write avg (bps) :        9356640.00      34143150.00      34205458.00
Peak threads        :                42               53               55`
)

func TestParseDateFromSummary(t *testing.T) {
	expected := "2023-09-01"
	result := ParseReport(report)
	assert.Equal(t, expected, result.Date.Format("2006-01-02"))
}

func TestParseRunDurationFromSummary(t *testing.T) {
	expected := 1.9 // 00:01:54 in HH:MM:SS format = 0*60 + 1 + 54/60 = 1.9 minutes
	result := ParseReport(report)
	assert.Equal(t, expected, result.Metrics["none"]["Run duration"])
}

func TestParseConfigsFromReport(t *testing.T) {
	expected := []string{"latest", "snapshot", "none"}
	result := ParseReport(report)

	var entities []string
	for entity := range result.Metrics {
		entities = append(entities, entity)
	}

	assert.ElementsMatch(t, expected, entities)
}

func TestParseMetricsFromReport(t *testing.T) {
	expected := 92.64
	result := ParseReport(report).Metrics["none"]["Min heap used (MB)"]
	assert.Equal(t, expected, result)

	expected = 53
	result = ParseReport(report).Metrics["latest"]["Peak threads"]
	assert.Equal(t, expected, result)
}
