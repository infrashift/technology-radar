// Package schema embeds the CUE schema for radar configuration and blip
// frontmatter so the CLI always validates against the schema it was built with.
package schema

import (
	_ "embed"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"cuelang.org/go/cue"
	cueerrors "cuelang.org/go/cue/errors"
)

// Filename is the repo-relative path of the embedded schema, used in error positions.
const Filename = "schema/radar.cue"

//go:embed radar.cue
var Source string

// Compile compiles the embedded schema in ctx.
func Compile(ctx *cue.Context) (cue.Value, error) {
	v := ctx.CompileString(Source, cue.Filename(Filename))
	if err := v.Err(); err != nil {
		return cue.Value{}, fmt.Errorf("compiling %s: %w", Filename, err)
	}
	return v, nil
}

// FormatErrors flattens a CUE error into one line per problem, prefixed with
// the position in the user's file (not the schema) when one is known.
//
// A failed disjunction (an enum like #Ring, or a field with a default) makes
// CUE report why every alternative failed. Those groups are collapsed into
// the errors that matter: "x is not one of a, b, c" for enums, and only the
// non-conflict errors otherwise.
func FormatErrors(err error) []string {
	errs := cueerrors.Errors(err)
	var lines []string
	seen := make(map[string]bool)
	emit := func(e cueerrors.Error, msg string) {
		path := userPath(e)
		if msg == "field not allowed" {
			// An unknown list-valued key is reported at its first element
			// ("relatedBlip.0"); name the key itself.
			path = trailingIndexes.ReplaceAllString(path, "")
		}
		if path != "" {
			msg = path + ": " + msg
		}
		if pos := userPos(e); pos != "" {
			msg = pos + ": " + msg
		}
		if !seen[msg] {
			seen[msg] = true
			lines = append(lines, msg)
		}
	}

	for i := 0; i < len(errs); i++ {
		format, _ := errs[i].Msg()
		if format != "%d errors in empty disjunction:" {
			emit(errs[i], message(errs[i]))
			continue
		}

		path := userPath(errs[i])
		j := i + 1
		for j < len(errs) && within(userPath(errs[j]), path) {
			j++
		}
		members := errs[i+1 : j]
		i = j - 1

		var real, conflicts []cueerrors.Error
		for _, m := range members {
			switch f, _ := m.Msg(); f {
			case "conflicting values %s and %s":
				conflicts = append(conflicts, m)
			case "incompatible list lengths (%d and %d)":
			default:
				real = append(real, m)
			}
		}
		switch {
		case len(real) > 0:
			for _, m := range real {
				emit(m, message(m))
			}
		case len(conflicts) == len(members) && len(conflicts) > 0:
			var allowed []string
			var got string
			for _, m := range conflicts {
				_, args := m.Msg()
				allowed = append(allowed, fmt.Sprint(args[0]))
				got = fmt.Sprint(args[1])
			}
			sort.Strings(allowed)
			emit(conflicts[0], fmt.Sprintf("%s is not one of %s", got, strings.Join(allowed, ", ")))
		default:
			for _, m := range members {
				emit(m, message(m))
			}
		}
	}
	return lines
}

var trailingIndexes = regexp.MustCompile(`(\.\d+)+$`)

// within reports whether path is parent or one of its descendants.
func within(path, parent string) bool {
	return parent == "" || path == parent || strings.HasPrefix(path, parent+".")
}

func message(e cueerrors.Error) string {
	format, args := e.Msg()
	return fmt.Sprintf(format, args...)
}

// userPath returns the error's field path without the schema's top-level
// "config" field or "#Blip" definition.
func userPath(e cueerrors.Error) string {
	path := e.Path()
	if len(path) > 0 && (path[0] == "config" || path[0] == "#Blip") {
		path = path[1:]
	}
	return strings.Join(path, ".")
}

// userPos returns the first position outside the schema file, if any.
func userPos(e cueerrors.Error) string {
	for _, p := range cueerrors.Positions(e) {
		if p.IsValid() && p.Filename() != Filename {
			return fmt.Sprintf("%s:%d:%d", p.Filename(), p.Line(), p.Column())
		}
	}
	return ""
}
