package ghrecon

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/google/go-github/v72/github"
)

func folderExists(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}
	return false
}

type EmailOccurrence struct {
	Email   string
	FoundIn []string
}

func findEmailsAndOccurrencesInDir(rootPath string) ([]EmailOccurrence, error) {
	emailLocations := make(map[string]map[string]bool)
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	normalizedRootPath := filepath.Clean(rootPath)

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Can't access %s: %v\n", path, err)
			return err
		}
		if !d.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Can't read %s: %v\n", path, err)
				return nil
			}

			currentFileEmails := emailRegex.FindAllString(string(content), -1)
			if len(currentFileEmails) > 0 {
				relativePath, errRel := filepath.Rel(normalizedRootPath, path)
				if errRel != nil {
					fmt.Printf("Can't find the relative path %s: %v\n", path, errRel)
					relativePath = path
				}

				for _, email := range currentFileEmails {
					if len(email) > 12 {
						if _, ok := emailLocations[email]; !ok {
							emailLocations[email] = make(map[string]bool)
						}
						emailLocations[email][relativePath] = true
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var results []EmailOccurrence
	for email, pathSet := range emailLocations {
		var paths []string
		for path := range pathSet {
			paths = append(paths, path)
		}
		results = append(results, EmailOccurrence{Email: email, FoundIn: paths})
	}

	return results, nil
}

type DeepResult struct {
	Repository string
	Owner      string
	Name       string
}

func (r Recon) Deep(username, excludeRepos string) (response []DeepResult) {
	excludeReposList := strings.Split(excludeRepos, ",")
	repos, resp, err := r.client.Repositories.ListByUser(
		r.ctx,
		username,
		&github.RepositoryListByUserOptions{
			Type: "all",
		},
	)
	if err != nil {
		r.logger.Error("Failed to fetch repositories", "err", err)
		return
	}

	r.PrintTitle("📦 Repositories")
	if len(repos) == 0 {
		r.PrintInfo("INFO", "No repositories found")
	} else {
		for _, repo := range repos {
			response = append(response, DeepResult{
				Repository: repo.GetCloneURL(),
				Owner:      repo.GetOwner().GetLogin(),
				Name:       repo.GetName(),
			})
		}
	}
	WaitForRateLimit(resp)

	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		r.PrintInfo("ERROR", "Git is not installed, please install it to use this feature")
		return
	}

	tmp_folder := "/tmp/ghrecon-" + username
	for _, repo := range response {
		if slices.Contains(excludeReposList, repo.Name) ||
			slices.Contains(excludeReposList, repo.Owner+"/"+repo.Name) {
			r.PrintInfo("INFO", "Skipping repository", repo.Owner+"/"+repo.Name)
			continue
		}
		r.PrintInfo(
			"Downloading repository",
			repo.Owner+"/"+repo.Name,
		)

		destination := tmp_folder + "/" + repo.Owner + "/" + repo.Name
		if folderExists(destination) {
			r.PrintInfo("INFO", "Directory already exists, skipping")
			continue
		}

		cmd := exec.Command(
			"git",
			"clone",
			repo.Repository,
			destination,
		)
		err := cmd.Run()
		if err != nil {
			r.logger.Error(
				"ERROR",
				"Failed to clone repository",
				"err",
				err,
				"repo",
				repo.Repository,
			)
			continue
		}
	}
	r.PrintInfo("INFO", "Cloned all repositories to "+tmp_folder)

	r.PrintInfo("INFO", "Now searching for emails in cloned repositories, this may take a while...")
	results, err := findEmailsAndOccurrencesInDir(tmp_folder)
	if err != nil {
		r.logger.Error("Failed to find emails in directory", "err", err)
		return
	}

	if len(results) == 0 {
		r.PrintInfo("INFO", "No emails found")
	} else {
		r.PrintInfo("INFO", "Found emails:")
		for _, email := range results {
			r.PrintInfo("Email", email.Email)
			r.PrintInfo("Found in", tmp_folder, email.FoundIn...)
		}
	}

	r.PrintNewline()
	return
}
