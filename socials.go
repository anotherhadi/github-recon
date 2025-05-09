package main

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/v72/github"
)

func socials(client *github.Client, username string) {
	resp, err := FetchGitHubAPI(client, "", "/users/"+username+"/social_accounts")
	if err != nil {
		log.Error("Failed to fetch user's social media accounts", "err", err)
	}

	type SocialAccount struct {
		Provider string `json:"provider"`
		URL      string `json:"url"`
	}
	type SocialAccounts []SocialAccount

	var socialAccounts SocialAccounts
	err = json.Unmarshal(resp, &socialAccounts)
	if err != nil {
		log.Error("Failed to unmarshal social media accounts", "err", err)
	}
	if len(socialAccounts) == 0 {
		fmt.Println(
			GreyStyle.Render("[")+
				RedStyle.Render("x")+
				GreyStyle.Render("]"),
			GreyStyle.Render("No social media accounts found\n"),
		)
	} else {
		fmt.Println(
			GreyStyle.Render("[")+
				GreenStyle.Render("+")+
				GreyStyle.Render("]"),
			GreyStyle.Render("Social Media Accounts:\n"),
		)
	}
	for _, account := range socialAccounts {
		PrintInfo("Provider", account.Provider)
		PrintInfo("URL", account.URL)
		fmt.Println()
	}
}
