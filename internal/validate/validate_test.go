package validate

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/infrashift/technology-radar/internal/build"
	"github.com/infrashift/technology-radar/internal/config"
)

func defaultConfig(t *testing.T) config.Config {
	t.Helper()
	cfg, err := config.Load(filepath.Join(t.TempDir(), "none.yaml"), false)
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func file(name, frontmatter string) build.RawBlipFile {
	return build.RawBlipFile{
		Filename:    name,
		Path:        "content/blips/" + name,
		Frontmatter: []byte(frontmatter),
	}
}

const rust = `id: rust
title: Rust
quadrant: languages-and-frameworks
ring: trial
history:
  - date: "2025-11"
    ring: trial
  - date: "2025-05"
    ring: assess
relatedBlips: [react]
`

const react = `id: react
title: React
quadrant: languages-and-frameworks
ring: adopt
isNew: true
ciIdentifier: "01944b6e-40c0-7000-8000-000000000000"
`

func TestBlipsValid(t *testing.T) {
	blips, err := Blips(defaultConfig(t), []build.RawBlipFile{file("rust.mdx", rust), file("react.mdx", react)})
	if err != nil {
		t.Fatalf("expected valid blips, got: %v", err)
	}
	if len(blips) != 2 {
		t.Fatalf("expected 2 blips, got %d", len(blips))
	}
	if blips[0].History[1].Ring != "assess" || blips[0].RelatedBlips[0] != "react" {
		t.Errorf("unexpected decode: %+v", blips[0])
	}
	if !blips[1].IsNew || blips[1].CIIdentifier == "" {
		t.Errorf("unexpected decode: %+v", blips[1])
	}
}

func TestBlipsIsNewDefaultsFalse(t *testing.T) {
	blips, err := Blips(defaultConfig(t), []build.RawBlipFile{file("x.mdx", "id: x\ntitle: X\nquadrant: tools\nring: hold\n")})
	if err != nil {
		t.Fatal(err)
	}
	if blips[0].IsNew {
		t.Error("expected isNew to default to false")
	}
}

func TestBlipsUsesConfiguredSlugs(t *testing.T) {
	cfg := defaultConfig(t)
	cfg.Rings[0].Slug = "use"

	_, err := Blips(cfg, []build.RawBlipFile{file("x.mdx", "id: x\ntitle: X\nquadrant: tools\nring: use\n")})
	if err != nil {
		t.Fatalf("expected configured ring slug to be accepted, got: %v", err)
	}
	_, err = Blips(cfg, []build.RawBlipFile{file("x.mdx", "id: x\ntitle: X\nquadrant: tools\nring: adopt\n")})
	if err == nil {
		t.Fatal("expected removed ring slug to be rejected")
	}
}

func TestBlipsInvalid(t *testing.T) {
	tests := map[string]struct {
		files []build.RawBlipFile
		want  string
	}{
		"unknown key": {
			files: []build.RawBlipFile{file("x.mdx", "id: x\ntitle: X\nquadrant: tools\nring: hold\nrelatedBlip: [y]\n")},
			want:  "content/blips/x.mdx:6:1: relatedBlip: field not allowed",
		},
		"missing title": {
			files: []build.RawBlipFile{file("x.mdx", "id: x\nquadrant: tools\nring: hold\n")},
			want:  "title",
		},
		"bad quadrant": {
			files: []build.RawBlipFile{file("x.mdx", "id: x\ntitle: X\nquadrant: toolz\nring: hold\n")},
			want:  "content/blips/x.mdx:4:11: quadrant",
		},
		"bad ciIdentifier": {
			files: []build.RawBlipFile{file("x.mdx", "id: x\ntitle: X\nquadrant: tools\nring: hold\nciIdentifier: nope\n")},
			want:  "ciIdentifier",
		},
		"bad history date": {
			files: []build.RawBlipFile{file("x.mdx", "id: x\ntitle: X\nquadrant: tools\nring: hold\nhistory:\n  - {date: \"2025-13\", ring: hold}\n")},
			want:  "history.0.date",
		},
		"id mismatch": {
			files: []build.RawBlipFile{file("x.mdx", "id: y\ntitle: X\nquadrant: tools\nring: hold\n")},
			want:  `id "y" must match the filename ("x")`,
		},
		"unknown related": {
			files: []build.RawBlipFile{file("x.mdx", "id: x\ntitle: X\nquadrant: tools\nring: hold\nrelatedBlips: [ghost]\n")},
			want:  `references unknown blip "ghost"`,
		},
		"self related": {
			files: []build.RawBlipFile{file("x.mdx", "id: x\ntitle: X\nquadrant: tools\nring: hold\nrelatedBlips: [x]\n")},
			want:  "lists the blip itself",
		},
		"history ring mismatch": {
			files: []build.RawBlipFile{file("x.mdx", "id: x\ntitle: X\nquadrant: tools\nring: hold\nhistory:\n  - {date: \"2025-11\", ring: adopt}\n")},
			want:  `latest history ring "adopt" does not match ring "hold"`,
		},
		"history out of order": {
			files: []build.RawBlipFile{file("x.mdx", "id: x\ntitle: X\nquadrant: tools\nring: hold\nhistory:\n  - {date: \"2025-05\", ring: hold}\n  - {date: \"2025-11\", ring: adopt}\n")},
			want:  "newest first",
		},
		"invalid yaml": {
			files: []build.RawBlipFile{file("x.mdx", "id: [x\n")},
			want:  "content/blips/x.mdx",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Blips(defaultConfig(t), tc.files)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("expected error containing %q, got:\n%v", tc.want, err)
			}
		})
	}
}

func TestBlipsCollectsAllProblems(t *testing.T) {
	_, err := Blips(defaultConfig(t), []build.RawBlipFile{
		file("a.mdx", "id: a\ntitle: A\nquadrant: tools\nring: nope\n"),
		file("b.mdx", "id: b\ntitle: B\nquadrant: tools\nring: hold\nrelatedBlips: [ghost]\n"),
	})
	var verr *Error
	if !errors.As(err, &verr) {
		t.Fatalf("expected *Error, got %v", err)
	}
	joined := strings.Join(verr.Problems, "\n")
	if !strings.Contains(joined, "a.mdx") || !strings.Contains(joined, "b.mdx") {
		t.Errorf("expected problems from both files, got:\n%s", joined)
	}
}
