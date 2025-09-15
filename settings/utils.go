package github_recon_settings

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"
)

// GetToken retrieves the GitHub token from the environment variable or config file
func GetToken() string {
	token := os.Getenv("GITHUB_RECON_TOKEN")
	if token != "" {
		return token
	}

	// Check the $HOME/.config/github-recon/env file for this variable
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "null"
	}
	godotenv.Load(filepath.Join(homedir, ".config/github-recon/env"))
	token = os.Getenv("GITHUB_RECON_TOKEN")
	if token != "" {
		return token
	}

	return "null"
}

func wordSepNormalizeFunc(f *flag.FlagSet, name string) flag.NormalizedName {
	from := []string{".", "_"}
	to := "-"
	for _, sep := range from {
		name = strings.ReplaceAll(name, sep, to)
	}
	return flag.NormalizedName(name)
}
