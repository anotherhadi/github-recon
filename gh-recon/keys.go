package ghrecon

import (
	"fmt"
)

type SSHKeyResult struct {
	ID        string
	Url       string
	Title     string
	CreatedAt string
	Key       string
	ReadOnly  string
	Verified  string
	LastUsed  string
	AddedBy   string
}

func (r Recon) SshKeys(username string) (response []SSHKeyResult) {
	sshKeys, resp, err := r.client.Users.ListKeys(r.ctx, username, nil)
	if err != nil {
		r.logger.Error("Failed to fetch ssh keys", "err", err)
	} else if len(sshKeys) == 0 {
		r.PrintTitle("🔑 SSH Keys")
		r.PrintInfo("INFO", "No SSH Keys found")
	} else {
		r.PrintTitle("🔑 SSH Keys")
		for i, key := range sshKeys {
			k := SSHKeyResult{
				ID:        fmt.Sprintf("%d", key.GetID()),
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
			r.PrintInfo("Key n°", fmt.Sprintf("%d", i))
			r.PrintInfo("ID", k.ID)
			r.PrintInfo("URL", k.Url)
			r.PrintInfo("Title", k.Title)
			r.PrintInfo("Created At", k.CreatedAt)
			r.PrintInfo("Key", k.Key)
			r.PrintInfo("Read Only", k.ReadOnly)
			r.PrintInfo("Verified", k.Verified)
			r.PrintInfo("Last Used", k.LastUsed)
			r.PrintInfo("Added By", k.AddedBy)
			if i != len(sshKeys)-1 {
				r.PrintNewline()
			}
		}
	}
	r.PrintNewline()
	WaitForRateLimit(resp)
	return
}

type GPGKeyEmail struct {
	Email    string
	Verified string
}
type GPGKeyResult struct {
	ID           string
	KeyID        string
	PublicKey    string
	CreatedAt    string
	PrimaryKeyID string
	RawKey       string
	Emails       []GPGKeyEmail
	Subkeys      []GPGKeyResult
}

func (r Recon) GpgKeys(username string) (response []GPGKeyResult) {
	gpgKeys, resp, err := r.client.Users.ListGPGKeys(r.ctx, username, nil)
	if err != nil {
		r.logger.Error("Failed to fetch user's gpg keys", "err", err)
	} else if len(gpgKeys) == 0 {
		r.PrintTitle("🗝️ GPG Keys")
		r.PrintInfo("INFO", "No GPG Keys found")
	} else {
		r.PrintTitle("🗝️ GPG Keys")
		for i, key := range gpgKeys {
			k := GPGKeyResult{
				ID:           fmt.Sprintf("%d", key.GetID()),
				KeyID:        key.GetKeyID(),
				PublicKey:    key.GetPublicKey(),
				CreatedAt:    key.GetCreatedAt().String(),
				PrimaryKeyID: fmt.Sprintf("%d", key.GetPrimaryKeyID()),
				RawKey:       key.GetRawKey(),
				Emails:       []GPGKeyEmail{},
				Subkeys:      []GPGKeyResult{},
			}
			for _, email := range key.Emails {
				email := GPGKeyEmail{
					Email:    email.GetEmail(),
					Verified: fmt.Sprintf("%t", email.GetVerified()),
				}
				k.Emails = append(k.Emails, email)
			}
			for _, subkey := range key.Subkeys {
				subkey := GPGKeyResult{
					ID:           fmt.Sprintf("%d", subkey.GetID()),
					KeyID:        subkey.GetKeyID(),
					PublicKey:    subkey.GetPublicKey(),
					CreatedAt:    subkey.GetCreatedAt().String(),
					PrimaryKeyID: fmt.Sprintf("%d", subkey.GetPrimaryKeyID()),
					RawKey:       subkey.GetRawKey(),
				}
				k.Subkeys = append(k.Subkeys, subkey)
			}
			response = append(response, k)
			r.PrintInfo("Key n°", fmt.Sprintf("%d", i))
			r.PrintInfo("ID", k.ID)
			r.PrintInfo("Key ID", k.KeyID)
			r.PrintInfo("Public Key", k.PublicKey)
			r.PrintInfo("Created At", k.CreatedAt)
			r.PrintInfo("Primary Key ID", k.PrimaryKeyID)
			r.PrintInfo("Raw Key", k.RawKey)
			r.PrintInfo("Emails", fmt.Sprintf("%d", len(k.Emails)))
			for j, email := range k.Emails {
				r.PrintInfo("  Email n°", fmt.Sprintf("%d", j))
				r.PrintInfo("  Email", email.Email)
				r.PrintInfo("  Verified", email.Verified)
				if j != len(k.Emails)-1 {
					r.PrintNewline()
				}
			}
			r.PrintInfo("Subkeys", fmt.Sprintf("%d", len(k.Subkeys)))
			for j, subkey := range k.Subkeys {
				r.PrintInfo("  Subkey n°", fmt.Sprintf("%d", j))
				r.PrintInfo("  Subkey ID", subkey.ID)
				r.PrintInfo("  Subkey Key ID", subkey.KeyID)
				r.PrintInfo("  Subkey Created At", subkey.CreatedAt)
				r.PrintInfo("  Subkey Primary Key ID", subkey.PrimaryKeyID)
				r.PrintInfo("  Subkey Raw Key", subkey.RawKey)
				if j != len(k.Subkeys)-1 {
					r.PrintNewline()
				}
			}
			if i != len(gpgKeys)-1 {
				r.PrintNewline()
			}
		}
	}
	r.PrintNewline()
	WaitForRateLimit(resp)
	return
}

type SSHSigningKeyResult struct {
	ID        string
	Title     string
	CreatedAt string
	Key       string
}

func (r Recon) SshSigningKeys(username string) (response []SSHSigningKeyResult) {
	signingKeys, resp, err := r.client.Users.ListSSHSigningKeys(
		r.ctx,
		username,
		nil,
	)
	if err != nil {
		r.logger.Error("Failed to fetch user's ssh signing keys", "err", err)
	} else if len(signingKeys) == 0 {
		r.PrintTitle("📝 SSH Signing Keys")
		r.PrintInfo("INFO", "No SSH Signing Keys found")
	} else {
		r.PrintTitle("📝 SSH Signing Keys")
		for i, key := range signingKeys {
			k := SSHSigningKeyResult{
				ID:        fmt.Sprintf("%d", key.GetID()),
				Title:     key.GetTitle(),
				CreatedAt: key.GetCreatedAt().String(),
				Key:       key.GetKey(),
			}
			r.PrintInfo("Key n°", fmt.Sprintf("%d", i))
			r.PrintInfo("ID", k.ID)
			r.PrintInfo("Title", k.Title)
			r.PrintInfo("Created At", k.CreatedAt)
			r.PrintInfo("Key", k.Key)
			if i != len(signingKeys)-1 {
				r.PrintNewline()
			}
			response = append(response, k)
		}
	}
	WaitForRateLimit(resp)
	r.PrintNewline()
	return response
}
