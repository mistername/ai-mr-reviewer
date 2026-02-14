package github

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/go-github/v68/github"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
)

type Client struct {
	client    *github.Client
	owner     string
	repo      string
	commitSHA string
	prNumber  int
}

func NewClient(token, owner, repo, prNumber, commitSHA string) (*Client, error) {
	httpClient := &http.Client{}
	client := github.NewClient(httpClient)
	client = client.WithAuthToken(token)

	prNum, err := strconv.Atoi(prNumber)
	if err != nil {
		return nil, fmt.Errorf("parse PR number: %w", err)
	}

	return &Client{
		client:    client,
		owner:     owner,
		repo:      repo,
		prNumber:  prNum,
		commitSHA: commitSHA,
	}, nil
}

func (c *Client) GetMergeRequestChanges() ([]domain.Diff, error) {
	files, _, err := c.client.PullRequests.ListFiles(context.Background(), c.owner, c.repo, c.prNumber, &github.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get PR files: %w", err)
	}

	var diffs []domain.Diff

	for _, f := range files {
		prevFilename := ""
		if f.PreviousFilename != nil {
			prevFilename = *f.PreviousFilename
		}

		diffs = append(diffs, domain.Diff{
			NewPath: *f.Filename,
			OldPath: prevFilename,
			Content: *f.Patch,
		})
	}

	return diffs, nil
}

func (c *Client) GetExistingComments() (map[string][]string, error) {
	existing := make(map[string][]string)

	reviewComments, _, err := c.client.PullRequests.ListComments(context.Background(), c.owner, c.repo, c.prNumber, &github.PullRequestListCommentsOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list PR review comments: %w", err)
	}

	for _, comment := range reviewComments {
		if comment.Path != nil && comment.Line != nil {
			key := fmt.Sprintf("%s:%d", *comment.Path, *comment.Line)
			existing[key] = append(existing[key], *comment.Body)
		}
	}

	issueComments, _, err := c.client.Issues.ListComments(context.Background(), c.owner, c.repo, c.prNumber, &github.IssueListCommentsOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list PR issue comments: %w", err)
	}

	for _, comment := range issueComments {
		if comment.Body != nil {
			existing["issue:"+*comment.Body] = append(existing["issue:"+*comment.Body], *comment.Body)
		}
	}

	return existing, nil
}

func (c *Client) AddMergeRequestDiscussion(file string, line int, note string) error {
	prComment := &github.PullRequestComment{
		Body:     &note,
		Path:     &file,
		Line:     github.Ptr(line),
		CommitID: &c.commitSHA,
		Side:     github.Ptr("RIGHT"),
	}

	_, _, err := c.client.PullRequests.CreateComment(context.Background(), c.owner, c.repo, c.prNumber, prComment)
	if err != nil {
		body := fmt.Sprintf("**File: %s**\n\n%s", file, note)
		issueComment := &github.IssueComment{
			Body: &body,
		}

		_, _, err = c.client.Issues.CreateComment(context.Background(), c.owner, c.repo, c.prNumber, issueComment)
		if err != nil {
			return fmt.Errorf("failed to add PR comment: %w", err)
		}
	}

	return nil
}

var _ domain.MRProviderPort = (*Client)(nil)
