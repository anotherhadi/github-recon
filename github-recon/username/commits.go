package recon

import (
	"fmt"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/anotherhadi/github-recon/utils"
	"github.com/google/go-github/v72/github"
)

type CommitsResult []CommitResult

type CommitResult struct {
	Name         string
	Email        string
	Occurrences  int
	FirstFoundIn string
}

func Commits(s github_recon_settings.Settings) (response CommitsResult) {
	results := make(map[string]CommitResult)

	collect := func(date string) error {
		for page := 1; page <= 10; page++ {
			result, resp, err := s.Client.Search.Commits(
				s.Ctx,
				fmt.Sprintf("author:%s author-date:%s", s.Target, date),
				&github.SearchOptions{
					Sort:        "author-date",
					Order:       "desc",
					ListOptions: github.ListOptions{PerPage: 100, Page: page},
				},
			)
			if err != nil {
				return fmt.Errorf("fetch page %d (%s): %w", page, date, err)
			}
			utils.WaitForRateLimit(s, resp)
			if len(result.Commits) == 0 {
				break
			}
			for _, item := range result.Commits {
				emails := []string{
					item.Commit.GetAuthor().GetEmail(),
					item.Commit.GetCommitter().GetEmail(),
					item.GetAuthor().GetEmail(),
					item.GetCommitter().GetEmail(),
				}

				names := []string{
					item.Commit.GetAuthor().GetName(),
					item.Commit.GetCommitter().GetName(),
					item.GetAuthor().GetName(),
					item.GetCommitter().GetName(),
				}

				for i := range names {
					if utils.SkipResult(names[i], emails[i]) {
						continue
					}
					if names[i] == "" || emails[i] == "" {
						continue
					}
					if _, seen := results[names[i]+" - "+emails[i]]; !seen {
						author := CommitResult{
							Name:        names[i],
							Email:       emails[i],
							Occurrences: 1,
							FirstFoundIn: item.GetRepository().Owner.GetLogin() + "/" + item.GetRepository().
								GetName(),
						}
						results[names[i]+" - "+emails[i]] = author
					} else if i == 0 {
						result := results[names[i]+" - "+emails[i]]
						result.Occurrences++
						results[names[i]+" - "+emails[i]] = result
					}
				}
			}
		}
		return nil
	}

	// Range of dates to bypass the limit of 1000 results
	for _, date := range []string{
		"<2023-01-01", "2023-01-01..2023-12-31",
		"2024-01-01..2024-05-31",
		"2024-06-01..2024-12-31",
		"2025-01-01..2025-05-31",
		"2025-06-01..2025-12-31",
		">2026-01-01",
	} {
		if err := collect(date); err != nil {
			s.Logger.Error("Failed to fetch commits", "err", err, "date", date)
		}
	}

	for _, result := range results {
		response = append(response, result)
	}

	return
}
