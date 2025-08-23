package recon

import (
	"math/rand"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/anotherhadi/github-recon/utils"
	"github.com/google/go-github/v72/github"
)

type SpoofingResult struct {
	Username  string
	Name      string
	Email     string
	Url       string
	AvatarURL string
}

func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Spoofing(s github_recon_settings.Settings) (response SpoofingResult) {
	// 1) Create repo (auto-init pour avoir une branche par défaut)
	name := "gh-recon-spoofing-" + RandomString(8)
	private := true
	autoInit := true
	repo, resp, err := s.Client.Repositories.Create(s.Ctx, "", &github.Repository{
		Name:     &name,
		Private:  &private,
		AutoInit: &autoInit,
	})
	if err != nil {
		s.Logger.Error("Error while creating repo", "err", err)
		return
	}
	utils.WaitForRateLimit(s, resp)

	branch := repo.GetDefaultBranch()
	if branch == "" {
		branch = "main"
	}
	refName := "heads/" + branch

	// 2) Commit vide avec auteur spoofé (tree identique au parent)
	author := &github.CommitAuthor{
		Name:  github.String("Spoofed Name"),
		Email: github.String(s.Target),
	}

	ref, resp, err := s.Client.Git.GetRef(s.Ctx, repo.Owner.GetLogin(), name, refName)
	if err != nil {
		s.Logger.Error("Error while getting ref", "err", err)
		s.Logger.Warn("The temp repo was left undeleted", "repo", repo.GetHTMLURL())
		return
	}
	utils.WaitForRateLimit(s, resp)

	parentCommit, resp, err := s.Client.Git.GetCommit(s.Ctx, repo.Owner.GetLogin(), name, ref.GetObject().GetSHA())
	if err != nil {
		s.Logger.Error("Error while getting parent commit", "err", err)
		s.Logger.Warn("The temp repo was left undeleted", "repo", repo.GetHTMLURL())
		return
	}
	utils.WaitForRateLimit(s, resp)

	commit := &github.Commit{
		Author:  author,
		Message: github.String("Spoofed empty commit"),
		Tree:    &github.Tree{SHA: parentCommit.Tree.SHA},
		Parents: []*github.Commit{parentCommit},
	}

	newCommit, resp, err := s.Client.Git.CreateCommit(s.Ctx, repo.Owner.GetLogin(), name, commit, nil)
	if err != nil {
		s.Logger.Error("Error while creating spoofed empty commit", "err", err)
		s.Logger.Warn("The temp repo was left undeleted", "repo", repo.GetHTMLURL())
		return
	}
	utils.WaitForRateLimit(s, resp)

	ref.Object.SHA = newCommit.SHA
	_, resp, err = s.Client.Git.UpdateRef(s.Ctx, repo.Owner.GetLogin(), name, ref, false)
	if err != nil {
		s.Logger.Error("Error while updating ref to spoofed commit", "err", err)
		s.Logger.Warn("The temp repo was left undeleted", "repo", repo.GetHTMLURL())
		return
	}
	utils.WaitForRateLimit(s, resp)

	// 3) Get user
	commits, _, err := s.Client.Repositories.ListCommits(s.Ctx, repo.Owner.GetLogin(), name, nil)
	if err != nil {
		s.Logger.Error("Error while listing commits", "err", err)
		s.Logger.Warn("The temp repo was left undeleted", "repo", repo.GetHTMLURL())
		return
	}

	if len(commits) > 0 {
		last := commits[0]
		response.Username = last.GetAuthor().GetLogin()
		response.Name = last.GetAuthor().GetName()
		response.Email = last.GetAuthor().GetEmail()
		response.Url = last.GetAuthor().GetHTMLURL()
		response.AvatarURL = last.GetAuthor().GetAvatarURL()
	}

	// 4) Cleanup
	_, err = s.Client.Repositories.Delete(s.Ctx, repo.Owner.GetLogin(), name)
	if err != nil {
		s.Logger.Error("Error while deleting repo", "err", err)
	}

	return
}
