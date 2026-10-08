CLI_NAME   := radar
BIN        := bin/$(CLI_NAME)
WEB_DIR    := web
CONTENT    := content/blips
GO_SRC     := $(shell find cmd internal schema \( -name '*.go' -o -name '*.cue' \) 2>/dev/null)
MDX_SRC    := $(shell find $(CONTENT) -name '*.mdx' 2>/dev/null)
JSON_OUT   := $(WEB_DIR)/src/data/blips-generated.json
CONFIG_OUT := $(WEB_DIR)/src/data/radar-config.json

.PHONY: all setup dev clean lint validate help test test-go test-e2e test-e2e-base build-web preview scaffold

# Sub-path used to test a GitHub Pages project-site build (https://<owner>.github.io/<repo>/).
PAGES_BASE ?= /technology-radar/

all: build-web ## Build everything (CLI, validate, generate, Astro)

help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

setup: ## Install all dependencies
	@echo "Installing Go dependencies..."
	go mod tidy
	@echo "Installing web dependencies..."
	cd $(WEB_DIR) && bun install
	@echo "Installing Playwright browsers..."
	cd $(WEB_DIR) && bunx playwright install chromium

# Build the Go CLI binary (rebuild when Go source or the embedded CUE schema changes)
$(BIN): $(GO_SRC) go.mod go.sum
	@echo "Building Go CLI..."
	@mkdir -p bin
	go build -o $@ ./cmd/radar

# Generate artifacts (rebuild when MDX, config, or CLI changes). `radar build`
# writes both files; the second rule only reruns it if the config JSON is missing.
$(JSON_OUT): $(BIN) $(MDX_SRC) $(wildcard .radar.yaml)
	@echo "Running data pipeline..."
	./$(BIN) build

$(CONFIG_OUT): $(JSON_OUT)
	@test -f $@ || ./$(BIN) build

validate: $(BIN) ## Validate config and content against the CUE schema
	./$(BIN) validate

lint: $(JSON_OUT) $(CONFIG_OUT) ## Run all linters (Go vet + Astro check)
	@echo "Running Go vet..."
	go vet ./...
	@echo "Running Astro check..."
	cd $(WEB_DIR) && bun run check

test-go: ## Run Go unit tests
	go test ./...

test-e2e: build-web ## Run Playwright end-to-end tests
	cd $(WEB_DIR) && bun run test:e2e

test-e2e-base: $(JSON_OUT) $(CONFIG_OUT) ## Build and run e2e tests under the GitHub Pages sub-path (PAGES_BASE)
	cd $(WEB_DIR) && BASE_PATH=$(PAGES_BASE) bun run build
	cd $(WEB_DIR) && BASE_PATH=$(PAGES_BASE) bun run test:e2e

test: test-go test-e2e ## Run all tests (Go + Playwright)

build-web: $(JSON_OUT) $(CONFIG_OUT) ## Build the Astro site (set BASE_PATH/SITE_URL for a sub-path deploy)
	@echo "Building Astro site..."
	cd $(WEB_DIR) && bun run build

dev: $(JSON_OUT) $(CONFIG_OUT) ## Start development server
	cd $(WEB_DIR) && bun run dev

preview: build-web ## Preview the built site
	cd $(WEB_DIR) && bun run preview

scaffold: $(BIN) ## Scaffold a new blip (usage: make scaffold NAME="My Blip")
	./$(BIN) scaffold "$(NAME)"

clean: ## Remove all build artifacts
	rm -rf bin/
	rm -rf $(WEB_DIR)/dist
	rm -rf $(WEB_DIR)/src/content/blips
	rm -f $(JSON_OUT) $(CONFIG_OUT)
	rm -rf $(WEB_DIR)/.astro
