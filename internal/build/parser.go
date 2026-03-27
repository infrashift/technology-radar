package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Blip represents the frontmatter of an MDX blip file.
type Blip struct {
	ID           string    `json:"id" yaml:"id"`
	Title        string    `json:"title" yaml:"title"`
	Quadrant     string    `json:"quadrant" yaml:"quadrant"`
	Ring         string    `json:"ring" yaml:"ring"`
	IsNew        bool      `json:"isNew" yaml:"isNew"`
	CIIdentifier string    `json:"ciIdentifier,omitempty" yaml:"ciIdentifier,omitempty"`
	History      []History `json:"history,omitempty" yaml:"history,omitempty"`
	RelatedBlips []string  `json:"relatedBlips,omitempty" yaml:"relatedBlips,omitempty"`
}

// History represents a historical ring placement.
type History struct {
	Date string `json:"date" yaml:"date"`
	Ring string `json:"ring" yaml:"ring"`
}

// RawBlipFile holds the original filename and raw bytes for copying.
type RawBlipFile struct {
	Filename string
	Content  []byte
}

// ParseAllBlips walks a directory and parses MDX frontmatter from all .mdx files.
func ParseAllBlips(dir string) ([]Blip, []RawBlipFile, error) {
	var blips []Blip
	var rawFiles []RawBlipFile

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".mdx") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, nil, fmt.Errorf("reading %s: %w", path, err)
		}

		blip, err := parseFrontmatter(content, entry.Name())
		if err != nil {
			return nil, nil, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}

		blips = append(blips, blip)
		rawFiles = append(rawFiles, RawBlipFile{
			Filename: entry.Name(),
			Content:  content,
		})
	}

	if len(blips) == 0 {
		return nil, nil, fmt.Errorf("no .mdx files found in %s", dir)
	}

	return blips, rawFiles, nil
}

func parseFrontmatter(content []byte, filename string) (Blip, error) {
	str := string(content)

	// Split on --- delimiters
	parts := strings.SplitN(str, "---", 3)
	if len(parts) < 3 {
		return Blip{}, fmt.Errorf("missing frontmatter delimiters in %s", filename)
	}

	frontmatter := strings.TrimSpace(parts[1])
	if frontmatter == "" {
		return Blip{}, fmt.Errorf("empty frontmatter in %s", filename)
	}

	var blip Blip
	if err := yaml.Unmarshal([]byte(frontmatter), &blip); err != nil {
		return Blip{}, fmt.Errorf("unmarshaling frontmatter in %s: %w", filename, err)
	}

	if blip.ID == "" {
		return Blip{}, fmt.Errorf("missing 'id' in frontmatter of %s", filename)
	}

	return blip, nil
}
