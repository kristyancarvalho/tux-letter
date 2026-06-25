# Changelog

## v3.1.0

- Generate one cohesive AI newsletter article instead of per-source cards.
- Extract full article content before summarization.
- Send a structured source bundle to OpenRouter and validate the cited output.
- Render inline numeric citations with a matching source footer.
- Add a deterministic article-style fallback digest when OpenRouter is
  unavailable.

## v3.0.1

- Prepare patch release metadata and AUR packaging for `v3.0.1`.
- Add deterministic source archive generation for release assets.
- Update AUR source packaging to consume the release source archive.

## v3.0.0

- Full rewrite as a lightweight single-binary newsletter service.
- Add URL-only sources, feed discovery, generic HTML scraping, OpenRouter model
  fallback, SMTP delivery, scheduled service mode, GitHub release assets and AUR
  packaging.
