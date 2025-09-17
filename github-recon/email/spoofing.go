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

func Spoofing(s github_recon_settings.Settings) (response *SpoofingResult) {
	response = &SpoofingResult{}
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

	defer func() {
		_, err = s.Client.Repositories.Delete(s.Ctx, repo.Owner.GetLogin(), name)
		if err != nil {
			s.Logger.Error("Error while deleting repo", "err", err)
		}
	}()

	branch := repo.GetDefaultBranch()
	if branch == "" {
		branch = "main"
	}
	refName := "heads/" + branch

	authorName := "GITHUB-RECON-SPOOFING"
	authorEmail := s.Target
	author := &github.CommitAuthor{
		Name:  &authorName,
		Email: &authorEmail,
	}

	ref, resp, err := s.Client.Git.GetRef(s.Ctx, repo.Owner.GetLogin(), name, refName)
	if err != nil {
		s.Logger.Error("Error while getting ref", "err", err)
		return
	}
	utils.WaitForRateLimit(s, resp)

	parentCommit, resp, err := s.Client.Git.GetCommit(s.Ctx, repo.Owner.GetLogin(), name, ref.GetObject().GetSHA())
	if err != nil {
		s.Logger.Error("Error while getting parent commit", "err", err)
		return
	}
	utils.WaitForRateLimit(s, resp)

	commitMessage := "Spoofed empty commit"
	commit := &github.Commit{
		Author:  author,
		Message: &commitMessage,
		Tree:    &github.Tree{SHA: parentCommit.Tree.SHA},
		Parents: []*github.Commit{parentCommit},
	}

	newCommit, resp, err := s.Client.Git.CreateCommit(s.Ctx, repo.Owner.GetLogin(), name, commit, nil)
	if err != nil {
		s.Logger.Error("Error while creating spoofed empty commit", "err", err)
		return
	}
	utils.WaitForRateLimit(s, resp)

	ref.Object.SHA = newCommit.SHA
	_, resp, err = s.Client.Git.UpdateRef(s.Ctx, repo.Owner.GetLogin(), name, ref, false)
	if err != nil {
		s.Logger.Error("Error while updating ref to spoofed commit", "err", err)
		return
	}
	utils.WaitForRateLimit(s, resp)

	commits, _, err := s.Client.Repositories.ListCommits(s.Ctx, repo.Owner.GetLogin(), name, nil)
	if err != nil {
		s.Logger.Error("Error while listing commits", "err", err)
		return
	}

	if len(commits) > 1 {
		last := commits[0]
		response.Username = last.GetAuthor().GetLogin()
		response.Name = last.GetAuthor().GetName()
		response.Email = last.GetAuthor().GetEmail()
		response.Url = last.GetAuthor().GetHTMLURL()
		response.AvatarURL = last.GetAuthor().GetAvatarURL()
	} else {
		s.Logger.Error("Only one commit found, something went wrong.", "commits", commits)
	}

	if response.Username == "" && response.Name == "" && response.Email == "" {
		return nil
	}
	return
}
