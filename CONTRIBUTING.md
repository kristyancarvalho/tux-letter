# Contributing to tux-letter

`tux-letter` uses an issue-driven workflow with `dev` as the integration branch and `release` as the stable, default branch.

The current release target is:

```text
v3.1.0
```

GitHub milestone:

```text
Release 3.1
```

## Branching

Use this workflow:

1. Start from `dev`.
2. Create a focused local staging branch.
3. Implement one issue at a time.
4. Validate locally.
5. Merge the staging branch back into local `dev`.
6. Re-run validation on `dev`.
7. Push only `dev`.

Do not commit directly to `release`.

Do not push staging branches unless the repository owner explicitly asks for it.

Local staging branch examples:

```text
stage/config-rewrite
stage/scraper-feed-discovery
stage/openrouter-fallback
stage/aur-packaging
```

Branches that should remain local-only:

```text
stage/*
staging/*
fix/*
wip/*
temp/*
```

## Issue-Driven Development

No implementation should happen without a GitHub issue.

Every issue for the 3.0 rewrite must be:

* assigned to `kristyancarvalho`;
* associated with the `Release 3.0` milestone;
* labeled with the correct type and area labels;
* small enough to be implemented and validated independently.

Recommended labels:

```text
type:feature
type:fix
type:docs
type:ci
type:release
type:packaging
area:architecture
area:config
area:scraper
area:ai
area:newsletter
area:service
area:docs
area:aur
priority:high
release:3.0
```

When using the GitHub CLI, prefer explicit issue creation:

```sh
gh issue create \
  --title "Rewrite configuration system for URL-only sources" \
  --body "Implement the config system described in ./specs/SPEC.md." \
  --assignee "kristyancarvalho" \
  --milestone "Release 3.0" \
  --label "type:feature,area:config,priority:high,release:3.0"
```

## Commit Style

Use English commit messages.

Preferred format:

```text
<type>/<scope>: <summary>; <action> (issue #<id>)
```

Examples:

```text
feat/config: add TOML config loader; implements (issue #12)
feat/scraper: add RSS feed discovery; implements (issue #13)
fix/ai: handle invalid OpenRouter responses; fixes (issue #14)
test/state: cover atomic state writes; implements (issue #15)
docs/readme: document AUR installation; updates (issue #16)
chore/release: prepare v3.0.0 metadata; updates (issue #17)
```

Common types:

```text
feat
fix
docs
test
refactor
chore
ci
build
release
```

Common scopes:

```text
config
scraper
feed
ai
newsletter
mail
state
service
cli
docs
aur
release
tests
```

## Specs and Local Planning

The `./specs/` directory is for local planning only.

It must never be committed.

Required local spec file for the rewrite:

```text
./specs/SPEC.md
```

The repository `.gitignore` must include:

```text
/specs/
```

Agents and contributors should read `./specs/SPEC.md` before making changes, but must not stage or commit it.

## Repository Structure

Preferred structure:

```text
cmd/tux-letter/
internal/
.github/workflows/
docs/
packaging/aur/
scripts/release/
```

Public documentation should stay simple.

Use:

```text
README.md
CONTRIBUTING.md
CHANGELOG.md
docs/release-notes.md
```

Avoid adding many public planning documents. Keep implementation planning in `./specs/`.

## Development Commands

Use the Makefile whenever possible.

```sh
make fmt
make test
make lint
make build
make coverage
make release-check
```

Direct Go equivalents:

```sh
gofmt -w .
go test ./...
go vet ./...
go build ./...
```

Before merging any local staging branch back into `dev`, run:

```sh
make fmt
make test
make lint
make build
```

Before release work, run:

```sh
make release-check
```

## Testing Rules

Tests should be deterministic.

Do not make unit tests depend on live websites, live OpenRouter calls, SMTP delivery, or the local system clock when avoidable.

Use fixtures for:

* RSS feeds;
* Atom feeds;
* HTML pages;
* JSON-LD blocks;
* OpenGraph metadata;
* AI responses;
* invalid AI responses.

Recommended fixture locations:

```text
internal/feed/testdata/
internal/scrape/testdata/
internal/ai/testdata/
```

Networked checks should be opt-in and should not run in normal CI unless explicitly designed as integration tests.

## Configuration Rules

The source list must be URL-only.

Do not reintroduce required per-site CSS selectors.

Good:

```toml
[sources]
urls = [
  "https://9to5linux.com",
  "https://archlinux.org/news"
]
```

Avoid:

```json
{
  "name": "Example",
  "url": "https://example.com",
  "selector": "h2.title a"
}
```

Selectors may be used internally as generic fallback heuristics, but they must not be required from the user config.

Secrets must be loaded from environment variables.

