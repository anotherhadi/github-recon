package recon

import (
	"time"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/anotherhadi/github-recon/utils"
)

type EmailResult struct {
	DateTime   string
	Target     string
	TargetType github_recon_settings.TargetType

	Commits  CommitsResult
	Spoofing *SpoofingResult
}

func Email(settings github_recon_settings.Settings) EmailResult {
	result := EmailResult{
		Target:     settings.Target,
		TargetType: settings.TargetType,
		DateTime:   time.Now().String(),
	}

	utils.PrintTitle(settings.Silent, "👤 Commits author")
	result.Commits = Commits(settings)
	utils.PrintStruct(settings, result.Commits, 0)

	if settings.SpoofEmail {
		if settings.Token == "null" {
			settings.Logger.Warn("Skipping email spoofing test, please provide a Github token")
		} else {
			utils.PrintTitle(settings.Silent, "🎭 Spoofing test")
			result.Spoofing = Spoofing(settings)
			if result.Spoofing != nil && result.Spoofing.AvatarURL != "" {
				utils.PrintAvatar(settings, result.Spoofing.AvatarURL)
			}
			utils.PrintStruct(settings, result.Spoofing, 0)
		}
	}

	return result
}
