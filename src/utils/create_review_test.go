package utils

import (
	"testing"
)

func TestCreateReview(t *testing.T) {
	t.Run("Should create pull request comment", func(t *testing.T) {
		CreateReview(&CreateReviewRequest{
			PullRequestNumber: 611,
			Body:              "123",
			Event:             "COMMENT",
			AccessToken:       "AccessToken",
		})
	})
}
