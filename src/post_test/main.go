package main

import (
	"os"
	"strconv"

	"../utils"
)

func main() {
	if len(os.Args) > 4 {
		pullRequestNumber, _ := strconv.Atoi(os.Args[1])
		body, _ := utils.ParseJunitText(os.Args[2])
		args := &utils.CreateReviewRequest{
			PullRequestNumber: pullRequestNumber,
			Body:              body,
			Event:             os.Args[3],
			AccessToken:       os.Args[4],
		}
		utils.CreateReview(args)
	}
}
