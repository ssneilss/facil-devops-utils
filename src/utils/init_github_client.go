package utils

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// InitGithubClient return the github client with oauth2 authentication
func InitGithubClient(ctx context.Context, accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}
