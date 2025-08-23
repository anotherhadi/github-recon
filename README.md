<div align="center">
    <img src="https://raw.githubusercontent.com/anotherhadi/github-recon/main/.github/assets/logo.png" width="120px" />
</div>

<br>

# Github-Recon 🔍

<p>
    <a href="https://github.com/anotherhadi/github-recon/releases"><img src="https://img.shields.io/github/release/anotherhadi/github-recon.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/anotherhadi/github-recon?tab=doc"><img src="https://godoc.org/github.com/anotherhadi/github-recon?status.svg" alt="GoDoc"></a>
    <a href="https://goreportcard.com/report/github.com/anotherhadi/github-recon"><img src="https://goreportcard.com/badge/github.com/anotherhadi/github-recon" alt="GoReportCard"></a>
</p>

- [🧾 Project Overview](#-project-overview)
- [🚀 Features](#-features)
- [⚠️ Disclaimer](#-disclaimer)
- [📦 Installation](#-installation)
  - [With Go](#with-go)
  - [With Nix/NixOS](#with-nixnixos)
- [🧪 Usage](#-usage)
  - [Flags](#flags)
  - [Token](#token)
- [💡 Examples](#-examples)
- [🕵️‍♂️ Cover your tracks](#-cover-your-tracks)
- [🤝 Contributing](#-contributing)
- [🙏 Credits](#-credits)

## 🧾 Project Overview

Retrieves and aggregates public OSINT data about a GitHub user using Go and the
GitHub API. Finds hidden emails in commit history, previous usernames, friends,
other GitHub accounts, and more.

## 🚀 Features

- Export results to JSON

**From usernames:**

- Retrieve basic user profile information (username, ID, avatar, bio, creation
  date)
- Display avatars directly in the terminal
- List organizations and roles
- Fetch SSH and GPG keys
- Enumerate social accounts
- Extract unique commit authors (name + email)
- Find close friends
- Deep scan option (clone repositories, run regex searches, analyze licenses,
  etc.)
- Use Levenshtein distance for matching usernames and emails

**From emails:**

- Search for a specific email across all GitHub commits
- Spoof an email to discover the associated user account

## ⚠️ Disclaimer

This tool is intended for educational purposes only. Use responsibly and ensure
you have permission to access the data you are querying.

## 📦 Installation

### With Go

```bash
go install github.com/anotherhadi/github-recon@latest
```

### With Nix/NixOS

<details>
<summary>Click to expand</summary>

**From anywhere (using the repo URL):**

```bash
nix run github:anotherhadi/github-recon -- [--flags value] target_username_or_email
```

**Permanent Installation:**

```bash
# add the flake to your flake.nix
{
  inputs = {
    github-recon.url = "github:anotherhadi/github-recon";
  };
}

# then add it to your packages
environment.systemPackages = with pkgs; [ # or home.packages
  github-recon
];
```

</details>

## 🧪 Usage

```bash
github-recon [--flags value] target_username_or_email
```

### Flags

```txt
-t, --token string           Github personal access token (e.g. ghp_aaa...). Can also be set via GITHUB_RECON_TOKEN environment variable. You also need to set the token in $HOME/.config/github-recon/env file if you want to use this tool without passing the token every time. (default "null")
-d, --deepscan               Enable deep scan (clone repos, regex search, analyse licenses, etc.)
    --max-size int           Limit the size of repositories to scan (in MB) (only for deep scan) (default 150)
-e, --exclude-repo strings   Exclude repos from deep scan (comma-separated list, only for deep scan)
-r, --refresh                Refresh the cache (only for deep scan)
-s, --show-source            Show where the information (authors, emails, etc) were found (only for deep scan)
-m, --max-distance int       Maximum Levenshtein distance for matching usernames & emails (only for deep scan) (default 20)
-S, --silent                 Suppress all non-essential output
    --spoof-email            Spoof email (only for email mode) (default true)
-a, --hide-avatar            Hide the avatar in the output
-j, --json string            Write results to specified JSON file
```

### Token

For the best experience, provide a **GitHub Personal Access Token**. Without a
token, you will quickly hit the **rate limit** and have to wait.

- For **basic usage**, you can create a token **without any permissions**.
- For the **email spoofing feature**, you need to add the **`repo`** and
  **`delete_repo`** permissions.

You can set the token in multiple ways:

- **Command-line flag**:

  ```bash
  github-recon -t "ghp_xxx..."
  ```

- **Environment variable**:

  ```bash
  export GITHUB_RECON_TOKEN=ghp_xxx...
  ```

- **Config file**: Create the file `~/.config/github-recon/env` and add:

  ```env
  GITHUB_RECON_TOKEN=ghp_xxx...
  ```

> [!WARNING]
> For safety, it is recommended to create the Personal Access Token on a
> **separate GitHub account** rather than your main account. This way, if
> anything goes wrong, your primary account remains safe.

## 💡 Examples

```bash
github-recon anotherhadi --token ghp_ABC123...
github-recon myemail@gmail.com # Find github accounts by email
github-recon anotherhadi --json output.json --deepscan # Clone the repo and search for leaked email
```

## 🕵️‍♂️ Cover your tracks

Understanding what information about you is publicly visible is the first step
to managing your online presence. github-recon can help you identify your own
publicly available data on GitHub. Here’s how you can take steps to protect your
privacy and security:

- **Review your public profile**: Regularly check your GitHub profile and
  repositories to ensure that you are not unintentionally exposing sensitive
  information.
- **Manage email exposure**: Use GitHub's settings to control which email
  addresses are visible on your profile and in commit history. You can also use
  a no-reply email address for commits. Delete/modify any sensitive information
  in your commit history.
- **Be Mindful of Repository Content**: Avoid including sensitive information in
  your repositories, such as API keys, passwords, emails or personal data. Use
  `.gitignore` to exclude files that contain sensitive information.

You can also use a tool like [TruffleHog](github.com/trufflesecurity/trufflehog)
to scan your repositories specifically for exposed secrets and tokens.

**Useful links:**

- [Blocking command line pushes that expose your personal email address](https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-personal-account-on-github/managing-email-preferences/blocking-command-line-pushes-that-expose-your-personal-email-address)
- [No-reply email address](https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-personal-account-on-github/managing-email-preferences/setting-your-commit-email-address)

## 🤝 Contributing

Feel free to contribute! See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 🙏 Credits

Some features and ideas in this project were inspired by the following tools:

- [gitrecon](https://github.com/GONZOsint/gitrecon) by GONZOsint
- [gitfive](https://github.com/mxrch/gitfive) by mxrch

Big thanks to their authors for sharing their work with the community.
