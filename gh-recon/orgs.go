package ghrecon

import (
	"fmt"
)

func (r Recon) Orgs(username string) {
	orgs, resp, err := r.client.Organizations.List(r.ctx, username, nil)
	if err != nil {
		r.logger.Error("Failed to fetch organizations", "err", err)
	} else if len(orgs) == 0 {
		PrintTitle("🏢 Organizations")
		r.logger.Info("No Organizations found\n")
	} else {
		PrintTitle("🏢 Organizations")
		for i, org := range orgs {
			PrintInfo("Orgs n°", fmt.Sprintf("%d", i))
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
	}
	WaitForRateLimit(resp)
}
