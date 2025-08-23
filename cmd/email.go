package main

import (
	recon "github.com/anotherhadi/github-recon/github-recon/email"
	github_recon_settings "github.com/anotherhadi/github-recon/settings"
)

type EmailResult struct {
	DateTime   string
	Target     string
	TargetType github_recon_settings.TargetType

	Commit recon.CommitsResult
}

func email(settings github_recon_settings.Settings, datetime string) {
	result := EmailResult{
		Target:     settings.Target,
		TargetType: settings.TargetType,
		DateTime:   datetime,
	}

	printTitle(settings.Silent, "👤 Commits author")
	result.Commit = recon.Email(settings)
	printStruct(settings, result.Commit, 0)
}
