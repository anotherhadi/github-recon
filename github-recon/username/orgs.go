package recon

import (
	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/anotherhadi/github-recon/utils"
)

type OrgsResult []OrgResult

type OrgResult struct {
	Name        string
	URL         string
	Description string
}

func Orgs(s github_recon_settings.Settings) (response OrgsResult) {
	orgs, resp, err := s.Client.Organizations.List(s.Ctx, s.Target, nil)
	if err != nil {
		s.Logger.Error("Failed to fetch organizations", "err", err)
		return
	}

	for _, org := range orgs {
		o := OrgResult{
			Name:        org.GetLogin(),
			URL:         org.GetURL(),
			Description: org.GetDescription(),
		}
		response = append(response, o)
	}

	utils.WaitForRateLimit(s, resp)

	return
}
