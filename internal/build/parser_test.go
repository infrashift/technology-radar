package build

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSplitFrontmatter(t *testing.T) {
	content := []byte("---\nid: \"test-blip\"\ntitle: \"A --- B\"\n---\n\nBody with --- dashes.\n\n---\n")

	fm, err := SplitFrontmatter(content)
	if err != nil {
		t.Fatalf("SplitFrontmatter failed: %v", err)
	}
	want := "id: \"test-blip\"\ntitle: \"A --- B\"\n"
	if string(fm) != want {
		t.Errorf("expected %q, got %q", want, fm)
	}
}

func TestSplitFrontmatterCRLFAndBOM(t *testing.T) {
	content := []byte("\xef\xbb\xbf---\r\nid: x\r\n---\r\nbody\r\n")

	fm, err := SplitFrontmatter(content)
	if err != nil {
		t.Fatalf("SplitFrontmatter failed: %v", err)
	}
	if string(fm) != "id: x\r\n" {
		t.Errorf("unexpected frontmatter %q", fm)
	}
}

func TestSplitFrontmatterErrors(t *testing.T) {
	tests := map[string]struct {
		content string
		want    string
	}{
		"no delimiters":  {"Just some content.\n", "missing frontmatter"},
		"text before":    {"intro\n---\nid: x\n---\n", "missing frontmatter"},
		"unterminated":   {"---\nid: x\nbody\n", "unterminated"},
		"empty":          {"---\n---\nbody\n", "empty frontmatter"},
		"inline dashes":  {"---\nid: x --- y\n", "unterminated"},
		"no newline end": {"---", "missing frontmatter"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := SplitFrontmatter([]byte(tc.content))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("expected error containing %q, got %v", tc.want, err)
			}
		})
	}
}

func TestReadBlipFiles(t *testing.T) {
	dir := t.TempDir()

	files := map[string]string{
		"blip-a.mdx": "---\nid: \"blip-a\"\n---\n\nContent A\n",
		"blip-b.mdx": "---\nid: \"blip-b\"\n---\n\nContent B\n",
		"readme.txt": "Not a blip",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "nested.mdx"), 0o755); err != nil {
		t.Fatal(err)
	}

	raw, err := ReadBlipFiles(dir)
	if err != nil {
		t.Fatalf("ReadBlipFiles failed: %v", err)
	}
	if len(raw) != 2 {
		t.Fatalf("expected 2 files, got %d", len(raw))
	}
	if raw[0].Filename != "blip-a.mdx" || raw[0].Path != filepath.Join(dir, "blip-a.mdx") {
		t.Errorf("unexpected file metadata: %+v", raw[0])
	}
	if string(raw[0].Frontmatter) != "id: \"blip-a\"\n" {
		t.Errorf("unexpected frontmatter %q", raw[0].Frontmatter)
	}
}

func TestReadBlipFilesBadFrontmatter(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "bad.mdx"), []byte("no frontmatter"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := ReadBlipFiles(dir)
	if err == nil || !strings.Contains(err.Error(), "bad.mdx") {
		t.Errorf("expected error naming bad.mdx, got %v", err)
	}
}

func TestReadBlipFilesEmptyDir(t *testing.T) {
	if _, err := ReadBlipFiles(t.TempDir()); err == nil {
		t.Fatal("expected error for empty directory")
	}
}
