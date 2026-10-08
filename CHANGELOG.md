# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Calendar Versioning](https://calver.org/) (`YYYY0M0D.MICRO`).

## [Unreleased]

### Added

- Configurable quadrant names and ring (blip state) names via `.radar.yaml`
- `radar-config.json` generated artifact to pass radar configuration (quadrants, rings, version) to the frontend
- `#Config` CUE schema for `.radar.yaml`: defaults, exactly 4 quadrants and 4 rings, unique slugs, non-empty labels, `YYYY.MM` version, and rejection of unknown keys. Can also be checked with `cue vet -c -d '#Config' schema/radar.cue .radar.yaml`
- Content validation rules: unknown frontmatter keys, blip id must match its filename, `relatedBlips` must reference existing blips, history must be newest first and agree with the current ring
- Validation errors point at the file, line, and column, and all problems are reported in one run
- GitHub Pages deployment workflow (`.github/workflows/pages.yml`) that builds, tests, and deploys under the repository's base path
- `make test-e2e-base` to test a build served from a sub-path
- Radar zoom state in the URL hash (`/radar#tools`), so breadcrumbs and links open the radar zoomed to a quadrant
- Radar volume (`version`) shown on the radar page
- Keyboard access to quadrant zoom controls, and focus states for blips
- 404 page and favicon
- E2E tests for viewport fit, mobile layout, zoom and deep links, tooltip stacking, base-path links, and blip layout

### Changed

- The radar page fits the viewport on desktop; the blip list scrolls on its own
- Config is loaded and validated through CUE instead of Viper; the schema is embedded in the `radar` binary
- `schema_path` config key and `--schema-path` flag removed (the schema is embedded)
- `RADAR_*` environment variable overrides removed; use `--config`, `--content-dir`, or `--web-dir`
- An explicit `--config` file that does not exist is now an error
- `internal/cuegen` replaced by `internal/validate`; quadrant and ring enums come from CUE unification instead of rewriting the schema text
- Upgraded Astro 5.7 → 7.3, `@astrojs/vue` 7, `@astrojs/mdx` 8, Playwright 1.64
- Blip layout is computed once at build time and passed to components
- Blip detail pages render the quadrant radar as static HTML (no client JavaScript)
- Fonts are self-hosted (`@fontsource-variable/inter`, `@fontsource-variable/jetbrains-mono`) instead of loaded from Google Fonts
- Shared radar SVG components (`RadarBackdrop`, `RadarBlip`, `RadarLabel`) replace duplicated code in the full and quadrant radars
- `radar scaffold` derives ids by slugifying the name and writes frontmatter with a YAML encoder
- Makefile tracks `.radar.yaml` and the embedded schema as build dependencies

### Fixed

- Blips near the ±180° edge of the left-hand quadrants could jump to the opposite edge of their quadrant
- The quadrant heading and number badges in the blip page sidebar always rendered gray
- Hover tooltips could be hidden behind blips drawn later
- The sticky quadrant radar on blip pages slid under the sticky header
- Sidebar blip lists were ordered by filename instead of blip number
- Long quadrant titles overflowed their label background
- Deleted or renamed blips stayed published, because stale copies were never removed from the Astro content directory
- `radar scaffold` produced ids that failed validation for names with punctuation (`Node.js`), and invalid YAML for titles containing quotes
- Frontmatter parsing broke on `---` inside a value
- `make help` omitted targets with digits in their names

### Removed

- Unused `d3-scale` and `@vitejs/plugin-vue` dependencies

### Security

- The Astro upgrade and lockfile refresh resolve all 68 `bun audit` advisories (1 critical, including Astro image-optimization RCE GHSA-26w7-cxv4-gfx2)

## [20251101.1] - 2025-11-01

### Added

- Interactive SVG radar visualization with four quadrants and four rings
- Smooth quadrant zoom via SVG viewBox CSS transitions
- Synchronized hover state between radar blips and sidebar lists
- Blip search and filter component
- Dark mode with theme toggle
- Go CLI (`radar`) with `build`, `validate`, and `scaffold` commands
- CUE schema validation for blip content at build time
- MDX content format with structured YAML frontmatter
- Astro static site generation with Vue 3 interactive islands
- Playwright end-to-end test suite
- Makefile-based build orchestration
- 12 initial technology blips across all quadrants and rings
