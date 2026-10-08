package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeConfig(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".radar.yaml")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadMissingFileUsesDefaults(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "missing.yaml"), false)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.ContentDir != "content/blips" {
		t.Errorf("expected default content_dir, got %q", cfg.ContentDir)
	}
	if got := cfg.QuadrantSlugs(); len(got) != 4 || got[0] != "techniques" {
		t.Errorf("unexpected default quadrants: %v", got)
	}
	if got := cfg.RingSlugs(); len(got) != 4 || got[0] != "adopt" {
		t.Errorf("unexpected default rings: %v", got)
	}
}

func TestLoadMissingRequiredFile(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "missing.yaml"), true); err == nil {
		t.Fatal("expected error for missing explicit config file")
	}
}

func TestLoadEmptyFileUsesDefaults(t *testing.T) {
	cfg, err := Load(writeConfig(t, ""), true)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Version != "2025.11" {
		t.Errorf("expected default version, got %q", cfg.Version)
	}
}

func TestLoadOverrides(t *testing.T) {
	cfg, err := Load(writeConfig(t, `
content_dir: blips
version: "2026.05"
rings:
  - {slug: use, label: Use}
  - {slug: try, label: Try}
  - {slug: watch, label: Watch}
  - {slug: avoid, label: Avoid}
`), true)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.ContentDir != "blips" || cfg.Version != "2026.05" {
		t.Errorf("overrides not applied: %+v", cfg)
	}
	if cfg.OutputJSON != "web/src/data/blips-generated.json" {
		t.Errorf("expected default output_json, got %q", cfg.OutputJSON)
	}
	if got := cfg.RingSlugs(); got[3] != "avoid" {
		t.Errorf("unexpected rings: %v", got)
	}
	if got := cfg.QuadrantSlugs(); got[0] != "techniques" {
		t.Errorf("expected default quadrants, got %v", got)
	}
}

func TestLoadInvalid(t *testing.T) {
	tests := map[string]struct {
		body string
		want string
	}{
		"unknown key": {
			body: "contnet_dir: blips\n",
			want: ".radar.yaml:1:1: contnet_dir: field not allowed",
		},
		"wrong count": {
			body: "quadrants:\n  - {slug: a, label: A}\n  - {slug: b, label: B}\n  - {slug: c, label: C}\n",
			want: "MinItems(4)",
		},
		"duplicate slug": {
			body: "rings:\n  - {slug: a, label: A}\n  - {slug: a, label: B}\n  - {slug: c, label: C}\n  - {slug: d, label: D}\n",
			want: "UniqueItems",
		},
		"empty label": {
			body: "rings:\n  - {slug: a, label: \"\"}\n  - {slug: b, label: B}\n  - {slug: c, label: C}\n  - {slug: d, label: D}\n",
			want: `.radar.yaml:2:22: rings.0.label: invalid value ""`,
		},
		"bad slug": {
			body: "rings:\n  - {slug: Adopt, label: A}\n  - {slug: b, label: B}\n  - {slug: c, label: C}\n  - {slug: d, label: D}\n",
			want: "rings.0.slug",
		},
		"bad version": {
			body: "version: latest\n",
			want: "version",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Load(writeConfig(t, tc.body), true)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("expected error containing %q, got:\n%v", tc.want, err)
			}
		})
	}
}
