# Contributing to gh-art

Welcome! We're glad you're interested in contributing to gh-art. This document provides guidelines and instructions for contributing.

## Reporting Bugs

Found a bug? Please [open an issue](../../issues/new?template=bug_report.md) using the bug report template. Include as much detail as possible to help us reproduce the problem.

## Suggesting Features

Have an idea for a new feature? [Open a feature request](../../issues/new?template=feature_request.md) and describe your proposal.

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/maxbeizer/gh-art.git
   cd gh-art
   ```

2. Make sure you have **Go 1.21+** installed.

3. Build the project:
   ```bash
   make build
   ```

4. Run the tests:
   ```bash
   make test
   ```

5. Run the linter:
   ```bash
   make lint
   ```

## Submitting Pull Requests

1. Fork the repository and create your branch from `main`.
2. Make your changes.
3. Run the full CI suite to make sure everything passes:
   ```bash
   make ci
   ```
4. Open a pull request with a clear description of your changes.

## Code Style

- Run `make fmt` before committing to ensure consistent formatting.
- Follow existing patterns and conventions in the codebase.

## Adding New Artworks

Artworks included in gh-art must be in the **public domain**. All artwork images
are converted to ASCII using [jp2a](https://github.com/cslarsen/jp2a).

### Artwork file format

Each artwork lives in its own `.txt` file inside the `artworks/` directory. The
file must start with YAML-style frontmatter delimited by `---` lines:

```text
---
name: my-artwork
title: My Artwork
artist: Some Artist
year: 1900
url: https://en.wikipedia.org/wiki/My_Artwork
---
<ASCII art content here>
```

**Required frontmatter fields:**

| Field    | Description                                      |
|----------|--------------------------------------------------|
| `name`   | Unique kebab-case identifier (used in `gh art show <name>`) |
| `title`  | Human-readable title of the artwork              |
| `artist` | Name of the original artist                      |
| `year`   | Year the artwork was created/completed           |
| `url`    | Wikipedia or reference URL                       |

### Guidelines

- The ASCII art content follows immediately after the closing `---` line.
- Keep dimensions reasonable for a typical terminal (~80–120 columns wide).
- Use [jp2a](https://github.com/cslarsen/jp2a) or a similar tool to convert images.
- Run `make build && ./bin/gh-art show my-artwork` to verify your artwork renders correctly.

### Custom artworks (local only)

Users can also add personal artworks without modifying the repository by placing
`.txt` files in `~/.config/gh-art/artworks/`, or by using `gh art import <file>`.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).
