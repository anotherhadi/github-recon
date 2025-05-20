package ghrecon

import (
	"fmt"
)

type OrgResult struct {
	Login       string
	ID          string
	URL         string
	Description string
}

func (r Recon) Orgs(username string) (response []OrgResult) {
	orgs, resp, err := r.client.Organizations.List(r.ctx, username, nil)
	if err != nil {
		r.logger.Error("Failed to fetch organizations", "err", err)
	} else if len(orgs) == 0 {
		r.PrintTitle("🏢 Organizations")
		r.PrintInfo("INFO", "No Organizations found")
	} else {
		r.PrintTitle("🏢 Organizations")
		for i, org := range orgs {
			o := OrgResult{
				Login:       org.GetLogin(),
				ID:          fmt.Sprintf("%d", org.GetID()),
				URL:         org.GetURL(),
				Description: org.GetDescription(),
			}
			r.PrintInfo("Organization n°", fmt.Sprintf("%d", i))
			r.PrintInfo("Login", o.Login)
			r.PrintInfo("ID", o.ID)
			r.PrintInfo("URL", o.URL)
			r.PrintInfo("Description", o.Description)
			r.PrintNewline()
			response = append(response, o)
		}
	}
	r.PrintNewline()
	WaitForRateLimit(resp)
	return
}
