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

This creates an MDX file in `content/blips/` with the correct frontmatter template. The id is derived from the name (`"Node.js"` becomes `node-js`). Use `--quadrant` and `--ring` to override the defaults (the first quadrant and the third ring). Edit the generated file to fill in the details.

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

- **Quadrants:** configurable via `.radar.yaml` (default: `techniques`, `platforms`, `tools`, `languages-and-frameworks`)
- **Rings:** configurable via `.radar.yaml` (default: `adopt`, `trial`, `assess`, `hold`)
- **ID format:** lowercase alphanumeric, may contain hyphens, must start with a letter or digit, and must match the filename (`my-technology.mdx`)
- **Unknown keys are rejected:** a misspelled field such as `relatedBlip:` is an error, not silently ignored
- **History:** dates in `YYYY-MM` format, newest first, and the first entry's ring must match `ring`
- **Related blips:** must name existing blips, and not the blip itself

## Validation

Validate the config and all blip content against the CUE schema without building:

```bash
make validate
```

Every problem is reported at once, with the file and line it comes from:

```
Error: validation failed: 2 problem(s):
  - content/blips/rust.mdx:4:11: quadrant: "toolz" is not one of "languages-and-frameworks", "platforms", "techniques", "tools"
  - content/blips/rust.mdx:7:15: relatedBlip: field not allowed
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

Run the E2E tests against a build served from a GitHub Pages-style sub-path (`/technology-radar/` by default, override with `PAGES_BASE=/my-repo/`):

```bash
make test-e2e-base
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
3. **CUE validation** — `.radar.yaml` and each blip's frontmatter are validated against `schema/radar.cue`, then cross-blip rules (ids, related blips, history) are checked
4. **Artifact generation** — Produces `blips-generated.json`, `radar-config.json`, and copies MDX files into the Astro content directory
5. **Astro build** — Generates the static site with Vue interactive components

## Configuration

The CLI is configured via `.radar.yaml` at the project root. Every key is optional; the values below are the defaults:

```yaml
content_dir: content/blips           # source blip MDX files
web_dir: web                         # Astro project directory
output_json: web/src/data/blips-generated.json
output_content: web/src/content/blips
version: "2025.11"                   # YYYY.MM, shown on the radar page as "Volume 2025.11"
```

The config is plain YAML; its schema and defaults live in `#Config` in `schema/radar.cue`, which is compiled into the `radar` binary. Unknown keys and invalid values are reported with their line number. You can also check the file with the [CUE CLI](https://cuelang.org/docs/introduction/installation/):

```bash
cue vet -c -d '#Config' schema/radar.cue .radar.yaml
```

Command-line flags override the file: `--config <file>` (a missing file is an error when passed explicitly), `--content-dir`, and `--web-dir`.

### Customizing Quadrants and Rings

You can customize the quadrant names and ring (blip state) names by adding `quadrants` and `rings` sections to `.radar.yaml`. Both sections are optional and default to the standard values shown below.

```yaml
quadrants:
  - slug: techniques
    label: "Techniques"
  - slug: platforms
    label: "Platforms"
  - slug: tools
    label: "Tools"
  - slug: languages-and-frameworks
    label: "Languages & Frameworks"

rings:
  - slug: adopt
    label: "Adopt"
  - slug: trial
    label: "Trial"
  - slug: assess
    label: "Assess"
  - slug: hold
    label: "Hold"
```

Each entry has two fields:

- **slug** — URL-safe identifier (lowercase letters, digits, hyphens) used in blip frontmatter, URL paths, and CUE schema validation. Blip content files reference this value in their `quadrant` and `ring` fields.
- **label** — Display name shown in the radar UI, search filters, breadcrumbs, and detail pages.

**Order matters.** Quadrant order maps to angular position on the radar (index 0 = top-right, 1 = bottom-right, 2 = bottom-left, 3 = top-left). Ring order maps to radial position (index 0 = innermost, 3 = outermost). Colors are assigned by position.

**Constraints:**

- Exactly 4 quadrants and 4 rings are required.
- Slugs must be unique and non-empty.
- Labels must be non-empty.

**Changing labels vs. slugs:**

- Changing a **label** (e.g., `"Tools"` to `"Developer Tools"`) only affects the display. No content files need updating — just rebuild the site with `make build-web`.
- Changing a **slug** (e.g., `"tools"` to `"dev-tools"`) requires updating the `quadrant` or `ring` field in every affected blip MDX file. The CUE validation step will catch any mismatches.

**Updating FAQ content:** The FAQ page (`web/src/data/faq.ts`) contains prose descriptions of the rings and quadrants. If you customize the names, update the FAQ text manually to match.

## Deploying to GitHub Pages

The site is fully static, so it can be hosted on GitHub Pages. `.github/workflows/pages.yml` builds, tests, and deploys it on every push to `main` (and builds and tests pull requests without deploying).

1. Push the repository to GitHub.
2. In the repository's **Settings → Pages**, set **Source** to **GitHub Actions**. This is a one-time step.
3. Push to `main`, or run the workflow by hand from the **Actions** tab.

The site is then available at `https://<owner>.github.io/<repo>/radar/`.

**Base path.** A project site is served from `/<repo>/` rather than `/`. The workflow reads the base path from your Pages configuration and passes it to the build as `BASE_PATH` (and the origin as `SITE_URL`). Every internal link goes through `withBase()` in `web/src/utils/url.ts`, so links work under either. With a custom domain the base path is empty and the site is served from `/`. When adding links in pages, components, or data files, always use `withBase('/radar/...')` (or `blipHref(blip)`) rather than a hard-coded `/radar/...` path; `make test-e2e-base` fails if a link is missing the base.

To build for a sub-path locally:

```bash
BASE_PATH=/technology-radar/ make build-web
```

**Repository link.** The GitHub icon in the header and the FAQ/Contribute links point at `repoUrl` in `web/src/data/site.ts`; update it if you fork the project.
