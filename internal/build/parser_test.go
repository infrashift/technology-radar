package build

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFrontmatterBasic(t *testing.T) {
	content := []byte(`---
id: "test-blip"
title: "Test Blip"
quadrant: "tools"
ring: "adopt"
isNew: true
history:
  - date: "2025-11"
    ring: "adopt"
---

Some content here.
`)

	blip, err := parseFrontmatter(content, "test-blip.mdx")
	if err != nil {
		t.Fatalf("parseFrontmatter failed: %v", err)
	}

	if blip.ID != "test-blip" {
		t.Errorf("expected id 'test-blip', got %q", blip.ID)
	}
	if blip.Title != "Test Blip" {
		t.Errorf("expected title 'Test Blip', got %q", blip.Title)
	}
	if blip.Quadrant != "tools" {
		t.Errorf("expected quadrant 'tools', got %q", blip.Quadrant)
	}
	if blip.Ring != "adopt" {
		t.Errorf("expected ring 'adopt', got %q", blip.Ring)
	}
	if !blip.IsNew {
		t.Error("expected isNew to be true")
	}
	if len(blip.History) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(blip.History))
	}
	if blip.History[0].Date != "2025-11" {
		t.Errorf("expected history date '2025-11', got %q", blip.History[0].Date)
	}
}

func TestParseFrontmatterWithCIIdentifier(t *testing.T) {
	content := []byte(`---
id: "ci-blip"
title: "CI Blip"
quadrant: "platforms"
ring: "trial"
isNew: false
ciIdentifier: "01944b6e-40c0-7000-8000-000000000000"
---

Content.
`)

	blip, err := parseFrontmatter(content, "ci-blip.mdx")
	if err != nil {
		t.Fatalf("parseFrontmatter failed: %v", err)
	}

	if blip.CIIdentifier != "01944b6e-40c0-7000-8000-000000000000" {
		t.Errorf("expected ciIdentifier, got %q", blip.CIIdentifier)
	}
}

func TestParseFrontmatterWithoutCIIdentifier(t *testing.T) {
	content := []byte(`---
id: "no-ci"
title: "No CI"
quadrant: "tools"
ring: "hold"
isNew: false
---

Content.
`)

	blip, err := parseFrontmatter(content, "no-ci.mdx")
	if err != nil {
		t.Fatalf("parseFrontmatter failed: %v", err)
	}

	if blip.CIIdentifier != "" {
		t.Errorf("expected empty ciIdentifier, got %q", blip.CIIdentifier)
	}
}

func TestParseFrontmatterMissingID(t *testing.T) {
	content := []byte(`---
title: "No ID"
quadrant: "tools"
ring: "hold"
---

Content.
`)

	_, err := parseFrontmatter(content, "no-id.mdx")
	if err == nil {
		t.Fatal("expected error for missing id")
	}
}

func TestParseFrontmatterMissingDelimiters(t *testing.T) {
	content := []byte(`no frontmatter here`)

	_, err := parseFrontmatter(content, "bad.mdx")
	if err == nil {
		t.Fatal("expected error for missing delimiters")
	}
}

func TestParseAllBlips(t *testing.T) {
	dir := t.TempDir()

	files := map[string]string{
		"alpha.mdx": `---
id: "alpha"
title: "Alpha"
quadrant: "tools"
ring: "adopt"
isNew: false
---

Alpha content.
`,
		"beta.mdx": `---
id: "beta"
title: "Beta"
quadrant: "platforms"
ring: "trial"
isNew: true
ciIdentifier: "01944b6e-40c0-7000-8000-000000000001"
---

Beta content.
`,
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	blips, rawFiles, err := ParseAllBlips(dir)
	if err != nil {
		t.Fatalf("ParseAllBlips failed: %v", err)
	}

	if len(blips) != 2 {
		t.Fatalf("expected 2 blips, got %d", len(blips))
	}
	if len(rawFiles) != 2 {
		t.Fatalf("expected 2 raw files, got %d", len(rawFiles))
	}

	// Find the beta blip and check ciIdentifier
	for _, b := range blips {
		if b.ID == "beta" {
			if b.CIIdentifier != "01944b6e-40c0-7000-8000-000000000001" {
				t.Errorf("expected ciIdentifier on beta, got %q", b.CIIdentifier)
			}
		}
	}
}

func TestParseAllBlipsEmptyDir(t *testing.T) {
	dir := t.TempDir()

	_, _, err := ParseAllBlips(dir)
	if err == nil {
		t.Fatal("expected error for empty directory")
	}
}
