package main

import (
	"time"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
)

func main() {
	settings := github_recon_settings.GetSettings()
	datetime := time.Now().String()

	if !settings.Silent {
		header()

		printStruct(settings, struct {
			Target     string
			TargetType string
			DateTime   string
		}{
			Target:     settings.Target,
			TargetType: string(settings.TargetType),
			DateTime:   datetime,
		}, 0)
	}

	if settings.TargetType == github_recon_settings.TargetUsername {
		username(settings, datetime)
	} else {
		email(settings, datetime)
	}
}
