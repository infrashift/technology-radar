package build

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
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
