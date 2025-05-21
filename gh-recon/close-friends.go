package ghrecon

import (
	"fmt"
	"sort"
)

type CloseFriendsResult struct {
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
func (r Recon) CloseFriends(username string) (response []CloseFriendsResult) {
	r.PrintTitle("🧑‍🤝‍🧑 Close Friends")

	following, resp, err := r.client.Users.ListFollowing(r.ctx, username, nil)
	if err != nil {
		r.logger.Error("Failed to fetch user's following list", "user", username, "err", err)
		r.PrintNewline()
		return
	}
	WaitForRateLimit(resp)

	if len(following) >= maxFollowingForTarget {
		r.PrintInfo(
			"INFO",
			fmt.Sprintf(
				"%s follows %d or more users (%d). Skipping close friends check.",
				username,
				maxFollowingForTarget,
				len(following),
			),
		)
		r.PrintNewline()
		return
	}

	if len(following) == 0 {
		r.PrintInfo("INFO", fmt.Sprintf("%s is not following anyone.", username))
		r.PrintNewline()
		return
	}

	for _, userBeingFollowedByTarget := range following {
		loginName := userBeingFollowedByTarget.GetLogin()
		if loginName == "" {
			r.logger.Warn("User in following list has an empty login", "target_user", username)
			continue
		}

		currentScore := 0

		userDetails, userResp, userErr := r.client.Users.Get(r.ctx, loginName)
		if userErr != nil {
			r.logger.Warn(
				"Failed to fetch details for followed user",
				"followed_user",
				loginName,
				"err",
				userErr,
			)
			if userResp != nil {
				WaitForRateLimit(userResp)
			}
			continue
		}
		WaitForRateLimit(userResp)

		if userDetails.GetFollowers() < maxFollowersForFollowing {
			currentScore += pointPerCriterion
		}

		followsTargetBack, checkErr := r.checkIfUserFollows(loginName, username)
		if checkErr != nil {
		} else if followsTargetBack {
			currentScore += pointPerCriterion
		}

		if currentScore > 0 {
			response = append(response, CloseFriendsResult{
				Login: loginName,
				Score: currentScore,
			})
		}
	}

	if len(response) == 0 {
		r.PrintInfo(
			"INFO",
			fmt.Sprintf("No close friends found for %s based on the criteria.", username),
		)
	} else {
		sort.Slice(response, func(i, j int) bool {
			return response[i].Score > response[j].Score
		})
		for i, friend := range response {
			r.PrintInfo(
				fmt.Sprintf("Friend n°%d", i+1),
				"@"+friend.Login,
				"Score: "+fmt.Sprintf("%d", friend.Score),
			)
		}
	}

	r.PrintNewline()
	return
}

// checkIfUserFollows checks if sourceUserLogin follows targetUserLogin.
func (r Recon) checkIfUserFollows(sourceUserLogin, targetUserLogin string) (bool, error) {
	isFollowing, resp, err := r.client.Users.IsFollowing(r.ctx, sourceUserLogin, targetUserLogin)
	if err != nil {
		r.logger.Warn("Error checking if user follows target",
			"source_user_checking", sourceUserLogin,
			"target_user_to_check", targetUserLogin,
			"err", err)
		if resp != nil {
			WaitForRateLimit(resp)
		}
		return false, err
	}

	if resp != nil {
		WaitForRateLimit(resp)
	}
	return isFollowing, nil
}
