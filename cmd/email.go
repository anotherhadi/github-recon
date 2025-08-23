package main

import (
	recon "github.com/anotherhadi/github-recon/github-recon/email"
	github_recon_settings "github.com/anotherhadi/github-recon/settings"
)

type EmailResult struct {
	DateTime   string
	Target     string
	TargetType github_recon_settings.TargetType

	Commits  recon.CommitsResult
	Spoofing recon.SpoofingResult
}

func email(settings github_recon_settings.Settings, datetime string) {
	result := EmailResult{
		Target:     settings.Target,
		TargetType: settings.TargetType,
		DateTime:   datetime,
	}

	printTitle(settings.Silent, "👤 Commits author")
	result.Commits = recon.Commits(settings)
	printStruct(settings, result.Commits, 0)

	if settings.SpoofEmail {
		if settings.Token == "null" {
			settings.Logger.Warn("Skipping email spoofing test, please provide a Github token")
		} else {
			printTitle(settings.Silent, "🎭 Spoofing test")
			result.Spoofing = recon.Spoofing(settings)
			if result.Spoofing.AvatarURL != "" {
				printAvatar(settings, result.Spoofing.AvatarURL)
			}
			printStruct(settings, result.Spoofing, 0)
		}
	}
}
