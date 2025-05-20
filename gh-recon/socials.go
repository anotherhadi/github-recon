package ghrecon

import (
	"encoding/json"
	"fmt"
)

type SocialResult struct {
	Provider string `json:"provider"`
	URL      string `json:"url"`
}

func (r Recon) Socials(username string) (response []SocialResult) {
	resp, err := FetchGitHubAPI(r.client, "", "/users/"+username+"/social_accounts")
	if err != nil {
		r.logger.Error("Failed to fetch socials", "err", err)
		return
	}

	var socialAccounts []SocialResult
	err = json.Unmarshal(resp, &socialAccounts)
	if err != nil {
		r.logger.Error("Failed to unmarshal socials", "err", err)
		return
	}

	if len(socialAccounts) == 0 {
		r.PrintTitle("🐥 Socials")
		r.PrintInfo("INFO", "No commits found")
	} else {
		r.PrintTitle("🐥 Socials")
		for i, account := range socialAccounts {
			r.PrintInfo("Social n°", fmt.Sprintf("%d", i))
			r.PrintInfo("Provider", account.Provider)
			r.PrintInfo("URL", account.URL)
		}
	}
	r.PrintNewline()

	return socialAccounts
}
