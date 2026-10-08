# Technology Radar

An interactive, Thoughtworks-inspired technology radar for visualizing your organization's technology landscape across adoption stages.

## Features

- **Interactive SVG radar** with smooth quadrant zoom via viewBox transitions
- **Synchronized hover** between radar blips and sidebar lists
- **Search and filter** across all technology entries
- **Dark mode** with theme toggle
- **Configurable quadrants and rings** with custom labels
- **Strict schema validation** of config and content using CUE at build time
- **Static site generation**; blip pages ship no JavaScript
- **Deployable to GitHub Pages**, including project sites served from a sub-path
- **Version history tracking** for ring movement over time
- **Scaffolding CLI** for creating new blip entries

## Tech Stack

| Layer | Technology |
|-------|-----------|
| CLI | Go (Cobra) |
| Config & Schema Validation | YAML config, CUE schema |
| Static Site | Astro 7 |
| Interactive UI | Vue 3 |
| Layout Math | D3 (d3-shape) |
| Content | MDX |
| E2E Tests | Playwright |
| Package Manager | bun |

## Quick Start

**Prerequisites:** Go 1.25+, [bun](https://bun.sh), GNU Make

```bash
git clone https://github.com/infrashift/technology-radar.git
cd technology-radar
make setup
make dev
```

See [GETTING-STARTED.md](GETTING-STARTED.md) for detailed setup and usage instructions.

## Project Structure

```
├── cmd/radar/          # CLI entrypoint
├── internal/           # Go packages (build, cli, config, validate)
├── schema/             # CUE schema for .radar.yaml and blip frontmatter (embedded in the CLI)
├── content/blips/      # MDX blip content files
├── web/                # Astro web project
│   ├── src/
│   │   ├── components/ # Vue + Astro components (radar/ holds the shared SVG pieces)
│   │   ├── config/     # Radar geometry and labels derived from the generated config
│   │   ├── lib/        # Build-time data loading
│   │   ├── pages/      # Astro page routes
│   │   ├── types/      # TypeScript type definitions
│   │   └── utils/      # Layout math, URL helpers
│   └── e2e/            # Playwright tests
├── .github/workflows/  # GitHub Pages build, test, and deploy
└── Makefile            # Build orchestration
```

## Available Commands

| Command | Description |
|---------|-------------|
| `make setup` | Install all dependencies |
| `make dev` | Start development server |
| `make build-web` | Build the Astro site |
| `make preview` | Preview the built site |
| `make validate` | Validate config and content against the CUE schema |
| `make lint` | Run all linters (Go vet + Astro check) |
| `make test` | Run all tests (Go + Playwright) |
| `make test-go` | Run Go unit tests |
| `make test-e2e` | Run Playwright end-to-end tests |
| `make test-e2e-base` | Build and run the e2e tests under the GitHub Pages sub-path |
| `make scaffold NAME="..."` | Scaffold a new blip |
| `make clean` | Remove all build artifacts |

## Deployment

The site is fully static and deploys to GitHub Pages with the included workflow. See [Deploying to GitHub Pages](GETTING-STARTED.md#deploying-to-github-pages).
