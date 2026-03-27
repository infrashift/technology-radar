package cuegen

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cueerrors "cuelang.org/go/cue/errors"

	"github.com/infrashift/technology-radar/internal/build"
)

// manifest mirrors the CUE #Manifest definition.
type manifest struct {
	Version string       `json:"version"`
	Blips   []build.Blip `json:"blips"`
}

// Validate checks the blip data against the CUE schema at the given path.
func Validate(schemaPath string, blips []build.Blip, version string) error {
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("reading schema %s: %w", schemaPath, err)
	}

	ctx := cuecontext.New()

	// Compile the CUE schema
	schemaVal := ctx.CompileString(string(schemaBytes), cue.Filename(schemaPath))
	if schemaVal.Err() != nil {
		return fmt.Errorf("compiling CUE schema: %w", schemaVal.Err())
	}

	// Build the manifest JSON
	m := manifest{
		Version: version,
		Blips:   blips,
	}
	dataBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshaling manifest: %w", err)
	}

	// Compile the data as CUE
	dataVal := ctx.CompileBytes(dataBytes)
	if dataVal.Err() != nil {
		return fmt.Errorf("compiling data: %w", dataVal.Err())
	}

	// Look up the #Manifest definition and unify
	manifestDef := schemaVal.LookupPath(cue.ParsePath("#Manifest"))
	if manifestDef.Err() != nil {
		return fmt.Errorf("looking up #Manifest: %w", manifestDef.Err())
	}

	unified := manifestDef.Unify(dataVal)
	if err := unified.Validate(cue.Final()); err != nil {
		return formatValidationErrors(err)
	}

	return nil
}

func formatValidationErrors(err error) error {
	var msgs []string
	for _, e := range cueerrors.Errors(err) {
		msgs = append(msgs, fmt.Sprintf("  - %s", e.Error()))
	}
	return fmt.Errorf("schema violations:\n%s", strings.Join(msgs, "\n"))
}
