package build

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// WriteJSON writes the blip data as JSON to the specified path.
func WriteJSON(path string, blips []Blip) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(blips, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	return nil
}

// CopyMDXFiles copies raw MDX files to the destination directory.
func CopyMDXFiles(files []RawBlipFile, destDir string) error {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", destDir, err)
	}

	for _, f := range files {
		destPath := filepath.Join(destDir, f.Filename)
		if err := os.WriteFile(destPath, f.Content, 0o644); err != nil {
			return fmt.Errorf("copying %s: %w", f.Filename, err)
		}
	}

	return nil
}
