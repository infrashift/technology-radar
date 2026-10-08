package cli

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/infrashift/technology-radar/internal/build"
	"github.com/infrashift/technology-radar/internal/config"
	"github.com/infrashift/technology-radar/internal/validate"
)

// scaffold runs "radar scaffold" in a fresh command tree against dir. There
// is no .radar.yaml in the test's working directory, so defaults apply.
func scaffold(t *testing.T, dir string, args ...string) error {
	t.Helper()
	cmd := newRootCmd()
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs(append([]string{"scaffold", "--content-dir", dir}, args...))
	return cmd.Execute()
}

// readScaffolded validates the scaffolded file against the schema and returns it.
func readScaffolded(t *testing.T, dir, name string) (build.Blip, string) {
	t.Helper()
	files, err := build.ReadBlipFiles(dir)
	if err != nil {
		t.Fatalf("reading scaffolded file: %v", err)
	}
	cfg, err := config.Load(filepath.Join(dir, "none.yaml"), false)
	if err != nil {
		t.Fatal(err)
	}
	blips, err := validate.Blips(cfg, files)
	if err != nil {
		t.Fatalf("scaffolded file does not validate: %v", err)
	}
	body, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	return blips[0], string(body)
}

func TestScaffoldBasic(t *testing.T) {
	dir := t.TempDir()
	if err := scaffold(t, dir, "My Test Blip"); err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	blip, body := readScaffolded(t, dir, "my-test-blip.mdx")
	if blip.ID != "my-test-blip" || blip.Title != "My Test Blip" {
		t.Errorf("unexpected id/title: %+v", blip)
	}
	if blip.Quadrant != "techniques" || blip.Ring != "assess" || !blip.IsNew {
		t.Errorf("unexpected defaults: %+v", blip)
	}
	if len(blip.History) != 1 || blip.History[0].Ring != "assess" {
		t.Errorf("unexpected history: %+v", blip.History)
	}
	if strings.Contains(body, "ciIdentifier") {
		t.Error("ciIdentifier should not appear when --ci-identifier is not provided")
	}
}

func TestScaffoldWithCIIdentifier(t *testing.T) {
	dir := t.TempDir()
	uuid := "01944b6e-40c0-7000-8000-000000000000"
	if err := scaffold(t, dir, "CI Blip", "--ci-identifier", uuid); err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	blip, _ := readScaffolded(t, dir, "ci-blip.mdx")
	if blip.CIIdentifier != uuid {
		t.Errorf("expected ciIdentifier %q, got %q", uuid, blip.CIIdentifier)
	}
}

func TestScaffoldWithQuadrantAndRing(t *testing.T) {
	dir := t.TempDir()
	if err := scaffold(t, dir, "Custom Blip", "--quadrant", "tools", "--ring", "adopt"); err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	blip, _ := readScaffolded(t, dir, "custom-blip.mdx")
	if blip.Quadrant != "tools" || blip.Ring != "adopt" {
		t.Errorf("expected tools/adopt, got %s/%s", blip.Quadrant, blip.Ring)
	}
}

func TestScaffoldPunctuationAndQuotes(t *testing.T) {
	dir := t.TempDir()
	if err := scaffold(t, dir, `Node.js "Edge": Runtime`); err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	blip, _ := readScaffolded(t, dir, "node-js-edge-runtime.mdx")
	if blip.Title != `Node.js "Edge": Runtime` {
		t.Errorf("title not preserved: %q", blip.Title)
	}
}

func TestScaffoldInvalidRing(t *testing.T) {
	err := scaffold(t, t.TempDir(), "Blip", "--ring", "nope")
	if err == nil || !strings.Contains(err.Error(), `invalid ring "nope"`) {
		t.Errorf("expected invalid ring error, got %v", err)
	}
}

func TestScaffoldNoSlugChars(t *testing.T) {
	if err := scaffold(t, t.TempDir(), "!!!"); err == nil {
		t.Error("expected error for a name with no letters or digits")
	}
}

func TestScaffoldFileAlreadyExists(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "existing.mdx"), []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := scaffold(t, dir, "existing")
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", err)
	}
}

func TestSlugify(t *testing.T) {
	tests := map[string]string{
		"My Test Blip":     "my-test-blip",
		"Node.js":          "node-js",
		"  C++ / CLI  ":    "c-cli",
		"AGENTS.md":        "agents-md",
		"Already-a-slug":   "already-a-slug",
		"Rust 2024 (beta)": "rust-2024-beta",
	}
	for in, want := range tests {
		if got := slugify(in); got != want {
			t.Errorf("slugify(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestExplicitConfigMustExist(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"validate", "--config", filepath.Join(t.TempDir(), "missing.yaml")})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing explicit config file")
	}
}
