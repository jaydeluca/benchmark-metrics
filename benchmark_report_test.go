package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestClient struct{}

func (c *TestClient) GetMostRecentCommit(repo, timestamp, branch string) (string, error) {
	return "lkasjdflkajsd", nil
}

func (c *TestClient) GetFileAtCommit(repository, filepath, commitSHA string) (string, error) {
	return "----------------------------------------------------------\n Run at Mon Feb 14 05:17:37 UTC 2022\n release : compares no agent, latest stable, and latest snapshot agents\n 5 users, 5000 iterations\n----------------------------------------------------------\nAgent               :              none           latest         snapshot\nRun duration        :          00:02:55         00:03:00         00:02:58\nAvg. CPU (user) %   :        0.54903346       0.51652056          0.51358\nMax. CPU (user) %   :        0.81407034        0.5959596        0.6111111\nAvg. mch tot cpu %  :         0.9860928        0.9813821        0.9850744\nStartup time (ms)   :             16806            18508            19178\nTotal allocated MB  :           9928.36         13951.66         13713.03\nMin heap used (MB)  :             61.50            74.50            76.34\nMax heap used (MB)  :            234.75           260.27           260.88\nThread switch rate  :           28917.9        33985.625        32923.316\nGC time (ms)        :        2113474464       2714794711       2588153729\nGC pause time (ms)  :              2248             2873             2735\nReq. mean (ms)      :             12.59            12.75            12.72\nReq. p95 (ms)       :             38.81            40.55            39.49\nIter. mean (ms)     :            170.88           175.07           174.24\nIter. p95 (ms)      :            276.19           263.54           260.11\nNet read avg (bps)  :        4790281.00       2451560.00       2461257.00\nNet write avg (bps) :        6908586.00      12662049.00      12728886.00\nPeak threads        :                42               53               50\n", nil
}

func TestFetchReportResultsInTimestampMappedReportData(t *testing.T) {
	commitCache := NewSingleFileCache("cache/test-commit-cache.json")
	reportCache := NewSingleFileCache("cache/test-report-cache.json")

	githubClient := &TestClient{}

	timeframe, _ := generateTimeframeSlice("2022-02-14", "2022-02-15", 1)

	benchmarkReport := BenchmarkReport{}
	benchmarkReport.FetchReports(timeframe, *commitCache, *reportCache, githubClient, "test-repo")
	assert.Contains(t, benchmarkReport.ReportData["2022-02-14"], "Mon Feb 14 05:17:37 UTC 2022")

	err := commitCache.DeleteCache()
	if err != nil {
		return
	}
	err = reportCache.DeleteCache()
	if err != nil {
		return
	}
}

func TestGenerateReportResultsInPopulatedMetrics(t *testing.T) {
	commitCache := NewSingleFileCache("cache/test-commit-cache2.json")
	reportCache := NewSingleFileCache("cache/test-report-cache2.json")

	githubClient := &TestClient{}

	timeframe, _ := generateTimeframeSlice("2022-02-14", "2022-02-15", 1)

	benchmarkReport := BenchmarkReport{}
	benchmarkReport.FetchReports(timeframe, *commitCache, *reportCache, githubClient, "test-repo")
	benchmarkReport.GenerateReport(timeframe)
	assert.Len(t, benchmarkReport.ResourceMetrics.ScopeMetrics[0].Metrics, 17)
	assert.Contains(t, benchmarkReport.MetricNames, "Max. CPU (user) %")

	err := commitCache.DeleteCache()
	if err != nil {
		return
	}
	err = reportCache.DeleteCache()
	if err != nil {
		return
	}
}
