package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/google/go-github/v72/github"
)

var (
	Grey  = lipgloss.Color("#7d7d7d")
	Green = lipgloss.Color("#a6e3a1")
	Red   = lipgloss.Color("#f38ba8")

	GreyStyle  = lipgloss.NewStyle().Foreground(Grey)
	GreenStyle = lipgloss.NewStyle().Foreground(Green)
	RedStyle   = lipgloss.NewStyle().Foreground(Red)
)

func header() {
	asciiArt := "        __                       \n  ___ _/ /  _______ _______  ___ \n / _ `/ _ \\/ __/ -_) __/ _ \\/ _ \\\n \\_, /_//_/_/  \\__/\\__/\\___/_//_/\n/___/                            "
	fmt.Println(GreyStyle.Render(lipgloss.JoinVertical(lipgloss.Right, asciiArt, "@anotherhadi\n")))
}

func parseUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if strings.Contains(username, " ") {
		return fmt.Errorf("username cannot contain spaces")
	}
	if strings.Contains(username, "@") {
		return fmt.Errorf("username cannot contain @")
	}
	return nil
}

func FetchGitHubAPI(github *github.Client, token, path string) ([]byte, error) {
	url := "https://api.github.com" + path
	userAgent := "GHRecon/1.0"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for %s: %w", url, err)
	}

	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := github.Client().Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request for %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"request for %s failed with status %d: %s",
			url,
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"error reading response body for %s: %w",
			url,
			err,
		)
	}

	return bodyBytes, nil
}

func PrintInfo(key, value string) {
	if value == "" || value == "0" || value == "0001-01-01 00:00:00 +0000 UTC" {
		return
	}
	fmt.Printf("%s %s\n", GreyStyle.Render(key+":"), GreenStyle.Render(value))
}

func WaitForRateLimit(resp *github.Response) {
	if resp.Rate.Remaining == 0 {
		log.Info(
			"Rate limit reached, waiting for reset... (time:" + resp.Rate.Reset.Time.String() + ")",
		)
		time.Sleep(time.Until(resp.Rate.Reset.Time) + time.Second)
	}
}
