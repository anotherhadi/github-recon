package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/v72/github"
)

func keys(client *github.Client, ctx context.Context, username string) {
	sshKeys, resp, err := client.Users.ListKeys(ctx, username, nil)
	if err != nil {
		log.Error("Failed to fetch user's keys", "err", err)
		os.Exit(1)
	}
	if len(sshKeys) == 0 {
		fmt.Println(
			GreyStyle.Render("[")+
				RedStyle.Render("x")+
				GreyStyle.Render("]"),
			GreyStyle.Render("No keys found for user "+username),
		)
	} else {
		fmt.Println(
			GreyStyle.Render("[")+
				GreenStyle.Render("+")+
				GreyStyle.Render("]"),
			GreyStyle.Render("Keys:\n"),
		)
	}
	for _, key := range sshKeys {
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
	WaitForRateLimit(resp)

	gpgKeys, resp, err := client.Users.ListGPGKeys(ctx, username, nil)
	if err != nil {
		log.Error("Failed to fetch user's GPG keys", "err", err)
	}

	if len(gpgKeys) == 0 {
		fmt.Println(
			GreyStyle.Render("[")+
				RedStyle.Render("x")+
				GreyStyle.Render("]"),
			GreyStyle.Render("No GPG keys found\n"),
		)
	} else {
		fmt.Println(
			GreyStyle.Render("[")+
				GreenStyle.Render("+")+
				GreyStyle.Render("]"),
			GreyStyle.Render("GPG Keys:\n"),
		)
	}
	for _, key := range gpgKeys {
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
		for _, email := range key.Emails {
			PrintInfo("\tEmail", email.GetEmail())
			PrintInfo("\tVerified", fmt.Sprintf("%t", email.GetVerified()))
		}
		PrintInfo("Subkeys", fmt.Sprintf("%d", len(key.Subkeys)))
		for _, subkey := range key.Subkeys {
			PrintInfo("\tSubkey ID", fmt.Sprintf("%d", subkey.GetID()))
			PrintInfo("\tSubkey Key ID", subkey.GetKeyID())
			PrintInfo("\tSubkey Created At", subkey.GetCreatedAt().String())
			PrintInfo("\tSubkey Expires At", subkey.GetExpiresAt().String())
			PrintInfo("\tSubkey Can Sign", fmt.Sprintf("%t", subkey.GetCanSign()))
			PrintInfo(
				"\tSubkey Can Encrypt Comms",
				fmt.Sprintf("%t", subkey.GetCanEncryptComms()),
			)
			PrintInfo(
				"\tSubkey Can Encrypt Storage",
				fmt.Sprintf("%t", subkey.GetCanEncryptStorage()),
			)
			PrintInfo("\tSubkey Can Certify", fmt.Sprintf("%t", subkey.GetCanCertify()))
			PrintInfo("\tSubkey Primary Key ID", fmt.Sprintf("%d", subkey.GetPrimaryKeyID()))
			PrintInfo("\tSubkey Raw Key", subkey.GetRawKey())
			PrintInfo("\tSubkey Public Key", subkey.GetPublicKey())
		}
		fmt.Println()
	}
	WaitForRateLimit(resp)

	signingKeys, resp, err := client.Users.ListSSHSigningKeys(
		ctx,
		username,
		nil,
	)
	if err != nil {
		log.Error("Failed to fetch user's signing keys", "err", err)
	}
	if len(signingKeys) == 0 {
		fmt.Println(
			GreyStyle.Render("[")+
				RedStyle.Render("x")+
				GreyStyle.Render("]"),
			GreyStyle.Render("No signing keys found\n"),
		)
	} else {
		fmt.Println(
			GreyStyle.Render("[")+
				GreenStyle.Render("+")+
				GreyStyle.Render("]"),
			GreyStyle.Render("Signing Keys:\n"),
		)
	}
	for _, key := range signingKeys {
		PrintInfo("ID", fmt.Sprintf("%d", key.GetID()))
		PrintInfo("Key", key.GetKey())
		PrintInfo("Title", key.GetTitle())
		PrintInfo("Created At", key.GetCreatedAt().String())
		fmt.Println()
	}

	WaitForRateLimit(resp)
}
