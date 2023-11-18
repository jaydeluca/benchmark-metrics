package main

import (
	"context"
	"github.com/google/go-github/v56/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetMostRecentSHAReturnsSHA(t *testing.T) {
	mockedHttpClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposCommitsByOwnerByRepo,
			[]github.RepositoryCommit{
				{
					SHA: github.String("testSHA"),
				},
			},
		),
	)

	ctx := context.Background()
	client := github.NewClient(mockedHttpClient)

	githubService := GithubService{
		repo:         "opentelemetry-java-instrumentation",
		owner:        "opentelemetry",
		gitHubClient: client,
	}

	testTime, _ := time.Parse(layout, "2023-01-01")
	result := githubService.GetMostRecentCommitSHA(ctx, testTime, "gh-pages")
	assert.Equal(t, "testSHA", result)
}

func TestGetFileAtCommitReturnsExpectedDecodedContent(t *testing.T) {
	mockedHttpClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposContentsByOwnerByRepoByPath,
			github.RepositoryContent{
				Content: github.String("dGVzdEVuY29kZQ=="),
			},
		),
	)

	ctx := context.Background()
	client := github.NewClient(mockedHttpClient)

	githubService := GithubService{
		repo:         "opentelemetry-java-instrumentation",
		owner:        "opentelemetry",
		gitHubClient: client,
	}

	result := githubService.GetFileContentsAtCommit(ctx, "testSha")
	assert.Equal(t, "testEncode", result.ToString())
	assert.Equal(t, "dGVzdEVuY29kZQ==", result.raw)
}
