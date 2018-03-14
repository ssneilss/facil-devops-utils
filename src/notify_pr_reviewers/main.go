package main

import (
	"context"
	"encoding/json"
	"facil-devops-utils/config"
	"facil-devops-utils/src/utils"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
)

type LabelFlags []string

func (i *LabelFlags) String() string {
	return ""
}

func (i *LabelFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func sendMessageToSlack(id string, pr *github.PullRequest) {
	ch := make(chan bool)

	go func() {
		title := pr.GetTitle()
		url := pr.GetHTMLURL()
		api := slack.New(config.SlackAPIToken)
		_, _, channelID, err := api.OpenIMChannel(id)
		if err != nil {
			log.Panic(err)
		}
		message := fmt.Sprintf("%s 幫測幫測 %s !!!", title, url)
		api.PostMessage(channelID, message, slack.PostMessageParameters{})

		ch <- true
	}()

	<-ch
}

func main() {
	var (
		Owner  string
		Repo   string
		Labels LabelFlags
	)

	flag.StringVar(&Owner, "owner", "", "Github Owner")
	flag.StringVar(&Repo, "repo", "", "Github Repo")
	flag.Var(&Labels, "filter", "Github PR Filter Label")
	flag.Parse()

	var ctx context.Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pullRequestReviews := utils.ListPRReviewers(ctx, &utils.ListPRInput{
		AccessToken: config.AccessToken,
		Owner:       Owner,
		Repo:        Repo,
		Labels:      Labels,
	})

	buff, _ := ioutil.ReadFile("./github_slack_accounts.json")
	var mapping map[string]string
	json.Unmarshal(buff, &mapping)

	c := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	for _, pullRequestReview := range pullRequestReviews {
		for _, githubUser := range pullRequestReview.Users {
			githubID := strconv.FormatInt(githubUser.GetID(), 10)
			slackID := mapping[githubID]
			cacheKey := fmt.Sprintf("%s%d", slackID, pullRequestReview.PullRequest.GetID())

			exists := c.Exists(cacheKey).Val()
			if exists == 1 {
				log.Println(pullRequestReview.PullRequest.GetHTMLURL(), "not going to notify")
				return
			}

			log.Println("Sending Message To:", githubUser.GetLogin(), slackID)

			hoursFromNow := time.Now().Sub(pullRequestReview.Issue.GetUpdatedAt()).Hours()
			log.Println(pullRequestReview.PullRequest.GetHTMLURL(), "last updated at", hoursFromNow, "hours ago")

			if hoursFromNow > 12 {
				// Inform every 10 mins
				c.Set(cacheKey, "true", 10*time.Minute)
				sendMessageToSlack(slackID, pullRequestReview.PullRequest)
			} else {
				// Inform every 20 mins
				c.Set(cacheKey, "true", 20*time.Minute)
				sendMessageToSlack(slackID, pullRequestReview.PullRequest)
			}
		}
	}
}
