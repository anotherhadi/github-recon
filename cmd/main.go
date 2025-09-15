package main

import (
	recon_email "github.com/anotherhadi/github-recon/github-recon/email"
	recon_username "github.com/anotherhadi/github-recon/github-recon/username"
	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/anotherhadi/github-recon/utils"
)

func main() {
	settings := github_recon_settings.GetSettings()

	if !settings.Silent {
		utils.Header()

		utils.PrintStruct(settings, struct {
			Target     string
			TargetType string
		}{
			Target:     settings.Target,
			TargetType: string(settings.TargetType),
		}, 0)
	}

	if settings.TargetType == github_recon_settings.TargetUsername {
		result := recon_username.Username(settings)
		writeJson(settings, result)
	} else {
		result := recon_email.Email(settings)
		writeJson(settings, result)
	}
}
