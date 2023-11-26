package internal

import (
	"context"
	b64 "encoding/base64"
	"github.com/google/go-github/v56/github"
	"log"
	"time"
)

type RepoSource interface {
	GetMostRecentCommitSHA(ctx context.Context, timestamp time.Time, branch string) string
	GetFileContentsAtCommit(ctx context.Context, sha string) FileContents
}

type GithubService struct {
	gitHubClient *github.Client
	owner        string
	repo         string
}

type FileContents struct {
	raw string
}

func (fc *FileContents) ToString() string {
	decoded, _ := b64.URLEncoding.DecodeString(fc.raw)
	return string(decoded)
}

func (gs *GithubService) GetMostRecentCommitSHA(ctx context.Context, timestamp time.Time, branch string) string {
	opts := &github.CommitsListOptions{
		Until: timestamp,
		SHA:   branch,
	}
	result, _, err := gs.gitHubClient.Repositories.ListCommits(ctx, gs.owner, gs.repo, opts)
	if err != nil {
		log.Fatal(err)
	}
	return *result[0].SHA
}

func (gs *GithubService) GetFileContentsAtCommit(ctx context.Context, sha string) FileContents {
	path := "benchmark-overhead/results/release/summary.txt"
	opts := &github.RepositoryContentGetOptions{
		Ref: sha,
	}
	fileContent, _, _, err := gs.gitHubClient.Repositories.GetContents(ctx, gs.owner, gs.repo, path, opts)
	if err != nil {
		log.Fatal(err)
	}
	return FileContents{raw: *fileContent.Content}
}
