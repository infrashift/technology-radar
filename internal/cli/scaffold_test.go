package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScaffoldBasic(t *testing.T) {
	dir := t.TempDir()

	cmd := rootCmd
	cmd.SetArgs([]string{
		"scaffold", "My Test Blip",
		"--content-dir", dir,
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	filename := filepath.Join(dir, "my-test-blip.mdx")
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading scaffolded file: %v", err)
	}

	body := string(content)

	assertions := map[string]string{
		"id":       `id: "my-test-blip"`,
		"title":    `title: "My Test Blip"`,
		"quadrant": `quadrant: "techniques"`,
		"ring":     `ring: "assess"`,
		"isNew":    `isNew: true`,
	}
	for field, expected := range assertions {
		if !strings.Contains(body, expected) {
			t.Errorf("expected %s field %q in output, got:\n%s", field, expected, body)
		}
	}

	if strings.Contains(body, "ciIdentifier") {
		t.Error("ciIdentifier should not appear when --ci-identifier is not provided")
	}
}

func TestScaffoldWithCIIdentifier(t *testing.T) {
	dir := t.TempDir()
	uuid := "01944b6e-40c0-7000-8000-000000000000"

	cmd := rootCmd
	cmd.SetArgs([]string{
		"scaffold", "CI Blip",
		"--content-dir", dir,
		"--ci-identifier", uuid,
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	filename := filepath.Join(dir, "ci-blip.mdx")
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading scaffolded file: %v", err)
	}

	body := string(content)
	expected := `ciIdentifier: "` + uuid + `"`
	if !strings.Contains(body, expected) {
		t.Errorf("expected ciIdentifier %q in output, got:\n%s", expected, body)
	}
}

func TestScaffoldWithQuadrantAndRing(t *testing.T) {
	dir := t.TempDir()

	cmd := rootCmd
	cmd.SetArgs([]string{
		"scaffold", "Custom Blip",
		"--content-dir", dir,
		"--quadrant", "tools",
		"--ring", "adopt",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	filename := filepath.Join(dir, "custom-blip.mdx")
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading scaffolded file: %v", err)
	}

	body := string(content)
	if !strings.Contains(body, `quadrant: "tools"`) {
		t.Errorf("expected quadrant tools, got:\n%s", body)
	}
	if !strings.Contains(body, `ring: "adopt"`) {
		t.Errorf("expected ring adopt, got:\n%s", body)
	}
}

func TestScaffoldFileAlreadyExists(t *testing.T) {
	dir := t.TempDir()

	// Create the file first
	existing := filepath.Join(dir, "existing.mdx")
	if err := os.WriteFile(existing, []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := rootCmd
	cmd.SetArgs([]string{
		"scaffold", "existing",
		"--content-dir", dir,
	})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for existing file, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", err)
	}
}
