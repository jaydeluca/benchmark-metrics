package main

import "fmt"

func FetchReports(timeframe []string, commitCache, reportCache SingleFileCache, client ReportSource, repo string, reports []string) map[string]string {
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

		for _, report := range reports {
			mapKey := fmt.Sprintf("%v-%v", timestamp, report)
			filePath := fmt.Sprintf("benchmark-overhead/results/%v/summary.txt", report)
			cached, _ = reportCache.RetrieveValue(mapKey)
			if cached == "" {
				contents, _ = client.GetFileAtCommit(repo, filePath, commit)
				err := reportCache.AddToCache(mapKey, contents)
				if err != nil {
					fmt.Println("Error adding to cache")
				}
			} else {
				contents = cached
			}
			results[mapKey] = contents
		}

	}
	return results
}
