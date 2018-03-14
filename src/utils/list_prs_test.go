package utils

import (
	"context"
	"facil-devops-utils/config"
	"fmt"
	"testing"
)

func TestListPRs(t *testing.T) {
	t.Run("Should list prs correctly", func(t *testing.T) {
		prs := ListPRs(context.TODO(), &ListPRInput{
			AccessToken: config.AccessToken,
			Owner:       "BonioTw",
			Repo:        "Facil",
			CommitID:    "6f6e2999cf06093413da59c4359eaf58dcdecf90",
		})
		if len(prs) < 1 {
			t.Fail()
		}
	})

	t.Run("Should list all prs correctly", func(t *testing.T) {
		prs := ListPRReviewers(context.TODO(), &ListPRInput{
			AccessToken: config.AccessToken,
			Owner:       "BonioTw",
			Repo:        "Facil",
			Labels:      []string{"Tested"},
		})
		for _, pr := range prs {
			for _, user := range pr.Users {
				fmt.Println(pr.PullRequest.GetHTMLURL(), user.GetLogin())
			}
		}
		if len(prs) < 1 {
			t.Fail()
		}
	})
}
