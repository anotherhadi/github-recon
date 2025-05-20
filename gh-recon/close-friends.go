package ghrecon

type CloseFriendsResult struct {
	Login string
	Score int
}

// CloseFriends returns a list of close friends of the user
// To derive this, we check the following:
// 1. The target has less than 50 Following
// 2. The target's following has less than 20 followers

func (r Recon) CloseFriends(username string) (response []CloseFriendsResult) {
	r.PrintTitle("🧑‍🤝‍🧑 Close Friends")

	following, resp, err := r.client.Users.ListFollowing(r.ctx, username, nil)
	if err != nil {
		r.logger.Fatal("Failed to fetch user's close friends", "err", err)
	}
	if len(following) > 50 {
		r.PrintInfo("INFO", "No commits found")
		r.PrintNewline()
		return []CloseFriendsResult{}
	}
	WaitForRateLimit(resp)
	for _, user := range following {
		followers, resp, err := r.client.Users.Get(r.ctx, user.GetLogin())
		WaitForRateLimit(resp)
		if err != nil {
			continue
		}
		if followers.GetFollowers() < 20 {
			response = append(response, CloseFriendsResult{
				Login: user.GetLogin(),
				Score: 1,
			})
		}
	}

	for _, friend := range response {
		r.PrintInfo("Username", "@"+friend.Login)
	}
	r.PrintNewline()
	return
}
