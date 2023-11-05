package main

import "fmt"

func FetchReports(timeframe []string, commitCache, reportCache SingleFileCache, client ReportSource, repo string) map[string]string {
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
	return results
}
