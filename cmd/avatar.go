package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
	gopixels "github.com/saran13raj/go-pixels"
)

func printAvatar(settings github_recon_settings.Settings, url string) {
	if settings.HideAvatar || url == "" || settings.Silent {
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	tmpfile, err := os.CreateTemp("", "avatar-*.png")
	if err != nil {
		return
	}
	defer os.Remove(tmpfile.Name())

	_, err = io.Copy(tmpfile, resp.Body)
	if err != nil {
		return
	}

	output, err := gopixels.FromImagePath(tmpfile.Name(), 30, 25, "halfcell", true)

	if err != nil {
		return
	}
	fmt.Println(output + "\n")
}
