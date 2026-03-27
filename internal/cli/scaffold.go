package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold <blip-name>",
	Short: "Scaffold a new MDX blip file",
	Args:  cobra.ExactArgs(1),
	RunE:  runScaffold,
}

func init() {
	scaffoldCmd.Flags().String("quadrant", "techniques", "quadrant (techniques, platforms, tools, languages-and-frameworks)")
	scaffoldCmd.Flags().String("ring", "assess", "ring (adopt, trial, assess, hold)")
	scaffoldCmd.Flags().String("ci-identifier", "", "optional UUIDv7 CI identifier")
	rootCmd.AddCommand(scaffoldCmd)
}

const blipTemplate = `---
id: "{{.ID}}"
title: "{{.Title}}"
quadrant: "{{.Quadrant}}"
ring: "{{.Ring}}"
isNew: true
{{- if .CIIdentifier}}
ciIdentifier: "{{.CIIdentifier}}"
{{- end}}
history:
  - date: "{{.Date}}"
    ring: "{{.Ring}}"
---

TODO: Add description for {{.Title}}.
`

func runScaffold(cmd *cobra.Command, args []string) error {
	name := args[0]
	contentDir := viper.GetString("content_dir")
	quadrant, _ := cmd.Flags().GetString("quadrant")
	ring, _ := cmd.Flags().GetString("ring")
	ciIdentifier, _ := cmd.Flags().GetString("ci-identifier")

	id := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	filename := filepath.Join(contentDir, id+".mdx")

	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("file already exists: %s", filename)
	}

	tmpl, err := template.New("blip").Parse(blipTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	data := struct {
		ID, Title, Quadrant, Ring, Date, CIIdentifier string
	}{
		ID:           id,
		Title:        name,
		Quadrant:     quadrant,
		Ring:         ring,
		Date:         time.Now().Format("2006-01"),
		CIIdentifier: ciIdentifier,
	}

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("writing template: %w", err)
	}

	fmt.Printf("Created %s\n", filename)
	return nil
}
