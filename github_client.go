package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GitHubClient struct {
	BaseURL string
	token   string
	client  *http.Client
}

func NewGitHubClient(token string) *GitHubClient {
	client := &http.Client{}
	return &GitHubClient{
		BaseURL: "https://api.github.com",
		token:   token,
		client:  client,
	}
}

type Commit struct {
	SHA string `json:"sha"`
}

func (c *GitHubClient) GetMostRecentCommit(repo, timestamp, branch string) (string, error) {
	apiURL := fmt.Sprintf("%s/repos/%s/commits", c.BaseURL, repo)

	params := map[string]string{
		"per_page": "1",
		"until":    timestamp,
		"order":    "desc",
		"sha":      branch,
	}

	response, err := c.get(apiURL, params)
	if err != nil {
		return "", err
	}

	if response.StatusCode == http.StatusOK {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}

		var commits []Commit
		err = json.Unmarshal(body, &commits)
		if err != nil {
			return "", err
		}

		if len(commits) > 0 {
			return commits[0].SHA, nil
		} else {
			fmt.Println("No commits found.")
			return "", nil
		}
	} else {
		return "", fmt.Errorf("error: %d", response.StatusCode)
	}
}

func (c *GitHubClient) GetFileAtCommit(repository, filepath, commitSHA string) (string, error) {
	apiURL := fmt.Sprintf("%s/repos/%s/contents/%s", c.BaseURL, repository, filepath)

	params := map[string]string{
		"ref": commitSHA,
	}

	response, err := c.get(apiURL, params)
	if err != nil {
		return "", err
	}

	if response.StatusCode == http.StatusOK {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}

		var contentResponse struct {
			Content string `json:"content"`
		}

		err = json.Unmarshal(body, &contentResponse)
		if err != nil {
			return "", err
		}

		decodedContent, err := base64.StdEncoding.DecodeString(contentResponse.Content)
		if err != nil {
			return "", err
		}

		return string(decodedContent), nil
	} else {
		return "", fmt.Errorf("error: %d", response.StatusCode)
	}
}

func (c *GitHubClient) get(url string, params map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return c.client.Do(req)
}
