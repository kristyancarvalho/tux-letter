# Tux Letter

[![CI](https://github.com/kristyancarvalho/tux-letter/actions/workflows/ci.yml/badge.svg)](https://github.com/kristyancarvalho/tux-letter/actions/workflows/ci.yml)
[![Release](https://github.com/kristyancarvalho/tux-letter/actions/workflows/release.yml/badge.svg)](https://github.com/kristyancarvalho/tux-letter/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/github/license/kristyancarvalho/tux-letter)](LICENSE)
[![Go version](https://img.shields.io/github/go-mod/go-version/kristyancarvalho/tux-letter)](go.mod)
[![Latest release](https://img.shields.io/github/v/release/kristyancarvalho/tux-letter)](https://github.com/kristyancarvalho/tux-letter/releases)
[![AUR](https://img.shields.io/aur/version/tux-letter)](https://aur.archlinux.org/packages/tux-letter)
[![Active milestone](https://img.shields.io/badge/milestone-Release%203.1-2563eb)](https://github.com/kristyancarvalho/tux-letter/milestones)

A lightweight, single-binary newsletter service for Linux and open-source news. Tux Letter collects articles from the websites you list, discovers their feeds, reads the full source articles, and synthesizes them into a single AI-written newsletter article delivered by email.

Tux Letter generates a single AI-written newsletter article synthesized from multiple Linux/open-source sources, with inline references and a source footer. It reads like an editorial intelligence dispatch rather than a feed of separate article cards.

Tux Letter renders a lightweight **cypherpunk-inspired** email digest — a terminal-styled, encrypted-bulletin look that stays comfortable to read in normal email clients.

## Features

- Single Go binary, no database, no Docker required
- URL-only source configuration — no per-site CSS selectors
- RSS/Atom feed discovery with a generic HTML fallback
- Full article content extraction before summarization
- One cohesive AI article with inline numeric citations and a source footer
- OpenRouter synthesis with model fallback and a deterministic non-AI fallback digest
- Cypherpunk HTML email + plain-text fallback
- One-shot and scheduled service modes
- Local atomic state for deduplication

## Install

Download a binary from [Releases](https://github.com/kristyancarvalho/tux-letter/releases), or:

```sh
go install github.com/kristyancarvalho/tux-letter/cmd/tux-letter@latest
```

Arch Linux (AUR):

```sh
yay -S tux-letter
```

## Quick start

```sh
cp tux-letter.example.toml tux-letter.toml
export OPENROUTER_API_KEY=sk-...
tux-letter validate-config
tux-letter sources test
tux-letter once
```

Preview the email identity without sending anything:

```sh
tux-letter preview --output preview.html
```

## Newsletter structure

Each issue is one cohesive article, not a list of per-source cards. Tux Letter
reads the full content of the selected articles, sends a structured source bundle
to OpenRouter, and asks for a single editorial article that groups related
developments. Sources are cited inline as `[1]`, `[2]`, `[3]` and listed in a
matching footer with their URLs. If OpenRouter is unavailable, a deterministic
fallback still produces one article-like digest with the same inline references
and source footer.

## Configuration

Config is discovered in this order: `--config`, `TUX_LETTER_CONFIG`, `./tux-letter.toml`,
`$XDG_CONFIG_HOME/tux-letter/config.toml`, `~/.config/tux-letter/config.toml`.
TOML is preferred; JSON is also supported. Sources are URL-only:

```toml
[sources]
urls = [
  "https://9to5linux.com",
  "https://archlinux.org/news",
  "https://www.linux.com",
]
```

Secrets are read from environment variables named in the config (e.g. `OPENROUTER_API_KEY`,
`TUX_LETTER_SMTP_PASS`) and are never stored in the config file. See
[`tux-letter.example.toml`](tux-letter.example.toml) for the full example.

## Running as a service

Run on a schedule with `tux-letter serve`. A systemd **user** service example:

```ini
[Unit]
Description=Tux Letter newsletter service
After=network-online.target

[Service]
Type=simple
ExecStart=%h/.local/bin/tux-letter serve
Restart=on-failure
RestartSec=30
Environment=OPENROUTER_API_KEY=replace-me

[Install]
WantedBy=default.target
```

```sh
systemctl --user enable --now tux-letter.service
```

## Development

```sh
make fmt      # format
make test     # run tests
make lint     # go vet
make build    # build bin/tux-letter
make coverage # coverage report
make release-check
```

## Releases

Tagging `v*` triggers the release workflow, which builds `linux/amd64` and `linux/arm64`
tarballs plus `checksums.txt`. See [docs/release-notes.md](docs/release-notes.md).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Tux Letter uses an issue-driven workflow with
`release` as the stable branch and `dev` as the integration branch.

## License

Tux Letter is released under the MIT License. See [LICENSE](LICENSE).
