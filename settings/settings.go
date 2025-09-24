package github_recon_settings

import (
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	"context"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/v72/github"
)

type TargetType string

const (
	TargetUsername TargetType = "Username"
	TargetEmail    TargetType = "Email"
)

type Settings struct {
	Token         string
	Target        string
	TargetType    TargetType
	ShowSource    bool
	Refresh       bool
	MaxRepoSize   int
	ExcludedRepos []string
	JsonOutput    string
	Silent        bool
	DeepScan      bool
	MaxDistance   int
	PrintAvatar   bool
	SpoofEmail    bool
	Trufflehog    bool

	// Internal
	Client *github.Client
	Logger *log.Logger
	Ctx    context.Context
}

func GetDefaultSettings() Settings {
	return Settings{
		Token:         "null",
		Target:        "",
		TargetType:    TargetUsername,
		ShowSource:    false,
		Refresh:       false,
		MaxRepoSize:   150,
		ExcludedRepos: []string{},
		JsonOutput:    "",
		Silent:        false,
		DeepScan:      false,
		MaxDistance:   20,
		PrintAvatar:   true,
		SpoofEmail:    true,
		Trufflehog:    true,

		Client: github.NewClient(nil),
		Logger: log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller:    false,
			ReportTimestamp: false,
		}),
		Ctx: context.WithValue(context.Background(), github.SleepUntilPrimaryRateLimitResetWhenRateLimited, true),
	}
}

func GetSettings() (settings Settings, err error) {
	settings = GetDefaultSettings()
	//// Flag settings
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "github-recon [flags] <target username or email>\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	flag.CommandLine.SortFlags = false

	//// Flags
	flag.StringVarP(&settings.Token, "token", "t", settings.Token, "Github personal access token (e.g. ghp_aaa...). Can also be set via GITHUB_RECON_TOKEN environment variable. You also need to set the token in $HOME/.config/github-recon/env file if you want to use this tool without passing the token every time.")

	// DeepScan
	flag.BoolVarP(&settings.DeepScan, "deepscan", "d", settings.DeepScan, "Enable deep scan (clone repos, regex search, analyse licenses, etc.)")
	flag.IntVar(
		&settings.MaxRepoSize,
		"max-size",
		settings.MaxRepoSize,
		"Limit the size of repositories to scan (in MB) (only for deep scan)",
	)
	flag.StringSliceVarP(
		&settings.ExcludedRepos,
		"exclude-repo",
		"e",
		settings.ExcludedRepos,
		"Exclude repos from deep scan (comma-separated list, only for deep scan)",
	)
	flag.BoolVarP(
		&settings.Refresh,
		"refresh",
		"r",
		settings.Refresh,
		"Refresh the cache (only for deep scan)",
	)
	flag.BoolVarP(
		&settings.ShowSource,
		"show-source",
		"s",
		settings.ShowSource,
		"Show where the information (authors, emails, etc) were found (only for deep scan)",
	)
	flag.IntVarP(
		&settings.MaxDistance,
		"max-distance",
		"m",
		settings.MaxDistance,
		"Maximum Levenshtein distance for matching usernames & emails (only for deep scan)",
	)
	flag.BoolVar(
		&settings.Trufflehog,
		"trufflehog",
		settings.Trufflehog,
		"Run trufflehog on cloned repositories (only for deep scan)",
	)

	flag.BoolVarP(&settings.Silent, "silent", "S", settings.Silent, "Suppress all non-essential output")
	flag.BoolVarP(&settings.SpoofEmail, "spoof-email", "", settings.SpoofEmail, "Spoof email (only for email mode)")
	flag.BoolVarP(&settings.PrintAvatar, "print-avatar", "a", settings.PrintAvatar, "Show the avatar in the output")
	flag.StringVarP(&settings.JsonOutput, "json", "j", settings.JsonOutput, "Write results to specified JSON file")

	//// Parse
	flag.Parse()

	//// Tail
	nonFlagArgs := flag.Args()
	if len(nonFlagArgs) > 1 {
		settings.Logger.Error("Please provide only one target (username or email)")
		flag.Usage()
		os.Exit(1)
	} else if len(nonFlagArgs) == 0 {
		settings.Logger.Error("Please provide a target (username or email)")
		flag.Usage()
		os.Exit(1)
	}

	settings.Target = flag.Arg(0)
	settings.Target = strings.TrimPrefix(settings.Target, "@") // Remove the @ of the username

	if strings.Contains(settings.Target, " ") {
		err = fmt.Errorf("target cannot contain spaces")
	}

	if strings.Contains(settings.Target, "@") {
		settings.TargetType = TargetEmail
	} else {
		settings.TargetType = TargetUsername
	}

	// If token is not set via flag, get it from env
	if settings.Token == "null" || settings.Token == "" {
		settings.Token = GetToken()
	}

	if settings.Token == "null" || settings.Token == "" {
		settings.Logger.Warn("No Github token provided. You might hit the rate limit. Check the help menu for more information.")
	} else {
		settings.Client = settings.Client.WithAuthToken(settings.Token)
	}

	return
}
