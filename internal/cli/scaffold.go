package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"

	"github.com/infrashift/technology-radar/internal/build"
)

func newScaffoldCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scaffold <blip-name>",
		Short: "Scaffold a new MDX blip file",
		Args:  cobra.ExactArgs(1),
		RunE:  runScaffold,
	}
	cmd.Flags().String("quadrant", "", "quadrant slug (default: first configured quadrant)")
	cmd.Flags().String("ring", "", "ring slug (default: third configured ring)")
	cmd.Flags().String("ci-identifier", "", "optional UUIDv7 CI identifier")
	return cmd
}

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

// slugify turns a display name into a blip id, e.g. "Node.js Edge" -> "node-js-edge".
func slugify(name string) string {
	return strings.Trim(nonSlugChars.ReplaceAllString(strings.ToLower(name), "-"), "-")
}

func runScaffold(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig(cmd)
	if err != nil {
		return err
	}

	name := args[0]
	quadrant, _ := cmd.Flags().GetString("quadrant")
	ring, _ := cmd.Flags().GetString("ring")
	ciIdentifier, _ := cmd.Flags().GetString("ci-identifier")

	if quadrant == "" {
		quadrant = cfg.Quadrants[0].Slug
	}
	if ring == "" {
		ring = cfg.Rings[2].Slug // third ring ("assess" by default)
	}
	if !slices.Contains(cfg.QuadrantSlugs(), quadrant) {
		return fmt.Errorf("invalid quadrant %q; allowed: %s", quadrant, strings.Join(cfg.QuadrantSlugs(), ", "))
	}
	if !slices.Contains(cfg.RingSlugs(), ring) {
		return fmt.Errorf("invalid ring %q; allowed: %s", ring, strings.Join(cfg.RingSlugs(), ", "))
	}

	id := slugify(name)
	if id == "" {
		return fmt.Errorf("cannot derive an id from %q; use letters or digits", name)
	}

	// Marshalling (rather than templating) keeps titles with quotes or colons valid YAML.
	var fm bytes.Buffer
	enc := yaml.NewEncoder(&fm)
	enc.SetIndent(2)
	err = enc.Encode(build.Blip{
		ID:           id,
		Title:        name,
		Quadrant:     quadrant,
		Ring:         ring,
		IsNew:        true,
		CIIdentifier: ciIdentifier,
		History:      []build.History{{Date: time.Now().Format("2006-01"), Ring: ring}},
	})
	if err != nil {
		return fmt.Errorf("encoding frontmatter: %w", err)
	}
	content := fmt.Sprintf("---\n%s---\n\nTODO: Add description for %s.\n", fm.String(), name)

	filename := filepath.Join(cfg.ContentDir, id+".mdx")
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("file already exists: %s", filename)
	}
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("writing %s: %w", filename, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", filename)
	return nil
}
