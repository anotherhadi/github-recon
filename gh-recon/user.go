package ghrecon

import (
	"fmt"
)

type UserResult struct {
	Username          string
	ID                string
	AvatarURL         string
	GravatarID        string
	Name              string
	Company           string
	Location          string
	Email             string
	Hireable          string
	Bio               string
	PublicRepos       string
	PublicGists       string
	Followers         string
	Following         string
	CreatedAt         string
	UpdatedAt         string
	SuspendedAt       string
	TotalPrivateRepos string
	PrivateGists      string
	DiskUsage         string
	Collaborators     string
	Plan              string
}

func (r Recon) User(username string) (response UserResult) {
	user, resp, err := r.Client.Users.Get(r.Ctx, username)
	if resp.StatusCode == 404 {
		r.Logger.Fatal("User not found")
	}
	if err != nil {
		r.Logger.Fatal("Failed to fetch user's information", "err", err)
	}

	r.PrintTitle("👤 User informations")
	u := UserResult{
		Username:          user.GetLogin(),
		ID:                fmt.Sprintf("%d", user.GetID()),
		AvatarURL:         user.GetAvatarURL(),
		GravatarID:        user.GetGravatarID(),
		Name:              user.GetName(),
		Company:           user.GetCompany(),
		Location:          user.GetLocation(),
		Email:             user.GetEmail(),
		Hireable:          fmt.Sprintf("%t", user.GetHireable()),
		Bio:               user.GetBio(),
		PublicRepos:       fmt.Sprintf("%d", user.GetPublicRepos()),
		PublicGists:       fmt.Sprintf("%d", user.GetPublicGists()),
		Followers:         fmt.Sprintf("%d", user.GetFollowers()),
		Following:         fmt.Sprintf("%d", user.GetFollowing()),
		CreatedAt:         user.GetCreatedAt().String(),
		UpdatedAt:         user.GetUpdatedAt().String(),
		SuspendedAt:       user.GetSuspendedAt().String(),
		TotalPrivateRepos: fmt.Sprintf("%d", user.GetTotalPrivateRepos()),
		PrivateGists:      fmt.Sprintf("%d", user.GetPrivateGists()),
		DiskUsage:         fmt.Sprintf("%d", user.GetDiskUsage()),
		Collaborators:     fmt.Sprintf("%d", user.GetCollaborators()),
		Plan:              user.GetPlan().GetName(),
	}
	r.PrintInfo("Username", u.Username)
	r.PrintInfo("ID", u.ID)
	r.PrintInfo("Avatar URL", u.AvatarURL)
	r.PrintInfo("Gravatar ID", u.GravatarID)
	r.PrintInfo("Name", u.Name)
	r.PrintInfo("Company", u.Company)
	r.PrintInfo("Location", u.Location)
	r.PrintInfo("Email", u.Email)
	r.PrintInfo("Hireable", u.Hireable)
	r.PrintInfo("Bio", u.Bio)
	r.PrintInfo("Public Repos", u.PublicRepos)
	r.PrintInfo("Public Gists", u.PublicGists)
	r.PrintInfo("Followers", u.Followers)
	r.PrintInfo("Following", u.Following)
	r.PrintInfo("Created At", u.CreatedAt)
	r.PrintInfo("Updated At", u.UpdatedAt)
	r.PrintInfo("Suspended At", u.SuspendedAt)
	r.PrintInfo("Total Private Repos", u.TotalPrivateRepos)
	r.PrintInfo("Private Gists", u.PrivateGists)
	r.PrintInfo("Disk Usage", u.DiskUsage)
	r.PrintInfo("Collaborators", u.Collaborators)
	r.PrintInfo("Plan", u.Plan)
	r.PrintNewline()

	WaitForRateLimit(resp)
	return u
}
