package build

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/infrashift/technology-radar/internal/config"
)

// WriteJSON writes the blip data as JSON to the specified path.
func WriteJSON(path string, blips []Blip) error {
	return writeJSONFile(path, blips)
}

// radarConfig is the JSON structure written to radar-config.json.
type radarConfig struct {
	Version   string         `json:"version"`
	Quadrants []config.Entry `json:"quadrants"`
	Rings     []config.Entry `json:"rings"`
}

// WriteRadarConfig writes the version, quadrant, and ring definitions as JSON for the frontend.
func WriteRadarConfig(path string, cfg config.Config) error {
	return writeJSONFile(path, radarConfig{
		Version:   cfg.Version,
		Quadrants: cfg.Quadrants,
		Rings:     cfg.Rings,
	})
}

func writeJSONFile(path string, v any) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("marshaling %s: %w", path, err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	return nil
}

// CopyMDXFiles copies raw MDX files to the destination directory and removes
// any .mdx files there that no longer exist in the source.
func CopyMDXFiles(files []RawBlipFile, destDir string) error {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", destDir, err)
	}

	keep := make(map[string]bool, len(files))
	for _, f := range files {
		keep[f.Filename] = true
		destPath := filepath.Join(destDir, f.Filename)
		if err := os.WriteFile(destPath, f.Content, 0o644); err != nil {
			return fmt.Errorf("copying %s: %w", f.Filename, err)
		}
	}

	entries, err := os.ReadDir(destDir)
	if err != nil {
		return fmt.Errorf("reading %s: %w", destDir, err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".mdx") || keep[e.Name()] {
			continue
		}
		if err := os.Remove(filepath.Join(destDir, e.Name())); err != nil {
			return fmt.Errorf("removing stale %s: %w", e.Name(), err)
		}
	}

	return nil
}
