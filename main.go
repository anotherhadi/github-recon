package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/v72/github"
)

func main() {
	var username string
	var token string
	flag.StringVar(&username, "username", "", "Target username")
	flag.StringVar(&token, "token", "", "Github token")
	flag.Parse()

	if username == "" {
		fmt.Println("Please provide a username with the --username flag")
		os.Exit(1)
	}
	err := parseUsername(username)
	if err != nil {
		log.Error("Invalid username", "err", err)
		os.Exit(1)
	}

	client := github.NewClient(nil)
	if token == "" {
		log.Info(
			"It's recommended to set a Github token for better rate limits. You can set it using the --token flag.",
		)
	} else {
		client = client.WithAuthToken(token)
	}

	ctx := context.Background()

	header()
	userInfo(client, ctx, username)
	orgs(client, ctx, username)
	keys(client, ctx, username)
	socials(client, username)
	commits(client, ctx, username)
}
