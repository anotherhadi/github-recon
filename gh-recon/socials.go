package ghrecon

import (
	"encoding/json"
	"fmt"
)

func (r Recon) Socials(username string) {
	resp, err := FetchGitHubAPI(r.client, "", "/users/"+username+"/social_accounts")
	if err != nil {
		r.logger.Error("Failed to fetch socials", "err", err)
		return
	}

	type SocialAccount struct {
		Provider string `json:"provider"`
		URL      string `json:"url"`
	}
	type SocialAccounts []SocialAccount

	var socialAccounts SocialAccounts
	err = json.Unmarshal(resp, &socialAccounts)
	if err != nil {
		r.logger.Error("Failed to unmarshal socials", "err", err)
		return
	}

	if len(socialAccounts) == 0 {
		PrintTitle("🐥 Socials")
		PrintTitle("No Socials found\n")
	} else {
		PrintTitle("🐥 Socials")
		for i, account := range socialAccounts {
			PrintInfo("Social n°", fmt.Sprintf("%d", i))
			PrintInfo("Provider", account.Provider)
			PrintInfo("URL", account.URL)
			fmt.Println()
		}
	}
}
