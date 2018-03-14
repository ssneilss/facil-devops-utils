package main

import (
	"testing"
)

// @TODO Get Slack users
// api := slack.New(config.SlackAPIToken)
// users, _ := api.GetUsers()
// for _, user := range users {
// 	fmt.Println(user.ID, user.Name)
// }

// @TODO Get Github user by account
// client := utils.InitGithubClient(context.TODO(), config.AccessToken)
// user, _, _ := client.Users.Get(context.TODO(), user)
// fmt.Println(user)

func TestNotifyPRUpdate(t *testing.T) {
	main()
}
