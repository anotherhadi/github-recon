package recon

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/anotherhadi/github-recon/utils"
	"github.com/google/go-github/v72/github"
)

type Authors []Author

type Author struct {
	Name        string
	Levenshtein int
	Email       string
	FoundIn     []string
}

type Emails []Email

type Email struct {
	Email       string
	Levenshtein int
	FoundIn     []string
}

type Secrets []Secret

type Secret struct {
	Repositorie string
	Raw         map[string]any
}

type DeepScanResult struct {
	Authors Authors
	Emails  Emails
	Secrets Secrets
}

type Repositorie struct {
	Repository string
	Owner      string
	Name       string
	Size       int
}

func DeepScan(s github_recon_settings.Settings) (response DeepScanResult) {
	repositories := []Repositorie{}
	repos, resp, err := s.Client.Repositories.ListByUser(
		s.Ctx,
		s.Target,
		&github.RepositoryListByUserOptions{
			Type: "all",
		},
	)

	if err != nil {
		s.Logger.Error("Failed to fetch repositories", "err", err)
		return
	}

	for _, repo := range repos {
		if slices.Contains(s.ExcludedRepos, repo.GetName()) ||
			slices.Contains(s.ExcludedRepos, repo.GetOwner().GetLogin()+"/"+repo.GetName()) {
			continue
		}

		maxRepoSize := s.MaxRepoSize * 1024

		if repo.GetSize() > maxRepoSize {
			s.Logger.Info("Skipping repository due to size", "repo", repo.GetOwner().GetLogin()+"/"+repo.GetName(), "size_MB", repo.GetSize()/1024, "max_size_MB", maxRepoSize/1024)
			continue
		}

		repositories = append(repositories, Repositorie{
			Repository: repo.GetCloneURL(),
			Owner:      repo.GetOwner().GetLogin(),
			Name:       repo.GetName(),
			Size:       repo.GetSize(),
		})
	}
	utils.WaitForRateLimit(s, resp)

	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		s.Logger.Error("Git is not installed", "err", err)
		return
	}

	tmp_folder := "/tmp/ghrecon-" + s.Target

	if utils.DoesFolderExists(tmp_folder) {
		if s.Refresh {
			s.Logger.Info("Deleting existing folder", "path", tmp_folder)
			err := os.RemoveAll(tmp_folder)
			if err != nil {
				s.Logger.Error("Failed to delete existing folder", "path", tmp_folder, "err", err)
				return
			}
		}
	}

	for _, repo := range repositories {
		destination := tmp_folder + "/" + repo.Owner + "/" + repo.Name
		if utils.DoesFolderExists(destination) {
			s.Logger.Info("Directory already downloaded, skipping", "repo", repo.Owner+"/"+repo.Name, "path", destination)
			continue
		}

		s.Logger.Info("Cloning repository", "repo", repo.Owner+"/"+repo.Name, "path", destination, "size_MB", repo.Size/1024)

		cmd := exec.Command(
			"git",
			"clone",
			repo.Repository,
			destination,
		)
		err := cmd.Run()
		if err != nil {
			s.Logger.Error(
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
	s.Logger.Info("Cloned all repositories", "path", tmp_folder)

	if len(repositories) == 0 {
		s.Logger.Info("No repositories found for the user, skipping deep scan.")
		return
	}

	authorOccurrences := Authors{}
	mapAuthorToIndex := make(map[string]int)
	for _, repo := range repositories {
		destination := tmp_folder + "/" + repo.Owner + "/" + repo.Name
		if !utils.DoesFolderExists(filepath.Join(destination, ".git")) {
			s.Logger.Error(
				"No .git directory found, cannot run git log.",
				"repo",
				repo.Owner+"/"+repo.Name,
				"path",
				destination,
			)
		} else {
			gitLogCmd := exec.Command("git", "log", "--all", "--format=%aN <%aE>")
			gitLogCmd.Dir = destination
			logOutput, logErr := gitLogCmd.Output()

			if logErr != nil {
				if exitErr, ok := logErr.(*exec.ExitError); ok {
					s.Logger.Error("Failed to execute git log (ExitError)", "repo", repo.Owner+"/"+repo.Name, "stderr", string(exitErr.Stderr), "err", logErr)
				} else {
					s.Logger.Error("Failed to execute git log", "repo", repo.Owner+"/"+repo.Name, "err", logErr)
				}
			} else {
				lines := strings.Split(string(logOutput), "\n")
				repoIdentifier := repo.Owner + "/" + repo.Name

				for _, line := range lines {
					trimmedLine := strings.TrimSpace(line)
					if trimmedLine == "" {
						continue
					}

					if index, exists := mapAuthorToIndex[trimmedLine]; exists {
						isRepoListed := false
						for _, foundRepo := range authorOccurrences[index].FoundIn {
							if foundRepo == repoIdentifier {
								isRepoListed = true
								break
							}
						}
						if !isRepoListed {
							authorOccurrences[index].FoundIn = append(authorOccurrences[index].FoundIn, repoIdentifier)
							slices.Sort(authorOccurrences[index].FoundIn)
						}
					} else {
						parts := strings.SplitN(trimmedLine, " <", 2)
						var authorName, authorEmail string
						if len(parts) == 2 {
							authorName = parts[0]
							authorEmail = strings.TrimSuffix(parts[1], ">")
						} else if len(parts) == 1 {
							authorName = "-"
							authorEmail = strings.TrimPrefix(strings.TrimSuffix(parts[0], ">"), "<")
						} else {
							s.Logger.Error("Malformed author line from git log", "line", trimmedLine, "repo", repoIdentifier)
							continue
						}

						authorOccurrences = append(authorOccurrences, Author{
							Name:        authorName,
							Email:       authorEmail,
							FoundIn:     []string{repoIdentifier},
							Levenshtein: levenshteinDistanceAuthor(s.Target, authorName, authorEmail),
						})
						mapAuthorToIndex[trimmedLine] = len(authorOccurrences) - 1
					}
				}
			}
		}
	}
	slices.SortFunc(authorOccurrences, func(a, b Author) int {
		if a.Levenshtein != b.Levenshtein {
			return a.Levenshtein - b.Levenshtein
		}
		return 1
	})

	authors := Authors{}
	for _, author := range authorOccurrences {
		if author.Levenshtein > s.MaxDistance {
			continue
		}
		if utils.SkipResult(author.Name, author.Email) {
			continue
		}
		authors = append(authors, author)
	}

	s.Logger.Info("Searching for emails in cloned repositories", "path", tmp_folder)
	emailsFound, err := findEmailsAndOccurrencesInDir(tmp_folder, s.Target)
	if err != nil {
		s.Logger.Error("Failed to find emails in directory", "err", err)
		return
	}
	slices.SortFunc(emailsFound, func(a, b Email) int {
		if a.Levenshtein != b.Levenshtein {
			return a.Levenshtein - b.Levenshtein
		}
		return 1
	})

	emails := Emails{}
	for _, email := range emailsFound {
		if email.Levenshtein > s.MaxDistance {
			continue
		}
		emails = append(emails, email)
	}

	response.Authors = authors
	response.Emails = emails

	s.Logger.Info("Searching for secrets in cloned repositories", "path", tmp_folder)
	if s.Trufflehog {
		cmd := exec.Command("trufflehog", "--version")
		if err := cmd.Run(); err != nil {
			s.Logger.Warn("Trufflehog is not installed, skipping secret scanning.")
		} else {
			secrets, err := truffleHog(tmp_folder)
			if err != nil {
				s.Logger.Error("Failed to run trufflehog", "err", err)
			} else {
				response.Secrets = secrets
			}
		}
	}

	return
}

func truffleHog(tmpFolder string) (Secrets, error) {
	allSecrets := Secrets{}

	directories, err := os.ReadDir(tmpFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to read tmp folder: %w", err)
	}

	for _, dir := range directories {
		if !dir.IsDir() {
			continue
		}

		innerPath := filepath.Join(tmpFolder, dir.Name())
		innerDirectories, err := os.ReadDir(innerPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read inner tmp folder: %w", err)
		}

		for _, innerDir := range innerDirectories {
			if !innerDir.IsDir() {
				continue
			}

			repoPath := filepath.Join(innerPath, innerDir.Name())
			cmd := exec.Command("trufflehog", "git", "file://"+repoPath, "--json", "--log-level=-1", "--results=verified")
			output, err := cmd.Output()

			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					if exitErr.ExitCode() > 1 {
						return nil, fmt.Errorf("failed to execute trufflehog (ExitError): %s", string(exitErr.Stderr))
					}
				} else {
					return nil, fmt.Errorf("failed to execute trufflehog: %w", err)
				}
			}

			decoder := json.NewDecoder(strings.NewReader(string(output)))
			for decoder.More() {
				var result map[string]any
				if err := decoder.Decode(&result); err != nil {
					return nil, fmt.Errorf("failed to parse trufflehog output: %w", err)
				}
				allSecrets = append(allSecrets, Secret{
					Repositorie: dir.Name() + "/" + innerDir.Name(),
					Raw:         result,
				})
			}
		}
	}

	return allSecrets, nil
}

func findEmailsAndOccurrencesInDir(rootPath string, username string) (Emails, error) {
	emailLocations := make(map[string]map[string]bool)
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	normalizedRootPath := filepath.Clean(rootPath)

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() {
			if strings.Contains(path, ".git/logs/") {
				return nil
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			currentFileEmails := emailRegex.FindAllString(string(content), -1)
			if len(currentFileEmails) > 0 {
				relativePath, errRel := filepath.Rel(normalizedRootPath, path)
				if errRel != nil {
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

	var results Emails
	for email, pathSet := range emailLocations {
		var paths []string
		for path := range pathSet {
			paths = append(paths, path)
		}
		results = append(results, Email{
			Email: email, FoundIn: paths,
			Levenshtein: utils.LevenshteinDistance(username, strings.SplitN(email, "@", 2)[0]),
		})
	}

	return results, nil
}

func levenshteinDistanceAuthor(target, name, email string) int {
	if strings.Contains(email, "@") {
		email = strings.SplitN(email, "@", 2)[0]
	}
	return slices.Min([]int{utils.LevenshteinDistance(target, name), utils.LevenshteinDistance(target, email)})
}
