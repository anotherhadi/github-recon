package recon

import (
	"time"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/anotherhadi/github-recon/utils"
)

type UsernameResult struct {
	DateTime   string // Now
	Target     string
	TargetType github_recon_settings.TargetType

	User    UserResult
	Socials SocialsResult
	Orgs    OrgsResult

	SshKeys        SshKeysResult
	SshSigningKeys SshSigningKeysResult
	GpgKeys        GpgKeysResult

	CloseFriends CloseFriendsResult

	Commits CommitsResult

	DeepScan DeepScanResult
}

func Username(settings github_recon_settings.Settings) (result UsernameResult, err error) {
	result = UsernameResult{
		Target:     settings.Target,
		TargetType: settings.TargetType,
		DateTime:   time.Now().String(),
	}

	utils.PrintTitle(settings.Silent, "👤 User informations")
	result.User, err = User(settings)
	if err != nil {
		return
	}
	if result.User == (UserResult{}) {
		return
	}
	utils.PrintAvatar(settings, result.User.AvatarURL)
	utils.PrintStruct(settings, result.User, 0)

	utils.PrintTitle(settings.Silent, "🐥 Socials")
	result.Socials = Socials(settings)
	utils.PrintStruct(settings, result.Socials, 0)

	utils.PrintTitle(settings.Silent, "🏢 Organizations")
	result.Orgs = Orgs(settings)
	utils.PrintStruct(settings, result.Orgs, 0)

	utils.PrintTitle(settings.Silent, "🔑 SSH Keys")
	result.SshKeys = SshKeys(settings)
	utils.PrintStruct(settings, result.SshKeys, 0)

	utils.PrintTitle(settings.Silent, "🖋️ SSH Signing Keys")
	result.SshSigningKeys = SshSigningKeys(settings)
	utils.PrintStruct(settings, result.SshSigningKeys, 0)

	utils.PrintTitle(settings.Silent, "🔐 GPG Keys")
	result.GpgKeys = GpgKeys(settings)
	utils.PrintStruct(settings, result.GpgKeys, 0)

	utils.PrintTitle(settings.Silent, "🤝 Close Friends")
	result.CloseFriends = CloseFriends(settings)
	utils.PrintStruct(settings, result.CloseFriends, 0)

	utils.PrintTitle(settings.Silent, "📝 Commits")
	result.Commits = Commits(settings)
	utils.PrintStruct(settings, result.Commits, 0)

	if settings.DeepScan {
		utils.PrintTitle(settings.Silent, "🔍 Deep Scan")
		result.DeepScan = DeepScan(settings)
		utils.PrintStruct(settings, result.DeepScan, 0)
	}

	return
}
