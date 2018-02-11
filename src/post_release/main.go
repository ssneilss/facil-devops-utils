package main

import (
	"flag"

	"../utils"
)

func main() {
	var (
		CommitID      string
		StatusID      int
		AssignedToID  int
		AccessToken   string
		RedmineAPIKey string
	)

	flag.StringVar(&CommitID, "commit-id", "", "Commit id on Github")
	flag.IntVar(&StatusID, "status-id", 0, "Status id to update on Redmine")
	flag.IntVar(&AssignedToID, "assigned-to-id", 0, "Assignee id on Redmine")
	flag.StringVar(&AccessToken, "token", "", "Access token of Github")
	flag.StringVar(&RedmineAPIKey, "redmine-key", "", "Redmine API Key")
	flag.Parse()

	args := &utils.UpdateRedmineIssueRequest{
		CommitID:      CommitID,
		StatusID:      StatusID,
		AssignedToID:  AssignedToID,
		AccessToken:   AccessToken,
		RedmineAPIKey: RedmineAPIKey,
	}
	utils.UpdateRedmineIssue(args)
}
