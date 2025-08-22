package utils

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	"github.com/google/go-github/v72/github"
)

func WaitForRateLimit(settings github_recon_settings.Settings, resp *github.Response) {
	if resp.Rate.Remaining == 0 {
		settings.Logger.Info(
			"Rate limit reached, waiting... (time:" + resp.Rate.Reset.Time.String() + ")",
		)
		time.Sleep(time.Until(resp.Rate.Reset.Time) + time.Second)
	}
}

func FetchGitHubAPI(github *github.Client, token, path string) ([]byte, error) {
	url := "https://api.github.com" + path
	userAgent := "GHRecon/1.0"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for %s: %w", url, err)
	}

	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := github.Client().Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request for %s: %w", url, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"request for %s failed with status %d: %s",
			url,
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"error reading response body for %s: %w",
			url,
			err,
		)
	}

	return bodyBytes, nil
}

func DoesFolderExists(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}
	return false
}

func LevenshteinDistance(s1, s2 string) int {
	len1 := len(s1)
	len2 := len(s2)

	dp := make([][]int, len1+1)
	for i := range dp {
		dp[i] = make([]int, len2+1)
	}

	for i := 0; i <= len1; i++ {
		dp[i][0] = i
	}

	for j := 0; j <= len2; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			dp[i][j] = int(
				math.Min(
					float64(dp[i-1][j]+1),
					math.Min(float64(dp[i][j-1]+1), float64(dp[i-1][j-1]+cost)),
				),
			)
		}
	}

	return dp[len1][len2]
}

func SkipResult(name, email string) bool {
	if name == "github-actions[bot]" || name == "dependabot[bot]" || name == "github-actions" {
		return true
	}
	if email == "github-actions[bot]@users.noreply.github.com" ||
		email == "github-actions@github.com" || email == "41898282+github-actions[bot]@users.noreply.github.com" || email == "49699333+dependabot[bot]@users.noreply.github.com" {
		return true
	}
	return false
}
