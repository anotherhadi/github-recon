package main

import (
	recon "github.com/anotherhadi/github-recon/github-recon/username"
	github_recon_settings "github.com/anotherhadi/github-recon/settings"
)

type UsernameResult struct {
	DateTime   string // Now
	Target     string
	TargetType github_recon_settings.TargetType

	User    recon.UserResult
	Socials recon.SocialsResult
	Orgs    recon.OrgsResult

	SshKeys        recon.SshKeysResult
	SshSigningKeys recon.SshSigningKeysResult
	GpgKeys        recon.GpgKeysResult

	CloseFriends recon.CloseFriendsResult

	Commits recon.CommitsResult

	DeepScan recon.DeepScanResult
}

func username(settings github_recon_settings.Settings, datetime string) {
	result := UsernameResult{
		Target:     settings.Target,
		TargetType: settings.TargetType,
		DateTime:   datetime,
	}

	printTitle(settings.Silent, "👤 User informations")
	result.User = recon.User(settings)
	printAvatar(settings, result.User.AvatarURL)
	printStruct(settings, result.User, 0)

	printTitle(settings.Silent, "🐥 Socials")
	result.Socials = recon.Socials(settings)
	printStruct(settings, result.Socials, 0)

	printTitle(settings.Silent, "🏢 Organizations")
	result.Orgs = recon.Orgs(settings)
	printStruct(settings, result.Orgs, 0)

	printTitle(settings.Silent, "🔑 SSH Keys")
	result.SshKeys = recon.SshKeys(settings)
	printStruct(settings, result.SshKeys, 0)

	printTitle(settings.Silent, "🖋️ SSH Signing Keys")
	result.SshSigningKeys = recon.SshSigningKeys(settings)
	printStruct(settings, result.SshSigningKeys, 0)

	printTitle(settings.Silent, "🔐 GPG Keys")
	result.GpgKeys = recon.GpgKeys(settings)
	printStruct(settings, result.GpgKeys, 0)

	printTitle(settings.Silent, "🤝 Close Friends")
	result.CloseFriends = recon.CloseFriends(settings)
	printStruct(settings, result.CloseFriends, 0)

	printTitle(settings.Silent, "📝 Commits")
	result.Commits = recon.Commits(settings)
	printStruct(settings, result.Commits, 0)

	if settings.DeepScan {
		printTitle(settings.Silent, "🔍 Deep Scan")
		result.DeepScan = recon.DeepScan(settings)
		printStruct(settings, result.DeepScan, 0)
	}

	writeJson(settings, result)
}
