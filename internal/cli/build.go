package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/infrashift/technology-radar/internal/build"
)

func newBuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Parse, validate, and generate radar artifacts",
		RunE:  runBuild,
	}
}

func runBuild(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig(cmd)
	if err != nil {
		return err
	}

	files, blips, err := loadBlips(cmd, cfg)
	if err != nil {
		return err
	}
	out := cmd.OutOrStdout()

	if err := build.WriteJSON(cfg.OutputJSON, blips); err != nil {
		return fmt.Errorf("writing JSON: %w", err)
	}
	fmt.Fprintf(out, "Generated %s\n", cfg.OutputJSON)

	// The frontend imports radar-config.json from next to the blip JSON.
	configPath := filepath.Join(filepath.Dir(cfg.OutputJSON), "radar-config.json")
	if err := build.WriteRadarConfig(configPath, cfg); err != nil {
		return fmt.Errorf("writing radar config: %w", err)
	}
	fmt.Fprintf(out, "Generated %s\n", configPath)

	if err := build.CopyMDXFiles(files, cfg.OutputContent); err != nil {
		return fmt.Errorf("copying MDX files: %w", err)
	}
	fmt.Fprintf(out, "Copied %d MDX files to %s\n", len(files), cfg.OutputContent)

	return nil
}
