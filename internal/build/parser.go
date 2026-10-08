package build

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// RawBlipFile holds a blip source file and its extracted frontmatter.
type RawBlipFile struct {
	Filename    string // base name, e.g. "rust.mdx"
	Path        string // path as read, used in error messages
	Content     []byte
	Frontmatter []byte // YAML between the --- delimiters, starting on line 2 of the file
}

// ReadBlipFiles reads every .mdx file in dir and extracts its frontmatter.
func ReadBlipFiles(dir string) ([]RawBlipFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}

	var files []RawBlipFile
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".mdx") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}

		fm, err := SplitFrontmatter(content)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}

		files = append(files, RawBlipFile{
			Filename:    entry.Name(),
			Path:        path,
			Content:     content,
			Frontmatter: fm,
		})
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no .mdx files found in %s", dir)
	}

	return files, nil
}

// SplitFrontmatter returns the YAML frontmatter of an MDX file. The file must
// start with a "---" line, and the frontmatter ends at the next line that is
// exactly "---".
func SplitFrontmatter(content []byte) ([]byte, error) {
	content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))

	first, rest, ok := bytes.Cut(content, []byte("\n"))
	if !ok || !isDelimiter(first) {
		return nil, fmt.Errorf("missing frontmatter: file must start with a --- line")
	}

	for off := 0; off < len(rest); {
		line, _, found := bytes.Cut(rest[off:], []byte("\n"))
		if isDelimiter(line) {
			fm := rest[:off]
			if len(bytes.TrimSpace(fm)) == 0 {
				return nil, fmt.Errorf("empty frontmatter")
			}
			return fm, nil
		}
		if !found {
			break
		}
		off += len(line) + 1
	}
	return nil, fmt.Errorf("unterminated frontmatter: no closing --- line")
}

func isDelimiter(line []byte) bool {
	return string(bytes.TrimRight(line, " \t\r")) == "---"
}
