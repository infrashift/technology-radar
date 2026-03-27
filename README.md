# Technology Radar

An interactive, Thoughtworks-inspired technology radar for visualizing your organization's technology landscape across adoption stages.

## Features

- **Interactive SVG radar** with smooth quadrant zoom via viewBox transitions
- **Synchronized hover** between radar blips and sidebar lists
- **Search and filter** across all technology entries
- **Dark mode** with theme toggle
- **Strict schema validation** using CUE at build time
- **Static site generation** with zero runtime overhead
- **Version history tracking** for ring movement over time
- **Scaffolding CLI** for creating new blip entries

## Tech Stack

| Layer | Technology |
|-------|-----------|
| CLI | Go (Cobra, Viper) |
| Schema Validation | CUE |
| Static Site | Astro |
| Interactive UI | Vue 3 |
| Layout Math | D3 (d3-scale, d3-shape) |
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
├── internal/           # Go packages (build, cli, config, cuegen)
├── schema/             # CUE schema definition
├── content/blips/      # MDX blip content files
├── web/                # Astro web project
│   ├── src/
│   │   ├── components/ # Vue + Astro components
│   │   ├── composables/# Vue composables (layout, animation, hover sync)
│   │   ├── pages/      # Astro page routes
│   │   └── types/      # TypeScript type definitions
│   └── e2e/            # Playwright tests
└── Makefile            # Build orchestration
```

## Available Commands

| Command | Description |
|---------|-------------|
| `make setup` | Install all dependencies |
| `make dev` | Start development server |
| `make build-web` | Build the Astro site |
| `make preview` | Preview the built site |
| `make validate` | Validate content against CUE schema |
| `make lint` | Run all linters (Go vet + Astro check) |
| `make test` | Run all tests (Go + Playwright) |
| `make test-go` | Run Go unit tests |
| `make test-e2e` | Run Playwright end-to-end tests |
| `make scaffold NAME="..."` | Scaffold a new blip |
| `make clean` | Remove all build artifacts |
