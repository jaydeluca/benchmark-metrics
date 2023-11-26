package internal

import (
	"context"
	"github.com/google/go-github/v56/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

type TestGithubService struct {
	gitHubClient *github.Client
	owner        string
	repo         string
}

func (c *TestGithubService) GetMostRecentCommitSHA(ctx context.Context, timestamp time.Time, branch string) string {
	return "lkasjdflkajsd"
}

func (c *TestGithubService) GetFileContentsAtCommit(ctx context.Context, sha string) FileContents {
	// -----------------------------------------------------------------------------
	// Run at Mon Feb 14 05:17:37 UTC 2022
	// release : compares no agent, latest stable, and latest snapshot agents
	// 5 users, 5000 iterations
	// -----------------------------------------------------------------------------
	// Agent               :              none           latest         snapshot
	// Run duration        :          00:02:55         00:03:00         00:02:58
	// Avg. CPU (user) %   :        0.54903346       0.51652056          0.51358
	// Max. CPU (user) %   :        0.81407034        0.5959596        0.6111111
	// Avg. mch tot cpu %  :         0.9860928        0.9813821        0.9850744
	// Startup time (ms)   :             16806            18508            19178
	// Total allocated MB  :           9928.36         13951.66         13713.03
	// Min heap used (MB)  :             61.50            74.50            76.34
	// Max heap used (MB)  :            234.75           260.27           260.88
	// Thread switch rate  :           28917.9        33985.625        32923.316
	// GC time (ms)        :        2113474464       2714794711       2588153729
	// GC pause time (ms)  :              2248             2873             2735
	// Req. mean (ms)      :             12.59            12.75            12.72
	// Req. p95 (ms)       :             38.81            40.55            39.49
	// Iter. mean (ms)     :            170.88           175.07           174.24
	// Iter. p95 (ms)      :            276.19           263.54           260.11
	// Net read avg (bps)  :        4790281.00       2451560.00       2461257.00
	// Net write avg (bps) :        6908586.00      12662049.00      12728886.00
	// Peak threads        :                42               53               50
	// -----------------------------------------------------------------------------
	return FileContents{raw: "LS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLQogUnVuIGF0IE1vbiBGZWIgMTQgMDU6MTc6MzcgVVRDIDIwMjIKIHJlbGVhc2UgOiBjb21wYXJlcyBubyBhZ2VudCwgbGF0ZXN0IHN0YWJsZSwgYW5kIGxhdGVzdCBzbmFwc2hvdCBhZ2VudHMKIDUgdXNlcnMsIDUwMDAgaXRlcmF0aW9ucwotLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tCkFnZW50ICAgICAgICAgICAgICAgOiAgICAgICAgICAgICAgbm9uZSAgICAgICAgICAgbGF0ZXN0ICAgICAgICAgc25hcHNob3QKUnVuIGR1cmF0aW9uICAgICAgICA6ICAgICAgICAgIDAwOjAyOjU1ICAgICAgICAgMDA6MDM6MDAgICAgICAgICAwMDowMjo1OApBdmcuIENQVSAodXNlcikgJSAgIDogICAgICAgIDAuNTQ5MDMzNDYgICAgICAgMC41MTY1MjA1NiAgICAgICAgICAwLjUxMzU4Ck1heC4gQ1BVICh1c2VyKSAlICAgOiAgICAgICAgMC44MTQwNzAzNCAgICAgICAgMC41OTU5NTk2ICAgICAgICAwLjYxMTExMTEKQXZnLiBtY2ggdG90IGNwdSAlICA6ICAgICAgICAgMC45ODYwOTI4ICAgICAgICAwLjk4MTM4MjEgICAgICAgIDAuOTg1MDc0NApTdGFydHVwIHRpbWUgKG1zKSAgIDogICAgICAgICAgICAgMTY4MDYgICAgICAgICAgICAxODUwOCAgICAgICAgICAgIDE5MTc4ClRvdGFsIGFsbG9jYXRlZCBNQiAgOiAgICAgICAgICAgOTkyOC4zNiAgICAgICAgIDEzOTUxLjY2ICAgICAgICAgMTM3MTMuMDMKTWluIGhlYXAgdXNlZCAoTUIpICA6ICAgICAgICAgICAgIDYxLjUwICAgICAgICAgICAgNzQuNTAgICAgICAgICAgICA3Ni4zNApNYXggaGVhcCB1c2VkIChNQikgIDogICAgICAgICAgICAyMzQuNzUgICAgICAgICAgIDI2MC4yNyAgICAgICAgICAgMjYwLjg4ClRocmVhZCBzd2l0Y2ggcmF0ZSAgOiAgICAgICAgICAgMjg5MTcuOSAgICAgICAgMzM5ODUuNjI1ICAgICAgICAzMjkyMy4zMTYKR0MgdGltZSAobXMpICAgICAgICA6ICAgICAgICAyMTEzNDc0NDY0ICAgICAgIDI3MTQ3OTQ3MTEgICAgICAgMjU4ODE1MzcyOQpHQyBwYXVzZSB0aW1lIChtcykgIDogICAgICAgICAgICAgIDIyNDggICAgICAgICAgICAgMjg3MyAgICAgICAgICAgICAyNzM1ClJlcS4gbWVhbiAobXMpICAgICAgOiAgICAgICAgICAgICAxMi41OSAgICAgICAgICAgIDEyLjc1ICAgICAgICAgICAgMTIuNzIKUmVxLiBwOTUgKG1zKSAgICAgICA6ICAgICAgICAgICAgIDM4LjgxICAgICAgICAgICAgNDAuNTUgICAgICAgICAgICAzOS40OQpJdGVyLiBtZWFuIChtcykgICAgIDogICAgICAgICAgICAxNzAuODggICAgICAgICAgIDE3NS4wNyAgICAgICAgICAgMTc0LjI0Ckl0ZXIuIHA5NSAobXMpICAgICAgOiAgICAgICAgICAgIDI3Ni4xOSAgICAgICAgICAgMjYzLjU0ICAgICAgICAgICAyNjAuMTEKTmV0IHJlYWQgYXZnIChicHMpICA6ICAgICAgICA0NzkwMjgxLjAwICAgICAgIDI0NTE1NjAuMDAgICAgICAgMjQ2MTI1Ny4wMApOZXQgd3JpdGUgYXZnIChicHMpIDogICAgICAgIDY5MDg1ODYuMDAgICAgICAxMjY2MjA0OS4wMCAgICAgIDEyNzI4ODg2LjAwClBlYWsgdGhyZWFkcyAgICAgICAgOiAgICAgICAgICAgICAgICA0MiAgICAgICAgICAgICAgIDUzICAgICAgICAgICAgICAgNTAK"}
}

func TestFetchReportResultsInTimestampMappedReportData(t *testing.T) {
	err := os.Mkdir("cache", 0755)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll("cache")

	commitCache := NewSingleFileCache("cache/test-commit-cache.json")
	reportCache := NewSingleFileCache("cache/test-report-cache.json")
	mockedHttpClient := mock.NewMockedHTTPClient()

	ctx := context.Background()
	client := github.NewClient(mockedHttpClient)

	githubService := &TestGithubService{
		gitHubClient: client,
		owner:        "test",
		repo:         "test",
	}

	timeframe, _ := generateTimeframeSlice("2022-02-14", "2022-02-15", 1)
	benchmarkReport := BenchmarkReport{}
	benchmarkReport.FetchReports(ctx, timeframe, *commitCache, *reportCache, githubService)
	assert.Contains(t, benchmarkReport.ReportData["2022-02-14"], "Mon Feb 14 05:17:37 UTC 2022")

	err = commitCache.DeleteCache()
	if err != nil {
		return
	}
	err = reportCache.DeleteCache()
	if err != nil {
		return
	}
}

func TestGenerateReportResultsInPopulatedMetrics(t *testing.T) {
	err := os.Mkdir("cache", 0755)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll("cache")

	commitCache := NewSingleFileCache("cache/test-commit-cache2.json")
	reportCache := NewSingleFileCache("cache/test-report-cache2.json")
	mockedHttpClient := mock.NewMockedHTTPClient()

	ctx := context.Background()
	client := github.NewClient(mockedHttpClient)

	githubService := &TestGithubService{
		gitHubClient: client,
		owner:        "test",
		repo:         "test",
	}

	timeframe, _ := generateTimeframeSlice("2022-02-14", "2022-02-15", 1)
	benchmarkReport := BenchmarkReport{}
	benchmarkReport.FetchReports(ctx, timeframe, *commitCache, *reportCache, githubService)
	benchmarkReport.GenerateReport(timeframe)
	assert.Len(t, benchmarkReport.ResourceMetrics.ScopeMetrics[0].Metrics, 17)
	assert.Contains(t, benchmarkReport.MetricNames, "Max. CPU (user) %")

	err = commitCache.DeleteCache()
	if err != nil {
		return
	}
	err = reportCache.DeleteCache()
	if err != nil {
		return
	}
}
