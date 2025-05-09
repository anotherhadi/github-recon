package main

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/v72/github"
)

func commits(client *github.Client, ctx context.Context, username string) {
	// Commits
	fmt.Println(
		GreyStyle.Render("[")+GreenStyle.Render("+")+GreyStyle.Render("]"),
		GreyStyle.Render("Commits:\n"),
	)

	seenNames := make(map[string]struct{})
	seenEmails := make(map[string]struct{})
	var names, emails []string

	collect := func(order string) error {
		opts := &github.SearchOptions{
			Sort:        "author-date",
			Order:       order, // "desc" ou "asc"
			ListOptions: github.ListOptions{PerPage: 100},
		}
		for page := 1; page <= 10; page++ {
			opts.ListOptions.Page = page
			result, resp, err := client.Search.Commits(
				ctx,
				fmt.Sprintf("author:%s", username),
				opts,
			)
			if err != nil {
				return fmt.Errorf("fetch page %d (%s): %w", page, order, err)
			}
			WaitForRateLimit(resp)
			if len(result.Commits) == 0 {
				break
			}
			for _, item := range result.Commits {
				a := item.Commit.GetAuthor()
				name := a.GetName()
				email := a.GetEmail()
				if _, seen := seenNames[name]; !seen {
					seenNames[name] = struct{}{}
					names = append(names, name)
				}
				if _, seen := seenEmails[email]; !seen {
					seenEmails[email] = struct{}{}
					emails = append(emails, email)
				}
			}
		}
		return nil
	}

	if err := collect("desc"); err != nil {
		log.Error("Failed to fetch commits", "err", err)
		return
	}
	if err := collect("asc"); err != nil {
		log.Error("Failed to fetch commits", "err", err)
		return
	}

	for i, name := range names {
		PrintInfo("Name "+fmt.Sprint(i+1), name)
	}
	for i, email := range emails {
		PrintInfo("Email "+fmt.Sprint(i+1), email)
	}
}
