package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/mattn/go-redmine"
)

// UpdateRedmineIssueRequest specifies the params to get pull request
type UpdateRedmineIssueRequest struct {
	CommitID      string
	AccessToken   string
	RedmineAPIKey string
	StatusID      int
	AssignedToID  int
}

type issueRequestBody struct {
	StatusID     int `json:"status_id"`
	AssignedToID int `json:"assigned_to_id"`
}

type issueUpdateRequest struct {
	Issue *issueRequestBody `json:"issue"`
}

type issueUpdateResponse struct {
	Issue *redmine.Issue `json:"issue"`
}

// UpdateRedmineIssue get the pull request by commit id
func UpdateRedmineIssue(r *UpdateRedmineIssueRequest) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	client := InitGithubClient(ctx, r.AccessToken)

	owner, repo := "BonioTw", "Facil"

	commit, _, _ := client.Repositories.GetCommit(ctx, owner, repo, r.CommitID)

	defer cancel()

	sha := commit.GetSHA()

	prs, _, _ := client.Search.Issues(ctx, sha, &github.SearchOptions{})

	endpoint := "http://redmine.bonio.com.tw"
	redmineClient := redmine.NewClient(endpoint, r.RedmineAPIKey)
	regex, _ := regexp.Compile(`http:\/\/redmine\.bonio\.com\.tw\/issues\/(\w*)`)

	for _, pr := range prs.Issues {
		fmt.Println("Found related PR:", pr.GetURL())
		matches := regex.FindAllStringSubmatch(pr.GetBody(), -1)
		if matches != nil {
			done := make(chan bool)

			for _, match := range matches {
				issueID := match[1]

				go func() {
					if issueID != "" {

						issueIDInt, _ := strconv.Atoi(issueID)
						issue, _ := redmineClient.Issue(issueIDInt)

						if statusID := issue.Status.Id; statusID > 11 {
							fmt.Println("Couldn't update redmine issue due to fasle status id:", statusID)
							done <- true
							return
						}

						fmt.Println("Start updating redmine issue:", issueID)
						APIEndpoint := endpoint + "/issues/" + issueID + ".json?key=" + r.RedmineAPIKey
						str, _ := json.Marshal(&issueUpdateRequest{
							Issue: &issueRequestBody{
								StatusID:     r.StatusID,
								AssignedToID: r.AssignedToID,
							},
						})
						req, err := http.NewRequest("PUT", APIEndpoint, strings.NewReader(string(str)))
						req.Header.Set("Content-Type", "application/json")
						redmineClient.Do(req)

						defer req.Body.Close()

						if err == nil {
							fmt.Println("Updated redmine issue successfully with params:", string(str))
						}

						done <- true
					}
				}()

				<-done

			}

		}
	}
}
