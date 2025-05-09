package main

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/v72/github"
)

func orgs(client *github.Client, ctx context.Context, username string) {
	orgs, resp, err := client.Organizations.List(ctx, username, nil)
	if err != nil {
		log.Error("Failed to fetch user's organizations", "err", err)
	}
	if len(orgs) == 0 {
		fmt.Println(
			GreyStyle.Render("[")+
				RedStyle.Render("x")+
				GreyStyle.Render("]"),
			GreyStyle.Render("No organizations found\n"),
		)
	} else {
		fmt.Println(
			GreyStyle.Render("[")+
				GreenStyle.Render("+")+
				GreyStyle.Render("]"),
			GreyStyle.Render("Organizations:\n"),
		)
	}
	for _, org := range orgs {
		PrintInfo("Login", org.GetLogin())
		PrintInfo("ID", fmt.Sprintf("%d", org.GetID()))
		PrintInfo("Node ID", org.GetNodeID())
		PrintInfo("URL", org.GetURL())
		PrintInfo("Repos URL", org.GetReposURL())
		PrintInfo("Events URL", org.GetEventsURL())
		PrintInfo("Hooks URL", org.GetHooksURL())
		PrintInfo("Issues URL", org.GetIssuesURL())
		PrintInfo("Members URL", org.GetMembersURL())
		PrintInfo("Public Members URL", org.GetPublicMembersURL())
		PrintInfo("Avatar URL", org.GetAvatarURL())
		PrintInfo("Description", org.GetDescription())
		fmt.Println()
	}
	WaitForRateLimit(resp)
}
