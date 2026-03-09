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

Artworks included in gh-art must be in the **public domain**. All artwork images are converted to ASCII using [jp2a](https://github.com/cslarsen/jp2a). If you'd like to add a new artwork, please ensure it meets these criteria and follow the existing artwork format.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).
