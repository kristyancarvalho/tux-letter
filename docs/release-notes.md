# Release Notes

## v3.0.0

Tux Letter 3.0 is a full rewrite into a lightweight, single Go binary with no
required database and no required Docker runtime.

### Highlights

- **URL-only sources.** Sources are declared as plain URLs. Per-site CSS
  selectors are gone; the scraper discovers feeds and extracts articles
  generically.
- **Feed discovery and parsing.** RSS 2.0, RDF and Atom feeds are discovered
  from page `<link>` tags and parsed into normalized articles.
- **Generic HTML scraping fallback.** When no feed is available, articles are
  extracted from JSON-LD, Open Graph metadata, `<article>` blocks and
  heading/anchor heuristics, with navigation, tag, login, ad and social links
  filtered out.
- **Normalization, dedupe and local state.** Articles are normalized
  (canonical URLs, collapsed whitespace, content hashes) and deduplicated
  across sources. A local JSON state file tracks seen articles and run metadata
  with atomic writes.
- **OpenRouter AI with model fallback.** Configured models are attempted in
  order; transient and rate-limit failures fall back to the next model and
  invalid output is retried once. A non-AI digest is generated when every model
  fails.
- **Cypherpunk newsletter.** Email-safe, terminal-inspired HTML digest with a
  plain-text fallback, rendered from validated structured data.
- **SMTP delivery.** Email is delivered over SMTP using environment variables;
  stdout output is used when email is disabled. No secrets are logged.
- **Scheduled service mode.** `tux-letter serve` runs on a timezone-aware daily
  schedule with graceful shutdown on `SIGINT`/`SIGTERM`.

### Commands

```sh
tux-letter once
tux-letter serve
tux-letter validate-config
tux-letter sources test
tux-letter preview
tux-letter --version
```

### Configuration

Configuration is TOML-first (JSON is also accepted) and discovered from
`--config`, `TUX_LETTER_CONFIG`, the working directory, or
`~/.config/tux-letter/`. Secrets are read from environment variables and are
never stored in committed config. See `tux-letter.example.toml`.

### Install

- Download a Linux binary from the GitHub release assets.
- `go install github.com/kristyancarvalho/tux-letter/cmd/tux-letter@latest`
- Arch Linux: install `tux-letter` from the AUR.

### Release assets

```text
tux-letter_3.0.0_linux_amd64.tar.gz
tux-letter_3.0.0_linux_arm64.tar.gz
checksums.txt
```

### License

Released under the MIT License.
