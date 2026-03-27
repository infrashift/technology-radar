package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/infrashift/technology-radar/internal/build"
	"github.com/infrashift/technology-radar/internal/cuegen"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Parse, validate, and generate radar artifacts",
	RunE:  runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func runBuild(cmd *cobra.Command, args []string) error {
	contentDir := viper.GetString("content_dir")
	schemaPath := viper.GetString("schema_path")
	outputJSON := viper.GetString("output_json")
	outputContent := viper.GetString("output_content")
	version := viper.GetString("version")

	// Parse MDX files
	blips, rawFiles, err := build.ParseAllBlips(contentDir)
	if err != nil {
		return fmt.Errorf("parsing blips: %w", err)
	}
	fmt.Printf("Parsed %d blips from %s\n", len(blips), contentDir)

	// Validate against CUE schema
	if err := cuegen.Validate(schemaPath, blips, version); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	fmt.Println("All blips passed CUE schema validation")

	// Write JSON artifact
	if err := build.WriteJSON(outputJSON, blips); err != nil {
		return fmt.Errorf("writing JSON: %w", err)
	}
	fmt.Printf("Generated %s\n", outputJSON)

	// Copy MDX files to Astro content directory
	if err := build.CopyMDXFiles(rawFiles, outputContent); err != nil {
		return fmt.Errorf("copying MDX files: %w", err)
	}
	fmt.Printf("Copied %d MDX files to %s\n", len(rawFiles), outputContent)

	return nil
}