Never commit real values for:

```text
OPENROUTER_API_KEY
TUX_LETTER_SMTP_PASS
TUX_LETTER_SMTP_USER
TUX_LETTER_EMAIL_FROM
TUX_LETTER_EMAIL_TO
```

## Scraping Rules

Scraping must be polite and resilient.

Requirements:

* use timeouts;
* set a clear User-Agent;
* prefer feeds when available;
* avoid crashing when one source fails;
* normalize URLs;
* deduplicate articles;
* keep logs useful;
* avoid excessive requests;
* avoid browser automation as the default path.

The scraper should prefer this order:

1. RSS/Atom feed;
2. JSON-LD;
3. Open Graph metadata;
4. semantic HTML article blocks;
5. generic heading/link ranking.

## AI Rules

OpenRouter model fallback is required.

The configured models must be attempted in order.

The app must handle:

* rate limits;
* temporary provider failures;
* malformed model responses;
* empty responses;
* invalid structured output.

Do not log API keys.

Do not send duplicate or unnecessary requests when a previous step already failed validation.

The final newsletter should be rendered from validated structured output.

## Runtime Rules

`tux-letter` should remain lightweight.

Avoid:

* required Docker runtime for normal usage;
* required external database server;
* heavy runtime daemons;
* unnecessary dependencies;
* hidden system mutations.

Required modes:

```sh
tux-letter once
tux-letter serve
tux-letter validate-config
tux-letter sources test
tux-letter --version
```

The app must shut down gracefully on `SIGINT` and `SIGTERM`.

## Documentation Rules

The README should be simple, polished and easy to follow.

It should include:

* badges/widgets;
* short project description;
* features;
* installation;
* quick start;
* config example;
* service mode;
* development commands;
* release information;
* license.

Recommended badges:

```md
[![CI](https://github.com/kristyancarvalho/tux-letter/actions/workflows/ci.yml/badge.svg)](https://github.com/kristyancarvalho/tux-letter/actions/workflows/ci.yml)
[![Release](https://github.com/kristyancarvalho/tux-letter/actions/workflows/release.yml/badge.svg)](https://github.com/kristyancarvalho/tux-letter/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/github/license/kristyancarvalho/tux-letter)](LICENSE)
[![Go version](https://img.shields.io/github/go-mod/go-version/kristyancarvalho/tux-letter)](go.mod)
[![Latest release](https://img.shields.io/github/v/release/kristyancarvalho/tux-letter)](https://github.com/kristyancarvalho/tux-letter/releases)
[![AUR](https://img.shields.io/aur/version/tux-letter)](https://aur.archlinux.org/packages/tux-letter)
[![Active milestone](https://img.shields.io/badge/milestone-Release%203.0-2563eb)](https://github.com/kristyancarvalho/tux-letter/milestones)
```

## CI

The CI workflow must run on:

* pushes to `dev`;
* pushes to `release`;
* pull requests targeting `dev` or `release`.

Required checks:

```sh
make fmt-check
make test
make lint
make build
```

Coverage should be generated when practical.

## Releases

Releases are created from `release`.

Release flow:

1. Finish all milestone issues.
2. Validate `dev`.
3. Merge `dev` into `release`.
4. Run release validation.
5. Tag from `release`.
6. Push the tag.
7. Let GitHub Actions build release assets.
8. Validate generated binaries.
9. Update and publish AUR package.

Tag format:

```text
v3.0.0
```

The release workflow must generate downloadable binaries and checksums.

Minimum release artifacts:

```text
tux-letter_3.0.1_linux_amd64.tar.gz
tux-letter_3.0.1_linux_arm64.tar.gz
checksums.txt
```

## AUR Packaging

AUR files live in:

```text
packaging/aur/
```

Required files:

```text
packaging/aur/PKGBUILD
packaging/aur/.SRCINFO
```

Validation commands:

```sh
make aur-srcinfo
make aur-verifysource
make aur-build
```

The AUR package must:

* build from source;
* install the `tux-letter` binary;
* use the MIT license;
* not require Docker;
* not require API keys at install time;
* not run the newsletter service automatically after install.

Publishing to AUR should happen only after the GitHub release is valid.

## License

`tux-letter` is released under the MIT License.

The repository must include:

```text
LICENSE
```

## Pull Request Checklist

Before opening or merging a PR:

* the related issue exists;
* the issue is assigned to `kristyancarvalho`;
* the issue is linked to `Release 3.0`;
* commits follow the project commit style;
* `make fmt` passes;
* `make test` passes;
* `make lint` passes;
* `make build` passes;
* docs are updated when behavior changes;
* no secrets were committed;
* `./specs/` was not staged;
* generated release/AUR artifacts are only committed when expected.
