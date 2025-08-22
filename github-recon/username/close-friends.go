package recon

import (
	"fmt"
	"sort"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/anotherhadi/github-recon/utils"
)

type CloseFriendsResult []CloseFriendResult

type CloseFriendResult struct {
	Login string
	Score int
}

const (
	maxFollowingForTarget    = 50
	maxFollowersForFollowing = 20
	pointPerCriterion        = 1
)

// CloseFriends returns a list of close friends of the user
// To derive this, we check the following:
// 1. The target has less than 50 Following
// 2. The target's following has less than 20 followers (+1 point)
// 3. The target's following follows the target (+1 point)

func CloseFriends(s github_recon_settings.Settings) (response CloseFriendsResult) {
	following, resp, err := s.Client.Users.ListFollowing(s.Ctx, s.Target, nil)
	if err != nil {
		s.Logger.Error("Failed to fetch user's following list", "err", err)
		return
	}

	utils.WaitForRateLimit(s, resp)

	if len(following) >= maxFollowingForTarget {
		s.Logger.Info("Skipping close friends check", "reason", fmt.Sprintf("Target follows %d or more users (%d)", maxFollowingForTarget, len(following)))
		return
	}

	if len(following) == 0 {
		return
	}

	for _, userBeingFollowedByTarget := range following {
		loginName := userBeingFollowedByTarget.GetLogin()
		if loginName == "" {
			continue
		}

		currentScore := 0

		userDetails, userResp, userErr := s.Client.Users.Get(s.Ctx, loginName)
		if userErr != nil {
			s.Logger.Warn(
				"Failed to fetch details for followed user",
				"followed_user",
				loginName,
				"err",
				userErr,
			)
			if userResp != nil {
				utils.WaitForRateLimit(s, userResp)
			}
			continue
		}
		utils.WaitForRateLimit(s, userResp)

		if userDetails.GetFollowers() < maxFollowersForFollowing {
			currentScore += pointPerCriterion
		}

		followsTargetBack, checkErr := checkIfUserFollows(s, loginName, s.Target)
		if checkErr != nil {
			continue
		} else if followsTargetBack {
			currentScore += pointPerCriterion
		}

		if currentScore > 0 {
			response = append(response, CloseFriendResult{
				Login: loginName,
				Score: currentScore,
			})
		}
	}

	if len(response) == 0 {
		return
	} else {
		sort.Slice(response, func(i, j int) bool {
			return response[i].Score > response[j].Score
		})
	}

	return
}

// checkIfUserFollows checks if sourceUserLogin follows targetUserLogin.
func checkIfUserFollows(s github_recon_settings.Settings, sourceUserLogin, targetUserLogin string) (bool, error) {
	isFollowing, resp, err := s.Client.Users.IsFollowing(s.Ctx, sourceUserLogin, targetUserLogin)
	if err != nil {
		s.Logger.Warn("Error checking if user follows target",
			"source_user_checking", sourceUserLogin,
			"target_user_to_check", targetUserLogin,
			"err", err)
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
