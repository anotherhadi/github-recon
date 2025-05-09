# GH-Recon

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

```bash
go get github.com/anotherhadi/gh-recon
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
