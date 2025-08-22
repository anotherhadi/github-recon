package main

import (
	"encoding/json"
	"os"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
)

func writeJson(s github_recon_settings.Settings, data any) {
	if s.JsonOutput == "" {
		return
	}
	file, err := os.Create(s.JsonOutput)
	if err != nil {
		s.Logger.Error("Failed to create JSON file", "err", err)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	as_json, _ := json.MarshalIndent(data, "", "\t")
	_, err = file.Write(as_json)
	if err != nil {
		s.Logger.Error("Failed to write to JSON file", "err", err)
		return
	}
	s.Logger.Info("JSON output written to file", "file", s.JsonOutput)
}
