package utils

import "testing"

func TestUpdateRedmineIssue(t *testing.T) {
	UpdateRedmineIssue(&UpdateRedmineIssueRequest{
		CommitID:      "CommitID",
		StatusID:      20,
		AssignedToID:  74,
		AccessToken:   "AccessToken",
		RedmineAPIKey: "RedmineAPIKey",
	})
}
