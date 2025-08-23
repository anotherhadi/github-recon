package recon

import (
	"errors"
	"fmt"
	"sort"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/anotherhadi/github-recon/utils"
	"github.com/google/go-github/v72/github"
)

type CloseFriendsResult []CloseFriendResult

type CloseFriendResult struct {
	Username string
	Score    int
}

const (
	maxTargetFollowing       = 50
	maxFollowersForCandidate = 20
	pointsPerCondition       = 1
)

// CloseFriends returns a list of "close friends" for the target user.
// A candidate is considered closer if:
// 1. The target follows fewer than 50 people.
// 2. The candidate has fewer than 20 followers (+1 point).
// 3. The candidate follows the target back (+1 point).
// 4. The candidate shares at least one organization with the target (+1 point).
func CloseFriends(s github_recon_settings.Settings) (results CloseFriendsResult) {
	targetFollowing, resp, err := s.Client.Users.ListFollowing(s.Ctx, s.Target, &github.ListOptions{PerPage: 100})
	if err != nil {
		s.Logger.Error("Failed to fetch target's following list", "err", err)
		return
	}
	utils.WaitForRateLimit(s, resp)

	targetOrgs, err := getOrgs(s, s.Target)
	if err != nil {
		s.Logger.Error("Failed to fetch target's organizations", "err", err)
		targetOrgs = []*github.Organization{}
	}

	if len(targetFollowing) > maxTargetFollowing {
		s.Logger.Info("Skipping close friends check",
			"reason",
			fmt.Sprintf("Target follows %d or more users (limit: %d)", len(targetFollowing), maxTargetFollowing),
		)
		return
	}

	if len(targetFollowing) == 0 {
		return
	}

	for _, candidate := range targetFollowing {
		candidateLogin := candidate.GetLogin()
		if candidateLogin == "" {
			continue
		}

		score := 0

		candidateDetails, userResp, err := s.Client.Users.Get(s.Ctx, candidateLogin)
		if err != nil {
			s.Logger.Warn("Failed to fetch details for candidate",
				"candidate", candidateLogin,
				"err", err,
			)
			if userResp != nil {
				utils.WaitForRateLimit(s, userResp)
			}
			continue
		}
		utils.WaitForRateLimit(s, userResp)

		// Condition: candidate has few followers
		if candidateDetails.GetFollowers() < maxFollowersForCandidate {
			score += pointsPerCondition
		}

		// Condition: candidate follows target back
		followsBack, err := checkIfUserFollows(s, candidateLogin, s.Target)
		if err == nil && followsBack {
			score += pointsPerCondition
		}

		// Condition: same organization
		candidateOrgs, _ := getOrgs(s, candidateLogin)
		if isInSameOrg(targetOrgs, candidateOrgs) {
			score += pointsPerCondition
		}

		// Add candidate if they matched at least one condition
		if score > 0 {
			results = append(results, CloseFriendResult{
				Username: candidateLogin,
				Score:    score,
			})
		}
	}

	if len(results) > 0 {
		sort.Slice(results, func(i, j int) bool {
			return results[i].Score > results[j].Score
		})
	}

	return
}

// checkIfUserFollows checks if sourceUser follows targetUser.
func checkIfUserFollows(s github_recon_settings.Settings, sourceUser, targetUser string) (bool, error) {
	isFollowing, resp, err := s.Client.Users.IsFollowing(s.Ctx, sourceUser, targetUser)
	if err != nil {
		s.Logger.Warn("Error checking if user follows target",
			"source", sourceUser,
			"target", targetUser,
			"err", err,
		)
		if resp != nil {
			utils.WaitForRateLimit(s, resp)
		}
		return false, err
	}
	if resp != nil {
		utils.WaitForRateLimit(s, resp)
	}
	return isFollowing, nil
}

func getOrgs(s github_recon_settings.Settings, user string) ([]*github.Organization, error) {
	orgs, resp, err := s.Client.Organizations.List(s.Ctx, user, nil)
	if err != nil {
		return nil, errors.New("failed to fetch organizations for user")
	}
	utils.WaitForRateLimit(s, resp)
	return orgs, nil
}

func isInSameOrg(orgsA, orgsB []*github.Organization) bool {
	for _, orgA := range orgsA {
		for _, orgB := range orgsB {
			if orgA.GetLogin() == orgB.GetLogin() {
				return true
			}
		}
	}
	return false
}
