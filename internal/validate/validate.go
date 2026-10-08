// Package validate checks blip frontmatter against the CUE schema and enforces
// rules that span blips (ids, related blips, history).
package validate

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/yaml"

	"github.com/infrashift/technology-radar/internal/build"
	"github.com/infrashift/technology-radar/internal/config"
	"github.com/infrashift/technology-radar/schema"
)

// Error lists every problem found during validation.
type Error struct {
	Problems []string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d problem(s):\n  - %s", len(e.Problems), strings.Join(e.Problems, "\n  - "))
}

// Blips validates each file's frontmatter against #Blip (with the configured
// quadrants and rings) and returns the decoded blips in file order. All
// problems are collected and returned together as an *Error.
func Blips(cfg config.Config, files []build.RawBlipFile) ([]build.Blip, error) {
	ctx := cuecontext.New()
	s, err := schema.Compile(ctx)
	if err != nil {
		return nil, err
	}
	s = s.FillPath(cue.ParsePath("config"), ctx.Encode(cfg))
	blipDef := s.LookupPath(cue.ParsePath("#Blip"))
	if err := blipDef.Err(); err != nil {
		return nil, fmt.Errorf("looking up #Blip: %w", err)
	}

	var problems []string
	blips := make([]build.Blip, 0, len(files))
	paths := make([]string, 0, len(files))
	for _, f := range files {
		// Frontmatter starts on line 2; the leading newline keeps reported
		// line numbers aligned with the file.
		ast, err := yaml.Extract(f.Path, append([]byte("\n"), f.Frontmatter...))
		if err != nil {
			problems = append(problems, err.Error())
			continue
		}
		v := blipDef.Unify(ctx.BuildFile(ast))
		if err := v.Validate(cue.Concrete(true)); err != nil {
			problems = append(problems, schema.FormatErrors(err)...)
			continue
		}
		var b build.Blip
		if err := v.Decode(&b); err != nil {
			problems = append(problems, fmt.Sprintf("%s: decoding: %v", f.Path, err))
			continue
		}
		if want := strings.TrimSuffix(f.Filename, ".mdx"); b.ID != want {
			problems = append(problems, fmt.Sprintf("%s: id %q must match the filename (%q)", f.Path, b.ID, want))
		}
		blips = append(blips, b)
		paths = append(paths, f.Path)
	}

	known := make(map[string]bool, len(files))
	for _, f := range files {
		known[strings.TrimSuffix(f.Filename, ".mdx")] = true
	}
	problems = append(problems, crossCheck(blips, paths, known)...)
	if len(problems) > 0 {
		return nil, &Error{Problems: problems}
	}
	return blips, nil
}

// crossCheck enforces rules that need more than one blip or more than one
// field. paths[i] is the source file of blips[i], and known holds every blip
// id (ids are unique because each must match its filename).
func crossCheck(blips []build.Blip, paths []string, known map[string]bool) []string {
	var problems []string
	for i, b := range blips {
		where := paths[i]

		for _, rel := range b.RelatedBlips {
			switch {
			case rel == b.ID:
				problems = append(problems, fmt.Sprintf("%s: relatedBlips lists the blip itself", where))
			case !known[rel]:
				problems = append(problems, fmt.Sprintf("%s: relatedBlips references unknown blip %q", where, rel))
			}
		}

		// History is most-recent first; the frontend derives moved-in/out
		// from the first two entries.
		if len(b.History) > 0 && b.History[0].Ring != b.Ring {
			problems = append(problems, fmt.Sprintf("%s: latest history ring %q does not match ring %q", where, b.History[0].Ring, b.Ring))
		}
		for i := 1; i < len(b.History); i++ {
			if b.History[i].Date >= b.History[i-1].Date {
				problems = append(problems, fmt.Sprintf("%s: history must be newest first with unique dates (%s is listed after %s)", where, b.History[i].Date, b.History[i-1].Date))
				break
			}
		}
	}
	return problems
}
