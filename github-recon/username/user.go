package recon

import (
	"fmt"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/anotherhadi/github-recon/utils"
)

type UserResult struct {
	Username          string
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

func User(s github_recon_settings.Settings) (response UserResult, err error) {
	user, resp, err := s.Client.Users.Get(s.Ctx, s.Target)
	if resp.StatusCode == 404 {
		return UserResult{}, fmt.Errorf("user not found with username: %s", s.Target)
	}
	if err != nil {
		return UserResult{}, fmt.Errorf("failed to fetch user's information")
	}

	u := UserResult{
		Username:          user.GetLogin(),
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

	utils.WaitForRateLimit(s, resp)
	return u, nil
}
