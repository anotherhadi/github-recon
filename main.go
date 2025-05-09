package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	ghrecon "github.com/anotherhadi/gh-recon/gh-recon"
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
	err := ghrecon.ParseUsername(username)
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

	styles := log.DefaultStyles()
	styles.Levels[log.InfoLevel] = styles.Levels[log.InfoLevel].Foreground(ghrecon.Grey)
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: false,
	})
	logger.SetStyles(styles)

	r := ghrecon.NewRecon(
		client,
		logger,
		ctx,
	)

	ghrecon.Header()
	r.User(username)
	r.Orgs(username)
	r.SshKeys(username)
	r.GpgKeys(username)
	r.SshSigningKeys(username)
	r.Socials(username)
	r.Commits(username)
}
