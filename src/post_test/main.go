package main

import (
	"facil-devops-utils/src/utils"
	"flag"
)

func main() {
	var (
		PullRequestNumber int
		JunitPath         string
		CommitID          string
		Event             string
		AccessToken       string
	)

	flag.IntVar(&PullRequestNumber, "pr", 0, "PullRequest number")
	flag.StringVar(&JunitPath, "junit", "junit.xml", "Github repository name")
	flag.StringVar(&CommitID, "commit-id", "", "Commit id on Github")
	flag.StringVar(&Event, "event", "COMMENT", "Event type")
	flag.StringVar(&AccessToken, "token", "", "Access token of Github")
	flag.Parse()

	body, _ := utils.ParseJunitText(JunitPath)
	args := &utils.CreateReviewRequest{
		PullRequestNumber: PullRequestNumber,
		Body:              body,
		CommitID:          CommitID,
		Event:             Event,
		AccessToken:       AccessToken,
	}
	utils.CreateReview(args)
}
