package build

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/infrashift/technology-radar/internal/config"
)

func TestWriteJSON(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "sub", "blips.json")

	blips := []Blip{
		{
			ID:       "test",
			Title:    "Test",
			Quadrant: "tools",
			Ring:     "adopt",
			IsNew:    false,
		},
		{
			ID:           "with-ci",
			Title:        "With CI",
			Quadrant:     "platforms",
			Ring:         "trial",
			IsNew:        true,
			CIIdentifier: "01944b6e-40c0-7000-8000-000000000000",
		},
	}

	if err := WriteJSON(outPath, blips); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}

	var parsed []Blip
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshaling output: %v", err)
	}

	if len(parsed) != 2 {
		t.Fatalf("expected 2 blips, got %d", len(parsed))
	}

	if parsed[1].CIIdentifier != "01944b6e-40c0-7000-8000-000000000000" {
		t.Errorf("expected ciIdentifier preserved in JSON, got %q", parsed[1].CIIdentifier)
	}
}

func TestWriteJSONOmitsEmptyCIIdentifier(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "blips.json")

	blips := []Blip{
		{
			ID:       "no-ci",
			Title:    "No CI",
			Quadrant: "tools",
			Ring:     "hold",
			IsNew:    false,
		},
	}

	if err := WriteJSON(outPath, blips); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}

	// Verify ciIdentifier is omitted from JSON when empty
	var raw []map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshaling: %v", err)
	}

	if _, exists := raw[0]["ciIdentifier"]; exists {
		t.Error("ciIdentifier should be omitted from JSON when empty")
	}
}

func TestWriteRadarConfig(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "radar-config.json")

	cfg := config.Config{
		Version: "2025.11",
		Quadrants: []config.Entry{
			{Slug: "techniques", Label: "Techniques"},
			{Slug: "platforms", Label: "Platforms"},
			{Slug: "tools", Label: "Tools"},
			{Slug: "languages-and-frameworks", Label: "Languages & Frameworks"},
		},
		Rings: []config.Entry{
			{Slug: "adopt", Label: "Adopt"},
			{Slug: "trial", Label: "Trial"},
			{Slug: "assess", Label: "Assess"},
			{Slug: "hold", Label: "Hold"},
		},
	}

	if err := WriteRadarConfig(outPath, cfg); err != nil {
		t.Fatalf("WriteRadarConfig failed: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}
	if !strings.Contains(string(data), `"Languages & Frameworks"`) {
		t.Errorf("expected unescaped ampersand in output, got:\n%s", data)
	}

	var parsed radarConfig
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshaling output: %v", err)
	}

	if parsed.Version != "2025.11" {
		t.Errorf("expected version '2025.11', got %q", parsed.Version)
	}
	if len(parsed.Quadrants) != 4 {
		t.Fatalf("expected 4 quadrants, got %d", len(parsed.Quadrants))
	}
	if parsed.Quadrants[0].Label != "Techniques" {
		t.Errorf("expected first quadrant label 'Techniques', got %q", parsed.Quadrants[0].Label)
	}
	if len(parsed.Rings) != 4 {
		t.Fatalf("expected 4 rings, got %d", len(parsed.Rings))
	}
	if parsed.Rings[3].Slug != "hold" {
		t.Errorf("expected last ring slug 'hold', got %q", parsed.Rings[3].Slug)
	}
}

func TestCopyMDXFiles(t *testing.T) {
	dir := t.TempDir()
	destDir := filepath.Join(dir, "output")

	files := []RawBlipFile{
		{Filename: "one.mdx", Content: []byte("content one")},
		{Filename: "two.mdx", Content: []byte("content two")},
	}

	if err := CopyMDXFiles(files, destDir); err != nil {
		t.Fatalf("CopyMDXFiles failed: %v", err)
	}

	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(destDir, f.Filename))
		if err != nil {
			t.Fatalf("reading %s: %v", f.Filename, err)
		}
		if string(data) != string(f.Content) {
			t.Errorf("content mismatch for %s", f.Filename)
		}
	}
}

func TestCopyMDXFilesRemovesStale(t *testing.T) {
	destDir := t.TempDir()
	for _, name := range []string{"stale.mdx", "keep.txt"} {
		if err := os.WriteFile(filepath.Join(destDir, name), []byte("old"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := CopyMDXFiles([]RawBlipFile{{Filename: "one.mdx", Content: []byte("new")}}, destDir); err != nil {
		t.Fatalf("CopyMDXFiles failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, "stale.mdx")); !os.IsNotExist(err) {
		t.Error("expected stale.mdx to be removed")
	}
	if _, err := os.Stat(filepath.Join(destDir, "keep.txt")); err != nil {
		t.Error("expected non-MDX files to be left alone")
	}
}
