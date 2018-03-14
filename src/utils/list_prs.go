package utils

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

type ListPRInput struct {
	AccessToken string
	Owner       string
	Repo        string
	CommitID    string
	Labels      []string
}

type PRReviewersOutput struct {
	PullRequest *github.PullRequest
	Users       []*github.User
	Issue       github.Issue
}

func ListPRReviewers(ctx context.Context, l *ListPRInput) []*PRReviewersOutput {
	client := InitGithubClient(ctx, l.AccessToken)

	var output []*PRReviewersOutput

	searchQuery := fmt.Sprintf("is:pr is:open repo:%s/%s", l.Owner, l.Repo)
	for _, label := range l.Labels {
		searchQuery += fmt.Sprintf(" -label:%s", label)
	}

	searchResult, _, _ := client.Search.Issues(ctx, searchQuery, &github.SearchOptions{ListOptions: github.ListOptions{PerPage: 500}})

	isIssueExistsByURL := make(map[string]bool)
	issueByURL := make(map[string]github.Issue)
	for _, i := range searchResult.Issues {
		isIssueExistsByURL[i.GetURL()] = true
		issueByURL[i.GetURL()] = i
	}

	prs, _, _ := client.PullRequests.List(ctx, l.Owner, l.Repo, &github.PullRequestListOptions{})
	for _, pr := range prs {
		result, _, _ := client.PullRequests.ListReviewers(ctx, l.Owner, l.Repo, pr.GetNumber(), &github.ListOptions{PerPage: 500})
		if isIssueExistsByURL[pr.GetIssueURL()] != true {
			continue
		}
		output = append(output, &PRReviewersOutput{
			PullRequest: pr,
			Users:       result.Users,
			Issue:       issueByURL[pr.GetIssueURL()],
		})
	}
	return output
}

func ListPRs(ctx context.Context, l *ListPRInput) []github.Issue {
	client := InitGithubClient(ctx, l.AccessToken)
	commit, _, _ := client.Repositories.GetCommit(ctx, l.Owner, l.Repo, l.CommitID)
	sha := commit.GetSHA()
	result, _, _ := client.Search.Issues(ctx, sha, &github.SearchOptions{})
	return result.Issues
}
