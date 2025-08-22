package recon

import (
	"fmt"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/anotherhadi/github-recon/utils"
)

type SshKeysResult []SshKeyResult

type SshKeyResult struct {
	Url       string
	Title     string
	CreatedAt string
	Key       string
	ReadOnly  string
	Verified  string
	LastUsed  string
	AddedBy   string
}

func SshKeys(s github_recon_settings.Settings) (response SshKeysResult) {
	sshKeys, resp, err := s.Client.Users.ListKeys(s.Ctx, s.Target, nil)
	if err != nil {
		s.Logger.Error("Failed to fetch ssh keys", "err", err)
		return
	}

	for _, key := range sshKeys {
		k := SshKeyResult{
			Url:       key.GetURL(),
			Title:     key.GetTitle(),
			CreatedAt: key.GetCreatedAt().String(),
			Key:       key.GetKey(),
			ReadOnly:  fmt.Sprintf("%t", key.GetReadOnly()),
			Verified:  fmt.Sprintf("%t", key.GetVerified()),
			LastUsed:  key.GetLastUsed().String(),
			AddedBy:   key.GetAddedBy(),
		}
		response = append(response, k)
	}

	utils.WaitForRateLimit(s, resp)
	return
}

type GpgKeyEmail struct {
	Email    string
	Verified string
}

type GpgKeysResult []GpgKeyResult

type GpgKeyResult struct {
	KeyID        string
	PublicKey    string
	CreatedAt    string
	PrimaryKeyID string
	RawKey       string
	Emails       []GpgKeyEmail
	Subkeys      []GpgKeyResult
}

func GpgKeys(s github_recon_settings.Settings) (response GpgKeysResult) {
	gpgKeys, resp, err := s.Client.Users.ListGPGKeys(s.Ctx, s.Target, nil)
	if err != nil {
		s.Logger.Error("Failed to fetch user's gpg keys", "err", err)
		return
	}

	for _, key := range gpgKeys {
		k := GpgKeyResult{
			KeyID:        key.GetKeyID(),
			PublicKey:    key.GetPublicKey(),
			CreatedAt:    key.GetCreatedAt().String(),
			PrimaryKeyID: fmt.Sprintf("%d", key.GetPrimaryKeyID()),
			RawKey:       key.GetRawKey(),
			Emails:       []GpgKeyEmail{},
			Subkeys:      []GpgKeyResult{},
		}
		for _, email := range key.Emails {
			email := GpgKeyEmail{
				Email:    email.GetEmail(),
				Verified: fmt.Sprintf("%t", email.GetVerified()),
			}
			k.Emails = append(k.Emails, email)
		}
		for _, subkey := range key.Subkeys {
			subkey := GpgKeyResult{
				KeyID:        subkey.GetKeyID(),
				PublicKey:    subkey.GetPublicKey(),
				CreatedAt:    subkey.GetCreatedAt().String(),
				PrimaryKeyID: fmt.Sprintf("%d", subkey.GetPrimaryKeyID()),
				RawKey:       subkey.GetRawKey(),
			}
			k.Subkeys = append(k.Subkeys, subkey)
		}
		response = append(response, k)
	}
	utils.WaitForRateLimit(s, resp)
	return
}

type SshSigningKeysResult []SshSigningKeyResult

type SshSigningKeyResult struct {
	Title     string
	CreatedAt string
	Key       string
}

func SshSigningKeys(s github_recon_settings.Settings) (response SshSigningKeysResult) {
	signingKeys, resp, err := s.Client.Users.ListSSHSigningKeys(
		s.Ctx,
		s.Target,
		nil,
	)
	if err != nil {
		s.Logger.Error("Failed to fetch user's ssh signing keys", "err", err)
		return
	}

	for _, key := range signingKeys {
		k := SshSigningKeyResult{
			Title:     key.GetTitle(),
			CreatedAt: key.GetCreatedAt().String(),
			Key:       key.GetKey(),
		}
		response = append(response, k)
	}

	utils.WaitForRateLimit(s, resp)
	return
}
