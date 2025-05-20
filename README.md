<div align="center">
    <img src="https://raw.githubusercontent.com/anotherhadi/gh-recon/main/.github/assets/logo.png" width="120px" />
</div>

<br>

# GH-Recon

<p>
    <a href="https://github.com/anotherhadi/gh-recon/releases"><img src="https://img.shields.io/github/release/anotherhadi/gh-recon.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/anotherhadi/gh-recon?tab=doc"><img src="https://godoc.org/github.com/anotherhadi/gh-recon?status.svg" alt="GoDoc"></a>
    <a href="https://goreportcard.com/report/github.com/anotherhadi/gh-recon"><img src="https://goreportcard.com/badge/github.com/anotherhadi/gh-recon" alt="GoReportCard"></a>
</p>

## Project Overview

Fetches and aggregates public OSINT data for a GitHub user, leveraging Go and the GitHub API.

## Features

- Retrieve basic user profile information (username, ID, avatar, bio, creation dates)
- List organizations and roles
- Fetch SSH and GPG keys
- Enumerate social accounts
- Extract unique commit authors (name + email)
- Find close friends
- Find Github accounts using an email address
- Export results to JSON
- Deep scan option (clone repositories, regex search, analyze licenses, etc.)

## Disclaimer

This tool is intended for educational purposes only. Use responsibly and ensure you have permission to access the data you are querying.

## Prerequisites

- Go 1.18+
- GitHub Personal Access Token (recommended for higher rate limits): Create a GitHub API token with no permissions/no scope. This will be equivalent to public GitHub access, but it will allow access to use the GitHub Search API.

## Installation

### With Go

```bash
go get github.com/anotherhadi/gh-recon
```

### With Nix/NixOS

<details>
<summary>Click to expand</summary>

**From anywhere (using the repo URL):**

```bash
nix run github:anotherhadi/gh-recon -- --username TARGET_USER [--token YOUR_TOKEN]
```

**Permanent Installation:**

```bash
# add the flake to your flake.nix
{
  inputs = {
    gh-recon.url = "github:anotherhadi/gh-recon";
  };
}

# then add it to your packages
environment.systemPackages = with pkgs; [ # or home.packages
  gh-recon
];
```

</details>

## Usage

```bash
gh-recon --username TARGET_USER [--token YOUR_TOKEN]
```

### Flags

```txt
  -deep
     Enable deep scan (clone repos, regex search, analyse licenses, etc.)
  -email string
     Search accounts by email address
  -json string
     Write results to specified JSON file
  -only-commits
     Display only commits with author info
  -silent
     Suppress all non-essential output
  -token string
     GitHub personal access token (e.g. ghp_...)
  -username string
     GitHub username to analyze
```

## Example

```bash
gh-recon --username anotherhadi --token ghp_ABC123...
gh-recon --email myemail@gmail.com --token ghp_ABC123...
gh-recon --username anotherhadi --json output.json --deep
```

## Contributing

Feel free to contribute! See [CONTRIBUTING.md](CONTRIBUTING.md) for details.
