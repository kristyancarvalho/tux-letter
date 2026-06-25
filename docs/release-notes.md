# Release Notes

## v3.1.1

A quality patch for the generated newsletter.

### Fixed

- The newsletter now always follows the configured `newsletter.language`. With
  `language = "pt-BR"` the article, headings, summary, identity strings, source
  labels and the non-AI fallback are written in Brazilian Portuguese instead of
  defaulting to English.
- Synthesis is deeper: the prompt now demands concrete technical detail,
  practical context and connected stories, so each issue reads like an editorial
  dispatch rather than a shallow roundup of titles.

### Changed

- The AI request uses a synthesis-appropriate timeout instead of reusing the
  short article-fetch timeout.
- The documented default OpenRouter models lead with a stronger synthesis model,
  keeping smaller models as fallbacks.

### Added

- Output-language validation: when a non-English language is configured, clearly
  English output is rejected and a stricter repair attempt is made before the
  deterministic fallback.

### Compatibility

- The configuration schema is unchanged. Existing config files keep working and
  no new fields are required.

### Validation

- `make fmt`
- `make test`
- `make lint`
- `make build`
- `make release-check`

## v3.1.0

Tux Letter now produces a single, cohesive newsletter article instead of a list
of per-source cards.

### Changed

- The newsletter is synthesized into one editorial article with inline numeric
  references and a matching source footer, rather than one summary card per
  source.
- The per-article `why_it_matters` card layout is no longer part of the final
  email.
- The OpenRouter prompt now sends a structured source bundle and asks for one
  unified article that groups related developments and cites sources inline.

### Added

- Full article content extraction (semantic `<article>`, JSON-LD article body,
  Open Graph and main-content heuristics) before summarization, with cleanup and
  a feed-summary fallback.
- Validation of the AI output: it must include a title, body and sources, carry
  inline citations, and never cite a source that was not provided. Invalid output
  triggers one stricter repair attempt before falling back.
- A deterministic non-AI fallback that still renders one article-like digest with
  inline references and a source footer.

### Compatibility

- The configuration schema is unchanged. Existing config files keep working and
  no new fields are required.

### Validation

- `make fmt`
- `make test`
- `make lint`
- `make build`
- `make release-check`

## v3.0.1

### Changed

- Prepare patch release metadata after `v3.0.0`.
- Add a project changelog for release summaries.
- Add a deterministic source archive helper for release packaging.
- Publish a source archive from the release workflow for downstream source
  packages.
- Correct GitHub issue templates so new issues use tux-letter areas and
  secret-safety language.

### Fixed

- Replace the broken `release-source-archive` Makefile target with a working
  archive generator.

### Packaging

- Update AUR metadata for `pkgver=3.0.1`.
- Switch the AUR package source to the GitHub Release source archive so the
  package can pin a stable checksum without relying on GitHub's generated tag
  archives.
- Keep Docker and OpenRouter credentials out of the AUR install path; the
  package builds from source and installs only the `tux-letter` binary and
  documentation.

### Validation

- `make test`
- `make build`
- `make release-check`

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
