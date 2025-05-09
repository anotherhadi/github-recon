package ghrecon

import (
	"fmt"
)

func (r Recon) SshKeys(username string) {
	sshKeys, resp, err := r.client.Users.ListKeys(r.ctx, username, nil)
	if err != nil {
		r.logger.Error("Failed to fetch ssh keys", "err", err)
	} else if len(sshKeys) == 0 {
		PrintTitle("🔑 SSH Keys")
		r.logger.Info("No SSH Keys found\n")
	} else {
		PrintTitle("🔑 SSH Keys")
		for i, key := range sshKeys {
			PrintInfo("Key n°", fmt.Sprintf("%d", i))
			PrintInfo("ID", fmt.Sprintf("%d", key.GetID()))
			PrintInfo("URL", key.GetURL())
			PrintInfo("Title", key.GetTitle())
			PrintInfo("Created At", key.GetCreatedAt().String())
			PrintInfo("Key", key.GetKey())
			PrintInfo("Read Only", fmt.Sprintf("%t", key.GetReadOnly()))
			PrintInfo("Verified", fmt.Sprintf("%t", key.GetVerified()))
			PrintInfo("Last Used", key.GetLastUsed().String())
			PrintInfo("Added By", key.GetAddedBy())
			fmt.Println()
		}
	}
	WaitForRateLimit(resp)
}

func (r Recon) GpgKeys(username string) {
	gpgKeys, resp, err := r.client.Users.ListGPGKeys(r.ctx, username, nil)
	if err != nil {
		r.logger.Error("Failed to fetch user's gpg keys", "err", err)
	} else if len(gpgKeys) == 0 {
		PrintTitle("🗝️ GPG Keys")
		r.logger.Info("No GPG Keys found\n")
	} else {
		PrintTitle("🗝️ GPG Keys")
		for i, key := range gpgKeys {
			PrintInfo("Key n°", fmt.Sprintf("%d", i))
			PrintInfo("ID", fmt.Sprintf("%d", key.GetID()))
			PrintInfo("Key ID", key.GetKeyID())
			PrintInfo("Public Key", key.GetPublicKey())
			PrintInfo("Created At", key.GetCreatedAt().String())
			PrintInfo("Expires At", key.GetExpiresAt().String())
			PrintInfo("Can Sign", fmt.Sprintf("%t", key.GetCanSign()))
			PrintInfo("Can Encrypt Comms", fmt.Sprintf("%t", key.GetCanEncryptComms()))
			PrintInfo("Can Encrypt Storage", fmt.Sprintf("%t", key.GetCanEncryptStorage()))
			PrintInfo("Can Certify", fmt.Sprintf("%t", key.GetCanCertify()))
			PrintInfo("Primary Key ID", fmt.Sprintf("%d", key.GetPrimaryKeyID()))
			PrintInfo("Raw Key", key.GetRawKey())
			PrintInfo("Emails", fmt.Sprintf("%d", len(key.Emails)))
			for j, email := range key.Emails {
				PrintInfo("  Email n°", fmt.Sprintf("%d", j))
				PrintInfo("  Email", email.GetEmail())
				PrintInfo("  Verified", fmt.Sprintf("%t", email.GetVerified()))
			}
			PrintInfo("Subkeys", fmt.Sprintf("%d", len(key.Subkeys)))
			for j, subkey := range key.Subkeys {
				PrintInfo("  Subkey n°", fmt.Sprintf("%d", j))
				PrintInfo("  Subkey ID", fmt.Sprintf("%d", subkey.GetID()))
				PrintInfo("  Subkey Key ID", subkey.GetKeyID())
				PrintInfo("  Subkey Created At", subkey.GetCreatedAt().String())
				PrintInfo("  Subkey Expires At", subkey.GetExpiresAt().String())
				PrintInfo("  Subkey Can Sign", fmt.Sprintf("%t", subkey.GetCanSign()))
				PrintInfo(
					"  Subkey Can Encrypt Comms",
					fmt.Sprintf("%t", subkey.GetCanEncryptComms()),
				)
				PrintInfo(
					"  Subkey Can Encrypt Storage",
					fmt.Sprintf("%t", subkey.GetCanEncryptStorage()),
				)
				PrintInfo("  Subkey Can Certify", fmt.Sprintf("%t", subkey.GetCanCertify()))
				PrintInfo("  Subkey Primary Key ID", fmt.Sprintf("%d", subkey.GetPrimaryKeyID()))
				PrintInfo("  Subkey Raw Key", subkey.GetRawKey())
				PrintInfo("  Subkey Public Key", subkey.GetPublicKey())
			}
			fmt.Println()
		}
	}
	WaitForRateLimit(resp)
}

func (r Recon) SshSigningKeys(username string) {
	signingKeys, resp, err := r.client.Users.ListSSHSigningKeys(
		r.ctx,
		username,
		nil,
	)
	if err != nil {
		r.logger.Error("Failed to fetch user's ssh signing keys", "err", err)
	} else if len(signingKeys) == 0 {
		PrintTitle("📝 SSH Signing Keys")
		r.logger.Info("No SSH Signing Keys found\n")
	} else {
		PrintTitle("📝 SSH Signing Keys")
		for i, key := range signingKeys {
			PrintInfo("Key n°", fmt.Sprintf("%d", i))
			PrintInfo("ID", fmt.Sprintf("%d", key.GetID()))
			PrintInfo("Key", key.GetKey())
			PrintInfo("Title", key.GetTitle())
			PrintInfo("Created At", key.GetCreatedAt().String())
			fmt.Println()
		}
	}
	WaitForRateLimit(resp)
}
