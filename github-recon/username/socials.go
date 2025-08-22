package recon

import (
	"encoding/json"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/anotherhadi/github-recon/utils"
)

type socialResultInput struct {
	Provider string `json:"provider"`
	URL      string `json:"url"`
}

type SocialsResult []socialResult

type socialResult struct {
	Provider string
	URL      string
}

func Socials(s github_recon_settings.Settings) (response SocialsResult) {
	resp, err := utils.FetchGitHubAPI(s.Client, "", "/users/"+s.Target+"/social_accounts")
	if err != nil {
		s.Logger.Error("Failed to fetch socials", "err", err)
		return
	}

	var socialAccounts []socialResultInput
	err = json.Unmarshal(resp, &socialAccounts)
	if err != nil {
		s.Logger.Error("Failed to unmarshal socials", "err", err)
		return
	}

	socials := []socialResult{}
	for _, account := range socialAccounts {
		socials = append(socials, socialResult{
			URL:      account.URL,
			Provider: account.Provider,
		})
	}

	return socials
}
