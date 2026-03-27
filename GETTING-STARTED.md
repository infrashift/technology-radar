# Getting Started

## Prerequisites

- **Go 1.25+** — [Install Go](https://go.dev/dl/)
- **bun** — [Install bun](https://bun.sh)
- **GNU Make** — typically pre-installed on Linux/macOS

## Installation

```bash
git clone https://github.com/infrashift/technology-radar.git
cd technology-radar
make setup
```

This installs Go dependencies, web dependencies (via bun), and Playwright browsers.

## Development Workflow

Start the dev server with hot reload:

```bash
make dev
```

Build the full site:

```bash
make build-web
```

Preview the production build locally:

```bash
make preview
```

## Adding a New Blip

Scaffold a new blip entry:

```bash
make scaffold NAME="My Technology"
```

This creates an MDX file in `content/blips/` with the correct frontmatter template. Edit the generated file to fill in the details.

### Blip Frontmatter Schema

```yaml
---
id: "my-technology"          # lowercase alphanumeric with hyphens
title: "My Technology"       # display name
quadrant: "tools"            # techniques | platforms | tools | languages-and-frameworks
ring: "assess"               # adopt | trial | assess | hold
isNew: true                  # whether this is a new entry (default: false)
history:                     # optional: track ring changes over time
  - date: "2025-11"
    ring: "assess"
relatedBlips:                # optional: IDs of related blips
  - "some-other-blip"
---

Write your description here in Markdown.
```

### Schema Constraints

- **Quadrants:** `techniques`, `platforms`, `tools`, `languages-and-frameworks`
- **Rings:** `adopt`, `trial`, `assess`, `hold`
- **ID format:** lowercase alphanumeric, may contain hyphens, must start with a letter or digit
- **History dates:** `YYYY-MM` format

## Validation

Validate all blip content against the CUE schema without building:

```bash
make validate
```

## Testing

Run all tests (Go unit tests + Playwright E2E):

```bash
make test
```

Run only Go unit tests:

```bash
make test-go
```

Run only Playwright E2E tests (builds the site first):

```bash
make test-e2e
```

## Linting

```bash
make lint
```

This runs `go vet` on Go source and `astro check` on the web project.

## Build Pipeline

The build process follows this pipeline:

1. **Go CLI build** — Compiles the `radar` CLI binary from `cmd/radar/`
2. **MDX parsing** — The CLI reads all MDX files from `content/blips/` and extracts YAML frontmatter
3. **CUE validation** — Extracted data is validated against `schema/radar.cue`
4. **Artifact generation** — Produces `blips-generated.json` and copies MDX files into the Astro content directory
5. **Astro build** — Generates the static site with Vue interactive components

## Configuration

The CLI is configured via `.radar.yaml` at the project root:

```yaml
content_dir: content/blips           # source blip MDX files
schema_path: schema/radar.cue        # CUE schema for validation
web_dir: web                         # Astro project directory
output_json: web/src/data/blips-generated.json
output_content: web/src/content/blips
version: "2025.11"
```
