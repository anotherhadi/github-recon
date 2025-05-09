package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/v72/github"
)

func userInfo(client *github.Client, ctx context.Context, username string) {
	fmt.Println(
		GreyStyle.Render("[")+
			GreenStyle.Render("-")+
			GreyStyle.Render("]"),
		GreyStyle.Render("Fetching user info...\n"),
	)

	user, resp, err := client.Users.Get(ctx, username)
	if resp.StatusCode == 404 {
		fmt.Println(
			GreyStyle.Render("[")+
				RedStyle.Render("x")+
				GreyStyle.Render("]"),
			GreyStyle.Render("Error:"),
			RedStyle.Render("User not found"),
		)
		os.Exit(1)
	}
	if err != nil {
		log.Error("Failed to fetch user's information", "err", err)
		os.Exit(1)
	}

	PrintInfo("Username", user.GetLogin())
	PrintInfo("ID", fmt.Sprintf("%d", user.GetID()))
	PrintInfo("Avatar URL", user.GetAvatarURL())
	PrintInfo("Gravatar ID", user.GetGravatarID())
	PrintInfo("Name", user.GetName())
	PrintInfo("Company", user.GetCompany())
	PrintInfo("Location", user.GetLocation())
	PrintInfo("Email", user.GetEmail())
	PrintInfo("Hireable", fmt.Sprintf("%t", user.GetHireable()))
	PrintInfo("Bio", user.GetBio())
	PrintInfo("Public Repos", fmt.Sprintf("%d", user.GetPublicRepos()))
	PrintInfo("Public Gists", fmt.Sprintf("%d", user.GetPublicGists()))
	PrintInfo("Followers", fmt.Sprintf("%d", user.GetFollowers()))
	PrintInfo("Following", fmt.Sprintf("%d", user.GetFollowing()))
	PrintInfo("Created At", user.GetCreatedAt().String())
	PrintInfo("Updated At", user.GetUpdatedAt().String())
	PrintInfo("Suspended At", user.GetSuspendedAt().String())
	PrintInfo("Type", user.GetType())
	PrintInfo("Site Admin", fmt.Sprintf("%t", user.GetSiteAdmin()))
	PrintInfo("Total Private Repos", fmt.Sprintf("%d", user.GetTotalPrivateRepos()))
	PrintInfo("Owned Private Repos", fmt.Sprintf("%d", user.GetOwnedPrivateRepos()))
	PrintInfo("Private Gists", fmt.Sprintf("%d", user.GetPrivateGists()))
	PrintInfo("Disk Usage", fmt.Sprintf("%d", user.GetDiskUsage()))
	PrintInfo("Collaborators", fmt.Sprintf("%d", user.GetCollaborators()))
	PrintInfo("Plan", user.GetPlan().GetName())
	fmt.Println()

	WaitForRateLimit(resp)
}
