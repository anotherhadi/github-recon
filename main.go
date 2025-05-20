package main

import (
	"context"
	"os"

	ghrecon "github.com/anotherhadi/gh-recon/gh-recon"
	"github.com/charmbracelet/log"
	"github.com/google/go-github/v72/github"
	flag "github.com/spf13/pflag"
)

func main() {
	var username string
	var token string
	var onlyCommitsLeak bool
	var fromEmail string
	var deep bool
	var silent bool
	var jsonFile string
	var excludeRepos string
	flag.StringVarP(&username, "username", "u", "", "GitHub username to analyze")
	flag.StringVarP(&token, "token", "t", "", "GitHub personal access token (e.g. ghp_...)")
	flag.StringVarP(&fromEmail, "email", "e", "", "Search accounts by email address")
	flag.BoolVarP(
		&onlyCommitsLeak,
		"only-commits",
		"c",
		false,
		"Display only commits with author info",
	)
	flag.BoolVarP(
		&deep,
		"deep",
		"d",
		false,
		"Enable deep scan (clone repos, regex search, analyse licenses, etc.)",
	)
	flag.BoolVarP(&silent, "silent", "s", false, "Suppress all non-essential output")
	flag.StringVarP(&jsonFile, "json", "j", "", "Write results to specified JSON file")
	flag.StringVar(
		&excludeRepos,
		"exclude-repo",
		"",
		"Exclude repos from deep scan (comma-separated list)",
	)
	flag.Parse()

	styles := log.DefaultStyles()
	styles.Levels[log.InfoLevel] = styles.Levels[log.InfoLevel].Foreground(ghrecon.Grey)
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: false,
	})
	logger.SetStyles(styles)

	if username == "" && fromEmail == "" {
		logger.Error(
			"Please provide a username with the --username (-u) flag or an email with the --email (-e) flag",
		)
		os.Exit(1)
	} else if username != "" {
		if err := ghrecon.ParseUsername(username); err != nil {
			logger.Error("Invalid username", "err", err)
			os.Exit(1)
		}
	}

	client := github.NewClient(nil)
	if token == "" {
		if !silent {
			logger.Info(
				"It's recommended to set a Github token for better rate limits. You can set it using the --token (-t) flag.",
			)
		}
	} else {
		client = client.WithAuthToken(token)
	}

	ctx := context.Background()

	r := ghrecon.NewRecon(
		client,
		logger,
		ctx,
		silent,
		jsonFile,
	)

	r.Header()

	if fromEmail != "" {
		emailsInfo := r.Email(fromEmail)
		r.WriteJson(
			map[string]any{
				"Authors": emailsInfo,
			},
		)
		return
	}

	if onlyCommitsLeak {
		commitsInfo := r.Commits(username)
		r.WriteJson(
			map[string]any{
				"Authors": commitsInfo,
			},
		)
		return
	}

	userInfo := r.User(username)
	orgsInfo := r.Orgs(username)
	sshKeysInfo := r.SshKeys(username)
	gpgKeysInfo := r.GpgKeys(username)
	sshSigningKeysInfo := r.SshSigningKeys(username)
	socialsInfo := r.Socials(username)
	closeFriendsInfo := r.CloseFriends(username)
	commitsInfo := r.Commits(username)

	results := map[string]any{
		"User":           userInfo,
		"Orgs":           orgsInfo,
		"SSHKeys":        sshKeysInfo,
		"GPGKeys":        gpgKeysInfo,
		"SSHSigningKeys": sshSigningKeysInfo,
		"Socials":        socialsInfo,
		"Commits":        commitsInfo,
		"CloseFriends":   closeFriendsInfo,
	}

	if deep {
		results["Deep"] = r.Deep(username, excludeRepos)
	}

	r.WriteJson(results)
}
