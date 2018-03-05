package utils

import (
	"context"
	"time"

	"github.com/google/go-github/github"
)

// CreateReviewRequest specifies the task definition used for CreateReview
type CreateReviewRequest struct {
	PullRequestNumber int
	Body              string
	Event             string
	CommitID          string
	AccessToken       string
}

func retry(attempts int, duration time.Duration, f func() error) error {
	if err := f(); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(duration)
			return retry(attempts, duration, f)
		}
	}
	return nil
}

// CreateReview create pull request comment on the target pull request
func CreateReview(r *CreateReviewRequest) (*github.PullRequestReview, error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	client := InitGithubClient(ctx, r.AccessToken)

	owner, repo, request :=
		"BonioTw",
		"Facil",
		&github.PullRequestReviewRequest{
			Body:     &r.Body,
			Event:    &r.Event,
			CommitID: &r.CommitID,
		}

	var (
		review *github.PullRequestReview
		err    error
	)

	defer cancel()

	retry(3, 2*time.Second, func() error {
		_review, _, _err := client.PullRequests.CreateReview(ctx, owner, repo, r.PullRequestNumber, request)
		review, err = _review, _err
		return err
	})

	return review, err
}
