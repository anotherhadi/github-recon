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
- Extract unique commit authors (name + email) in both chronological orders

## Disclaimer

This tool is intended for educational purposes only. Use responsibly and ensure you have permission to access the data you are querying.

## Prerequisites

- Go 1.18+
- GitHub Personal Access Token (recommended for higher rate limits)

## Installation

### With Go

```bash
go get github.com/anotherhadi/gh-recon
```

### With Nix/NixOS

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

## Usage

```bash
gh-recon --username TARGET_USER [--token YOUR_TOKEN]
```

### Flags

- `--username`: GitHub username to inspect (required)
- `--token`: Personal Access Token (optional but recommended)

## Example

```bash
gh-recon --username anotherhadi --token ghp_ABC123...
```

## Todo

Feel free to contribute! Here are some ideas:

- Fetch names in License files
- Fetch emails in README files/comments
