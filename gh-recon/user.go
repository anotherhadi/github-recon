package ghrecon

import (
	"fmt"
)

func (r Recon) User(username string) {
	user, resp, err := r.client.Users.Get(r.ctx, username)
	if resp.StatusCode == 404 {
		r.logger.Fatal("User not found")
	}
	if err != nil {
		r.logger.Fatal("Failed to fetch user's information", "err", err)
	}

	PrintTitle("👤 User informations")
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
